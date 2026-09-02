package dvls

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testVaultID = "test-vault-id"

const (
	testRoleID     = "11111111-2222-3333-4444-555555555555"
	testAssigneeID = "66666666-7777-8888-9999-000000000000"
	testUserID     = "12121212-3434-5656-7878-909090909090"
)

func newTestClient(t *testing.T, mux *http.ServeMux) *Client {
	t.Helper()

	mux.HandleFunc("/api/is-logged", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("true"))
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := &Client{
		baseUri:    server.URL,
		client:     server.Client(),
		credential: credentials{token: "test-token"},
	}
	client.initServices()

	return client
}
