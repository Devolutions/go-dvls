package dvls

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Login(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"tokenId":"mock-token-123"}`))
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("true"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient("test-key", "test-secret", server.URL)
	require.NoError(t, err)
	assert.Equal(t, "mock-token-123", client.credential.token)
	assert.NotNil(t, client.Entries)
	assert.NotNil(t, client.Vaults)
}

func TestNewClientWithContext_Login(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"tokenId":"mock-token-123"}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClientWithContext(ctx, "test-key", "test-secret", server.URL)
	require.NoError(t, err)
	assert.Equal(t, "mock-token-123", client.credential.token)
	assert.NotNil(t, client.Entries)
	assert.NotNil(t, client.Vaults)
}

func TestNewClientWithContext_LoginCanceled(t *testing.T) {
	loginReceived := make(chan struct{})
	release := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
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

type sessionMock struct {
	logins        atomic.Int32
	isLoggedCalls atomic.Int32
	vaultCalls    atomic.Int32
	isLoggedBody  string
	validToken    string
}

func (m *sessionMock) vault(w http.ResponseWriter, r *http.Request) {
	m.vaultCalls.Add(1)
	if m.validToken == "" || r.Header.Get("tokenId") != m.validToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"result":1}`))
}

func newSessionTestClient(t *testing.T, mux *http.ServeMux, mock *sessionMock) *Client {
	t.Helper()

	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Empty(t, r.Header.Get("tokenId"), "login must not send the expired token")
		mock.logins.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tokenId":"new-token"}`))
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		mock.isLoggedCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mock.isLoggedBody))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &Client{
		baseUri:    server.URL,
		client:     server.Client(),
		credential: &credentials{appKey: "test-key", appSecret: "test-secret", token: "test-token"},
	}
}

func TestRequest_NoSessionCheckOnSuccess(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "true", validToken: "test-token"}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", mock.vault)
	client := newSessionTestClient(t, mux, mock)

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
	require.NoError(t, err)
	assert.EqualValues(t, 1, mock.vaultCalls.Load())
	assert.Zero(t, mock.isLoggedCalls.Load(), "a successful request must not call the is-logged endpoint")
	assert.Zero(t, mock.logins.Load())
}

func TestRequest_ReloginAfterTokenExpiry(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "false", validToken: "new-token"}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", mock.vault)
	client := newSessionTestClient(t, mux, mock)

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
	require.NoError(t, err)
	assert.Equal(t, "new-token", client.credential.token)
	assert.EqualValues(t, 2, mock.vaultCalls.Load())
	assert.EqualValues(t, 1, mock.logins.Load())
}

func TestRequest_RetryResendsBody(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "false", validToken: "new-token"}
	var mu sync.Mutex
	var bodies []string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		mu.Lock()
		bodies = append(bodies, string(body))
		mu.Unlock()
		mock.vault(w, r)
	})
	client := newSessionTestClient(t, mux, mock)

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodPost, bytes.NewBufferString(`{"name":"vault"}`))
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{`{"name":"vault"}`, `{"name":"vault"}`}, bodies)
}

func TestRequest_OtherStatusNotRetried(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "false"}
	var calls atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	})
	client := newSessionTestClient(t, mux, mock)

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
	require.Error(t, err)
	assert.EqualValues(t, 1, calls.Load(), "only a 401 may trigger a retry")
	assert.Zero(t, mock.isLoggedCalls.Load())
	assert.Zero(t, mock.logins.Load())
}

func TestRequest_UnauthorizedWithValidSession(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "true"}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", mock.vault)
	client := newSessionTestClient(t, mux, mock)

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
	assert.True(t, hasStatusCode(err, http.StatusUnauthorized), "expected a 401, got %v", err)
	assert.EqualValues(t, 1, mock.vaultCalls.Load(), "a 401 on a valid session must not be retried")
	assert.EqualValues(t, 1, mock.isLoggedCalls.Load())
	assert.Zero(t, mock.logins.Load(), "a 401 on a valid session must not trigger a login")
}

func TestRequest_ConcurrentUnauthorizedSingleLogin(t *testing.T) {
	const requests = 8

	mock := &sessionMock{isLoggedBody: "false", validToken: "new-token"}
	var arrived sync.WaitGroup
	arrived.Add(requests)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("tokenId") != mock.validToken {
			arrived.Done()
			arrived.Wait()
		}
		mock.vault(w, r)
	})
	client := newSessionTestClient(t, mux, mock)

	errs := make(chan error, requests)
	for range requests {
		go func() {
			_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
			errs <- err
		}()
	}
	for range requests {
		assert.NoError(t, <-errs)
	}

	assert.EqualValues(t, 1, mock.logins.Load(), "concurrent 401s must share one login")
	assert.EqualValues(t, 1, mock.isLoggedCalls.Load())
}

func TestNewClient_ServicesShareTheClientSession(t *testing.T) {
	var mu sync.Mutex
	logins := 0
	valid := map[string]bool{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		logins++
		token := fmt.Sprintf("token-%d", logins)
		valid[token] = true
		w.Write([]byte(`{"tokenId":"` + token + `"}`))
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if valid[r.Header.Get("tokenId")] {
			w.Write([]byte("true"))
			return
		}
		w.Write([]byte("false"))
	})
	mux.HandleFunc("/api/private-instance-information", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if !valid[r.Header.Get("tokenId")] {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"result":1,"data":{}}`))
	})
	mux.HandleFunc("/api/v1/vault/vault-id", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if !valid[r.Header.Get("tokenId")] {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"id":"vault-id"}`))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient("test-key", "test-secret", server.URL)
	require.NoError(t, err)

	mu.Lock()
	valid = map[string]bool{}
	mu.Unlock()

	_, err = client.GetPrivateServerInfo()
	require.NoError(t, err)
	_, err = client.Vaults.Get("vault-id")
	require.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 2, logins, "a session renewed through the client must be the one its services use")
}

func TestNewClientWithApiKey(t *testing.T) {
	mux := http.NewServeMux()
	loginCalled := false
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClientWithApiKey("test-api-key", server.URL)
	require.NoError(t, err)
	assert.True(t, client.credential.useApiKey())
	assert.Equal(t, "test-api-key", client.credential.apiKey)
	assert.NotNil(t, client.Entries)
	assert.NotNil(t, client.Vaults)
	assert.False(t, loginCalled, "API key client must not call the login endpoint")
}

func TestRequestWithApiKey_BearerHeader(t *testing.T) {
	mux := http.NewServeMux()
	loginCalled := false
	isLoggedCalled := false
	mux.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
	})
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		isLoggedCalled = true
	})

	var gotAuth, gotTokenId string
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotTokenId = r.Header.Get("tokenId")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClientWithApiKey("test-api-key", server.URL)
	require.NoError(t, err)

	_, err = client.Request(server.URL+"/api/data", http.MethodGet, nil)
	require.NoError(t, err)

	assert.Equal(t, "Bearer test-api-key", gotAuth)
	assert.Empty(t, gotTokenId)
	assert.False(t, loginCalled, "API key client must not call the login endpoint")
	assert.False(t, isLoggedCalled, "API key client must not call the is-logged endpoint")
}

func TestRequestWithApiKey_UnauthorizedNotRetried(t *testing.T) {
	mock := &sessionMock{isLoggedBody: "false"}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", mock.vault)
	client := newSessionTestClient(t, mux, mock)
	client.credential = &credentials{apiKey: "test-api-key"}

	_, err := client.Request(client.baseUri+"/api/v1/vault", http.MethodGet, nil)
	assert.True(t, hasStatusCode(err, http.StatusUnauthorized), "expected a 401, got %v", err)
	assert.EqualValues(t, 1, mock.vaultCalls.Load(), "API key client must not retry a 401")
	assert.Zero(t, mock.isLoggedCalls.Load(), "API key client must not call the is-logged endpoint")
	assert.Zero(t, mock.logins.Load(), "API key client must not call the login endpoint")
}

func TestIsLogged_True(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	logged, err := client.isLogged()
	require.NoError(t, err)
	assert.True(t, logged)
}

func TestIsLogged_False(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("false"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := &Client{
		baseUri:    server.URL,
		client:     server.Client(),
		credential: &credentials{token: "test-token"},
	}

	logged, err := client.isLogged()
	require.NoError(t, err)
	assert.False(t, logged)
}
