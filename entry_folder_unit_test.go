package dvls

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFolderGetEntries(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"1","name":"Folder1","type":"Folder","subType":"Folder","path":"test","data":{"domain":"d1"}},
				{"id":"2","name":"Cred1","type":"Credential","subType":"Default","path":"test","data":{"username":"u1"}},
				{"id":"3","name":"Folder2","type":"Folder","subType":"Database","path":"test","data":{"domain":"d2"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 3,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	entries, err := client.Entries.Folder.GetEntries(testVaultID, GetEntriesOptions{})
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "Folder1", entries[0].Name)
	assert.Equal(t, "Folder2", entries[1].Name)
}

func TestFolderGetByName(t *testing.T) {
	entryID := "folder-found"
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"folder-found","name":"Target","type":"Folder","subType":"Folder","path":"test","data":{"domain":"d"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 1,
			"pageSize": 20
		}`))
	})
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry/%s", testVaultID, entryID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"id": "folder-found",
			"name": "Target",
			"type": "Folder",
			"subType": "Folder",
			"path": "test",
			"data": {"domain": "d", "username": "u"}
		}`))
	})

	client := newTestClient(t, mux)

	entry, err := client.Entries.Folder.GetByName(testVaultID, "Target", GetByNameOptions{})
	require.NoError(t, err)
	assert.Equal(t, "folder-found", entry.Id)
	assert.Equal(t, "Target", entry.Name)
}

func TestFolderGetByName_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 0,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Entries.Folder.GetByName(testVaultID, "NonExistent", GetByNameOptions{})
	assert.ErrorIs(t, err, ErrEntryNotFound)
}
