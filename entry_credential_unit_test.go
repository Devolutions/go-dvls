package dvls

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCredentialJSON = `{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"cred","path":"A","type":"Credential","subType":"Default","description":"","tags":["a b","c"],"data":{"username":"u","password":"p","domain":"d"},"createdOn":"2026-11-21T17:38:35Z"}`

func TestCredentialGetEntries_Filters(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		query := r.URL.Query()
		assert.Equal(t, "Credential", query.Get("type"))
		assert.Equal(t, "cred", query.Get("name"))
		assert.Equal(t, "A", query.Get("path"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":2,"totalPage":1,"data":[
			{"id":"1","name":"cred","path":"A","type":"Credential","subType":"Default","data":{"username":"u1"}},
			{"id":"2","name":"cred","path":"A\\B","type":"Credential","subType":"Unknown","data":{}}
		]}`))
	})

	client := newTestClient(t, mux)

	name, path := "cred", "A"
	entries, err := client.Entries.Credential.GetEntries(testVaultID, GetEntriesOptions{Name: &name, Path: &path})
	require.NoError(t, err)
	require.Len(t, entries, 1, "an unsupported subtype is skipped")
	assert.Equal(t, testVaultID, entries[0].VaultId)
	data, ok := entries[0].GetCredentialDefaultData()
	require.True(t, ok)
	assert.Equal(t, "u1", data.Username)
}

func TestCredentialGetEntries_RootPathFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.False(t, r.URL.Query().Has("path"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":2,"totalPage":1,"data":[
			{"id":"1","name":"root","path":"","type":"Credential","subType":"Default","data":{}},
			{"id":"2","name":"nested","path":"A","type":"Credential","subType":"Default","data":{}}
		]}`))
	})

	client := newTestClient(t, mux)

	root := ""
	entries, err := client.Entries.Credential.GetEntries(testVaultID, GetEntriesOptions{Path: &root})
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "root", entries[0].Name)
}

func TestCredentialGetById(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testCredentialJSON))
	})

	client := newTestClient(t, mux)

	entry, err := client.Entries.Credential.GetById(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, []string{"a b", "c"}, entry.Tags)
	require.NotNil(t, entry.CreatedOn)
	assert.Equal(t, 2026, entry.CreatedOn.Year())
	data, ok := entry.GetCredentialDefaultData()
	require.True(t, ok)
	assert.Equal(t, EntryCredentialDefaultData{Username: "u", Password: "p", Domain: "d"}, *data)
}

func TestCredentialNew(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{
			"name":    "cred",
			"path":    "A",
			"type":    "Credential",
			"subType": "AccessCode",
			"tags":    []any{"t"},
			"data":    map[string]any{"password": "secret"},
		}, body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(testCredentialJSON))
	})

	client := newTestClient(t, mux)

	id, err := client.Entries.Credential.New(Entry{
		VaultId: testVaultID,
		Name:    "cred",
		Path:    "A",
		Type:    EntryCredentialType,
		SubType: EntryCredentialSubTypeAccessCode,
		Tags:    []string{"t"},
		Data:    &EntryCredentialAccessCodeData{Password: "secret"},
	})
	require.NoError(t, err)
	assert.Equal(t, testEntryID, id)
}

func TestCredentialUpdate_SendsRequiredFieldsAndRefetches(t *testing.T) {
	gets := 0

	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPut:
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, map[string]any{
				"name":        "cred",
				"description": "",
				"path":        "",
				"tags":        []any{},
				"data":        map[string]any{"username": "u", "domain": ""},
			}, body, "an empty secret is left out so the stored one is kept; an empty domain is sent to clear it")
			w.Write([]byte(`{"id":"` + testEntryID + `","name":"cred","type":"Credential","subType":"Default","data":{"username":"u"}}`))
		case http.MethodGet:
			gets++
			w.Write([]byte(testCredentialJSON))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, mux)

	updated, err := client.Entries.Credential.Update(Entry{
		Id:      testEntryID,
		VaultId: testVaultID,
		Name:    "cred",
		Type:    EntryCredentialType,
		SubType: EntryCredentialSubTypeDefault,
		Data:    &EntryCredentialDefaultData{Username: "u"},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, gets)
	data, _ := updated.GetCredentialDefaultData()
	assert.Equal(t, "p", data.Password)
}

func TestCredentialDeleteById(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, mux)

	require.NoError(t, client.Entries.Credential.DeleteById(testVaultID, testEntryID))
}

func TestCredentialGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError error
	}{
		{name: "Found", data: `{"id":"` + testEntryID + `","name":"cred","type":"Credential","subType":"Default","data":{}},{"id":"2","name":"cred","type":"Credential","subType":"AccessCode","data":{}}`},
		{name: "NotFound", data: `{"id":"2","name":"cred","type":"Credential","subType":"AccessCode","data":{}}`, wantError: ErrEntryNotFound},
		{name: "Multiple", data: `{"id":"1","name":"cred","type":"Credential","subType":"Default","data":{}},{"id":"2","name":"cred","type":"Credential","subType":"Default","data":{}}`, wantError: ErrMultipleEntriesFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "cred", r.URL.Query().Get("name"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":2,"totalPage":1,"data":[` + tt.data + `]}`))
			})
			mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(testCredentialJSON))
			})

			client := newTestClient(t, mux)

			entry, err := client.Entries.Credential.GetByName(testVaultID, "cred", EntryCredentialSubTypeDefault, GetByNameOptions{})
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			data, _ := entry.GetCredentialDefaultData()
			assert.Equal(t, "p", data.Password, "the match is read back with its secrets")
		})
	}
}

func TestCredentialValidation(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	_, err := client.Entries.Credential.New(Entry{Name: "x", Type: EntryCredentialType, SubType: EntryCredentialSubTypeDefault})
	assert.ErrorContains(t, err, "VaultId")

	_, err = client.Entries.Credential.New(Entry{VaultId: testVaultID, Name: "x", Type: EntryCredentialType, SubType: "Nope"})
	assert.ErrorContains(t, err, "unsupported entry subtype")

	_, err = client.Entries.Credential.Update(Entry{VaultId: testVaultID, Name: "x", Type: EntryCredentialType, SubType: EntryCredentialSubTypeDefault})
	assert.ErrorContains(t, err, "entry Id is required")
}
