package dvls

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		credential: credentials{token: "test-token"},
	}

	logged, err := client.isLogged()
	require.NoError(t, err)
	assert.False(t, logged)
}
