package dvls

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Login(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Empty(t, r.Header.Get("tokenId"))

		var body map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]string{"appKey": "test-key", "appSecret": "test-secret"}, body)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"mock-token-123"}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient("test-key", "test-secret", server.URL)
	require.NoError(t, err)
	assert.Equal(t, "mock-token-123", client.credential.currentToken())
	assert.NotNil(t, client.Entries)
	assert.NotNil(t, client.Vaults)
}

func TestNewClient_LoginRejected(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"The application key or secret is invalid."}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	_, err := NewClient("test-key", "wrong", server.URL)
	require.Error(t, err)
	assert.ErrorContains(t, err, "The application key or secret is invalid.")
}

func TestNewClient_LoginWithoutToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	_, err := NewClient("test-key", "test-secret", server.URL)
	require.Error(t, err)
	assert.ErrorContains(t, err, "no tokenId")
}

func TestNewClientWithContext_Login(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"mock-token-123"}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClientWithContext(ctx, "test-key", "test-secret", server.URL)
	require.NoError(t, err)
	assert.Equal(t, "mock-token-123", client.credential.currentToken())
}

func TestNewClientWithContext_LoginCanceled(t *testing.T) {
	loginReceived := make(chan struct{})
	release := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		close(loginReceived)
		<-release
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		_, err := NewClientWithContext(ctx, "test-key", "test-secret", server.URL)
		errCh <- err
	}()

	select {
	case <-loginReceived:
	case <-time.After(5 * time.Second):
		t.Fatal("login request never reached the server")
	}
	cancel()

	select {
	case err := <-errCh:
		require.Error(t, err)
		assert.ErrorContains(t, err, context.Canceled.Error())
	case <-time.After(5 * time.Second):
		t.Fatal("NewClientWithContext did not return after the context was canceled")
	}
}

func TestNewClientWithContext_LoginDeadlineExceeded(t *testing.T) {
	release := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		<-release
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client, err := NewClientWithContext(ctx, "test-key", "test-secret", server.URL)
	require.Error(t, err)
	assert.ErrorContains(t, err, context.DeadlineExceeded.Error())
	assert.Equal(t, Client{}, client)
}

func TestRequest_ReloginOnUnauthorized(t *testing.T) {
	var logins atomic.Int32
	var bodies []string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Empty(t, r.Header.Get("tokenId"), "login must not send the expired token")
		logins.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"new-token"}`))
	})
	mux.HandleFunc("/api/v2/vault", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(body))
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("tokenId") != "new-token" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Missing or invalid token."}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"v1"}`))
	})

	client := newTestClient(t, mux)

	resp, err := client.Request(client.baseUri+"/api/v2/vault", http.MethodPost, strings.NewReader(`{"name":"a"}`))
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"v1"}`, string(resp.Response))
	assert.Equal(t, int32(1), logins.Load())
	assert.Equal(t, "new-token", client.credential.currentToken())
	assert.Equal(t, []string{`{"name":"a"}`, `{"name":"a"}`}, bodies, "the retried request must resend the body")
}

func TestRequest_ReloginRetriesOnce(t *testing.T) {
	var calls atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"new-token"}`))
	})
	mux.HandleFunc("/api/v2/me", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Missing or invalid token."}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Request(client.baseUri+"/api/v2/me", http.MethodGet, nil)
	require.Error(t, err)
	assert.True(t, hasStatusCode(err, http.StatusUnauthorized))
	assert.Equal(t, int32(2), calls.Load())
}

func TestRequest_ReloginFailure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"The application key or secret is invalid."}`))
	})
	mux.HandleFunc("/api/v2/me", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	})

	client := newTestClient(t, mux)

	_, err := client.Request(client.baseUri+"/api/v2/me", http.MethodGet, nil)
	require.Error(t, err)
	assert.ErrorContains(t, err, "failed to refresh login token")
}

func TestRequest_ReloginSkippedWhenTokenAlreadyRefreshed(t *testing.T) {
	var logins atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		logins.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"new-token"}`))
	})

	client := newTestClient(t, mux)
	client.credential.token = "token-from-another-request"

	err := client.reloginWithContext(context.Background(), "test-token")
	require.NoError(t, err)
	assert.Equal(t, int32(0), logins.Load())
	assert.Equal(t, "token-from-another-request", client.credential.currentToken())
}

func TestRequest_ApiKeyDoesNotRelogin(t *testing.T) {
	var logins atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		logins.Add(1)
	})
	mux.HandleFunc("/api/v2/me", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClientWithApiKey("test-api-key", server.URL)
	require.NoError(t, err)

	_, err = client.Request(server.URL+"/api/v2/me", http.MethodGet, nil)
	require.Error(t, err)
	assert.Equal(t, int32(0), logins.Load())
}

func TestNewClientWithApiKey(t *testing.T) {
	var logins atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		logins.Add(1)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClientWithApiKey("test-api-key", server.URL)
	require.NoError(t, err)
	assert.True(t, client.credential.useApiKey())
	assert.NotNil(t, client.Entries)
	assert.Equal(t, int32(0), logins.Load(), "API key client must not call the login endpoint")
}

func TestRequestWithApiKey_BearerHeader(t *testing.T) {
	var gotAuth, gotTokenId string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/server/info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		gotAuth = r.Header.Get("Authorization")
		gotTokenId = r.Header.Get("tokenId")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"2027.1.0.0"}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClientWithApiKey("test-api-key", server.URL)
	require.NoError(t, err)

	info, err := client.GetServerInfo()
	require.NoError(t, err)
	assert.Equal(t, "2027.1.0.0", info.Version)
	assert.Equal(t, "Bearer test-api-key", gotAuth)
	assert.Empty(t, gotTokenId)
}

func TestRequest_ContentTypeOption(t *testing.T) {
	var gotContentType string

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		gotContentType = r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Request(client.baseUri+"/upload", http.MethodPost, strings.NewReader("raw"), RequestOptions{ContentType: "application/octet-stream"})
	require.NoError(t, err)
	assert.Equal(t, "application/octet-stream", gotContentType)

	_, err = client.Request(client.baseUri+"/upload", http.MethodPost, strings.NewReader("{}"))
	require.NoError(t, err)
	assert.Equal(t, "application/json", gotContentType)
}

func TestRequest_NoContent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault/"+testVaultID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, mux)

	require.NoError(t, client.Vaults.Delete(testVaultID))
}

func TestRequestError_Messages(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":[{"field":"name","message":"name is required."},{"field":"contentType","message":"Invalid value 'x'."}]}`))
	})
	mux.HandleFunc("/api/v2/vault/"+testVaultID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"The requested resource was not found."}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Vaults.New(Vault{})
	require.Error(t, err)
	assert.ErrorContains(t, err, "(status 400)")
	assert.ErrorContains(t, err, "name: name is required.; contentType: Invalid value 'x'.")

	_, err = client.Vaults.Get(testVaultID)
	require.Error(t, err)
	assert.True(t, IsNotFound(err))
	assert.ErrorContains(t, err, "The requested resource was not found.")
}
