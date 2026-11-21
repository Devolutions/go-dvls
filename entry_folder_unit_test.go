package dvls

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFolderGetEntries_TypeFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "Folder", r.URL.Query().Get("type"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[{"id":"1","name":"F","path":"F","type":"Folder","subType":"Server","data":{"domain":"d","username":"u","credentialMode":"Custom"}}]}`))
	})

	client := newTestClient(t, mux)

	folders, err := client.Entries.Folder.GetEntries(testVaultID, GetEntriesOptions{})
	require.NoError(t, err)
	require.Len(t, folders, 1)
	data, ok := folders[0].GetFolderData()
	require.True(t, ok)
	assert.Equal(t, EntryFolderData{Domain: "d", Username: "u"}, *data)
}

func TestFolderGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError error
	}{
		{name: "Found", data: `{"id":"` + testEntryID + `","name":"F","path":"A","type":"Folder","subType":"Folder","data":{}}`},
		{name: "NotFound", data: "", wantError: ErrEntryNotFound},
		{name: "Multiple", data: `{"id":"1","name":"F","path":"A","type":"Folder","subType":"Folder","data":{}},{"id":"2","name":"F","path":"A\\F","type":"Folder","subType":"Folder","data":{}}`, wantError: ErrMultipleEntriesFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "F", r.URL.Query().Get("name"))
				assert.Equal(t, "A", r.URL.Query().Get("path"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":2,"totalPage":1,"data":[` + tt.data + `]}`))
			})
			mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id":"` + testEntryID + `","name":"F","path":"A","type":"Folder","subType":"Folder","data":{}}`))
			})

			client := newTestClient(t, mux)

			path := "A"
			folder, err := client.Entries.Folder.GetByName(testVaultID, "F", GetByNameOptions{Path: &path})
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testEntryID, folder.Id)
		})
	}
}

func TestFolderNewAndUpdate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"name": "F", "path": "A", "type": "Folder", "subType": "Folder", "data": map[string]any{}}, body, "a nil Data is sent as an empty object")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"` + testEntryID + `"}`))
	})
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "A", body["path"])
			assert.Equal(t, map[string]any{}, body["data"], "empty folder credentials are not sent, so the stored ones are kept")
			w.Write([]byte(`{}`))
		case http.MethodGet:
			w.Write([]byte(`{"id":"` + testEntryID + `","name":"F","path":"A","type":"Folder","subType":"Folder","data":{}}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, mux)

	id, err := client.Entries.Folder.New(Entry{VaultId: testVaultID, Name: "F", Path: "A", Type: EntryFolderType, SubType: EntryFolderSubTypeFolder})
	require.NoError(t, err)

	folder, err := client.Entries.Folder.Update(Entry{Id: id, VaultId: testVaultID, Name: "F", Path: "A", Type: EntryFolderType, SubType: EntryFolderSubTypeFolder, Data: &EntryFolderData{}})
	require.NoError(t, err)
	assert.Equal(t, "A", folder.Path)

	require.NoError(t, client.Entries.Folder.DeleteById(testVaultID, id))
}
