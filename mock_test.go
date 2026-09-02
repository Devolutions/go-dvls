package dvls

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testVaultID = "test-vault-id"

const testEntryID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

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
