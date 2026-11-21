package dvls

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testVaultJSON = `{"id":"` + testVaultID + `","name":"Prod","description":"d","contentType":"Everything","securityLevel":"Standard","type":"Shared","visibility":"Default"}`

func TestVaultsList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Empty(t, r.URL.Query().Get("name"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testVaultJSON + `]}`))
	})

	client := newTestClient(t, mux)

	vaults, err := client.Vaults.List()
	require.NoError(t, err)
	require.Len(t, vaults, 1)
	assert.Equal(t, Vault{Id: testVaultID, Name: "Prod", Description: "d", ContentType: VaultContentTypeEverything, Type: "Shared", SecurityLevel: VaultSecurityLevelStandard, Visibility: VaultVisibilityDefault}, vaults[0])
}

func TestVaultsGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault/"+testVaultID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testVaultJSON))
	})

	client := newTestClient(t, mux)

	vault, err := client.Vaults.Get(testVaultID)
	require.NoError(t, err)
	assert.Equal(t, "Prod", vault.Name)
}

func TestVaultsGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantId    string
		wantError error
	}{
		{name: "Found", data: testVaultJSON, wantId: testVaultID},
		{name: "NotFound", data: "", wantError: ErrVaultNotFound},
		{name: "CaseOnlyDuplicates", data: `{"id":"a","name":"prod"},{"id":"b","name":"Prod"}`, wantId: "b"},
		{name: "Multiple", data: `{"id":"a","name":"Prod"},{"id":"b","name":"Prod"}`, wantError: ErrMultipleVaultsFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/vault", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "Prod", r.URL.Query().Get("name"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":9,"totalPage":1,"data":[` + tt.data + `]}`))
			})

			client := newTestClient(t, mux)

			vault, err := client.Vaults.GetByName("Prod")
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantId, vault.Id)
		})
	}
}

func TestVaultsNew_OmitsUnsetSettings(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"name": "Prod", "contentType": "Default"}, body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(testVaultJSON))
	})

	client := newTestClient(t, mux)

	vault, err := client.Vaults.New(Vault{Name: "Prod", ContentType: VaultContentTypeDefault})
	require.NoError(t, err)
	assert.Equal(t, testVaultID, vault.Id)
}

func TestVaultsUpdate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/vault/"+testVaultID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, map[string]any{"name": "Prod", "description": "", "contentType": "Secrets", "securityLevel": "High", "visibility": "Never"}, body)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testVaultJSON))
	})

	client := newTestClient(t, mux)

	_, err := client.Vaults.Update(Vault{Id: testVaultID, Name: "Prod", ContentType: VaultContentTypeSecrets, SecurityLevel: VaultSecurityLevelHigh, Visibility: VaultVisibilityPrivate})
	require.NoError(t, err)
}
