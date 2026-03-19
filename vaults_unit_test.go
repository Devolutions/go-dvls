package dvls

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultsList_Pagination(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")

		switch page {
		case "1":
			json.NewEncoder(w).Encode(vaultListResponse{
				Data:        []Vault{{Id: "vault-1", Name: "Page1Vault"}},
				CurrentPage: 1,
				PageSize:    1,
				TotalCount:  2,
				TotalPage:   2,
			})
		case "2":
			json.NewEncoder(w).Encode(vaultListResponse{
				Data:        []Vault{{Id: "vault-2", Name: "Page2Vault"}},
				CurrentPage: 2,
				PageSize:    1,
				TotalCount:  2,
				TotalPage:   2,
			})
		}
	})

	client := newTestClient(t, mux)

	vaults, err := client.Vaults.List()
	require.NoError(t, err)
	require.Len(t, vaults, 2)
	assert.Equal(t, "vault-1", vaults[0].Id)
	assert.Equal(t, "vault-2", vaults[1].Id)
}

func TestVaultsGetByName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vaultListResponse{
			Data: []Vault{
				{Id: "vault-1", Name: "Alpha"},
				{Id: "vault-2", Name: "Beta"},
				{Id: "vault-3", Name: "Gamma"},
			},
			CurrentPage: 1,
			TotalPage:   1,
		})
	})

	client := newTestClient(t, mux)

	vault, err := client.Vaults.GetByName("Beta")
	require.NoError(t, err)
	assert.Equal(t, "vault-2", vault.Id)
	assert.Equal(t, "Beta", vault.Name)
}

func TestVaultsGetByName_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vaultListResponse{
			Data:        []Vault{{Id: "vault-1", Name: "Alpha"}},
			CurrentPage: 1,
			TotalPage:   1,
		})
	})

	client := newTestClient(t, mux)

	_, err := client.Vaults.GetByName("NonExistent")
	assert.ErrorIs(t, err, ErrVaultNotFound)
}

func TestVaultsNew_DefaultToEverything(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vault", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req vaultRequest
		json.Unmarshal(body, &req)
		assert.Equal(t, VaultContentTypeEverything, req.ContentType)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Vault{
			Id:          "vault-id",
			Name:        req.Name,
			ContentType: VaultContentTypeEverything,
		})
	})

	client := newTestClient(t, mux)

	created, err := client.Vaults.New(Vault{
		Name:        "DefaultVault",
		ContentType: VaultContentTypeDefault,
	})
	require.NoError(t, err)
	assert.Equal(t, VaultContentTypeEverything, created.ContentType)
}
