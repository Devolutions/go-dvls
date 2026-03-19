package dvls

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialGetEntries(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"1","name":"Cred1","type":"Credential","subType":"Default","path":"test","data":{"username":"u1"}},
				{"id":"2","name":"Folder1","type":"Folder","subType":"Folder","path":"test","data":{"domain":"d1"}},
				{"id":"3","name":"Cred2","type":"Credential","subType":"ApiKey","path":"test","data":{"apiId":"a1"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 3,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	entries, err := client.Entries.Credential.GetEntries(testVaultID, GetEntriesOptions{})
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "Cred1", entries[0].Name)
	assert.Equal(t, "Cred2", entries[1].Name)
}

func TestCredentialGetEntries_PathFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"1","name":"InBase","type":"Credential","subType":"Default","path":"base","data":{"username":"u"}},
				{"id":"2","name":"InSub","type":"Credential","subType":"Default","path":"base\\sub","data":{"username":"u"}},
				{"id":"3","name":"InSimilar","type":"Credential","subType":"Default","path":"baseother","data":{"username":"u"}},
				{"id":"4","name":"InRoot","type":"Credential","subType":"Default","path":"","data":{"username":"u"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 4,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	basePath := "base"
	entries, err := client.Entries.Credential.GetEntries(testVaultID, GetEntriesOptions{Path: &basePath})
	require.NoError(t, err)
	require.Len(t, entries, 2)

	var names []string
	for _, e := range entries {
		names = append(names, e.Name)
	}
	assert.Contains(t, names, "InBase")
	assert.Contains(t, names, "InSub")
	assert.NotContains(t, names, "InSimilar")
	assert.NotContains(t, names, "InRoot")
}

func TestCredentialGetEntries_RootPathFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"1","name":"InBase","type":"Credential","subType":"Default","path":"base","data":{"username":"u"}},
				{"id":"2","name":"InRoot","type":"Credential","subType":"Default","path":"","data":{"username":"u"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 2,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	rootPath := ""
	entries, err := client.Entries.Credential.GetEntries(testVaultID, GetEntriesOptions{Path: &rootPath})
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "InRoot", entries[0].Name)
}

func TestCredentialGetByName(t *testing.T) {
	entryID := "entry-found"
	mux := http.NewServeMux()
	callCount := 0
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"entry-found","name":"Target","type":"Credential","subType":"Default","path":"test","data":{"username":"u"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 1,
			"pageSize": 20
		}`))
	})
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry/%s", testVaultID, entryID), func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"id": "entry-found",
			"name": "Target",
			"type": "Credential",
			"subType": "Default",
			"path": "test",
			"data": {"username": "u", "password": "p"}
		}`))
	})

	client := newTestClient(t, mux)

	testPath := "test"
	entry, err := client.Entries.Credential.GetByName(testVaultID, "Target", EntryCredentialSubTypeDefault, GetByNameOptions{Path: &testPath})
	require.NoError(t, err)
	assert.Equal(t, "entry-found", entry.Id)
	assert.Equal(t, "Target", entry.Name)
	assert.Equal(t, 1, callCount)
}

func TestCredentialGetByName_NotFound(t *testing.T) {
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

	_, err := client.Entries.Credential.GetByName(testVaultID, "NonExistent", EntryCredentialSubTypeDefault, GetByNameOptions{})
	assert.ErrorIs(t, err, ErrEntryNotFound)
}

func TestCredentialGetByName_Multiple(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/api/v1/vault/%s/entry", testVaultID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"result": 1,
			"data": [
				{"id":"1","name":"Dup","type":"Credential","subType":"Default","path":"","data":{"username":"u1"}},
				{"id":"2","name":"Dup","type":"Credential","subType":"Default","path":"","data":{"username":"u2"}}
			],
			"currentPage": 1,
			"totalPage": 1,
			"totalCount": 2,
			"pageSize": 20
		}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Entries.Credential.GetByName(testVaultID, "Dup", EntryCredentialSubTypeDefault, GetByNameOptions{})
	assert.ErrorIs(t, err, ErrMultipleEntriesFound)
}

func TestCredentialValidation_EmptyVaultId(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	_, err := client.Entries.Credential.New(Entry{
		Type:    EntryCredentialType,
		SubType: EntryCredentialSubTypeDefault,
		Data:    &EntryCredentialDefaultData{},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "VaultId")
}

func TestCredentialValidation_UnsupportedSubType(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	_, err := client.Entries.Credential.New(Entry{
		VaultId: testVaultID,
		Type:    EntryCredentialType,
		SubType: "InvalidSubType",
		Data:    &EntryCredentialDefaultData{},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}
