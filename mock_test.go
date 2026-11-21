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

const testEntryID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

const (
	testVaultEntriesPath = "/api/v2/vault/" + testVaultID + "/entry"
	testEntryPath        = testVaultEntriesPath + "/" + testEntryID
)

func newTestClient(t *testing.T, mux *http.ServeMux) *Client {
	t.Helper()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := &Client{
		baseUri:    server.URL,
		client:     server.Client(),
		credential: &credentials{appKey: "test-key", appSecret: "test-secret", token: "test-token"},
	}
	client.initServices()

	return client
}
