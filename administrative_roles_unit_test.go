package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRoleJSON = `{"id":"` + testRoleID + `","name":"Vault member","description":"d","isBuiltIn":true,"assignableScopes":["Workspace","Vault","BusinessUnit"],"permissions":["VaultsContentView"]}`

func TestAdministrativeRolesList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/role", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testRoleJSON + `]}`))
	})

	client := newTestClient(t, mux)

	roles, err := client.AdministrativeRoles.List()
	require.NoError(t, err)
	require.Len(t, roles, 1)
	assert.Equal(t, AdministrativeRole{
		Id:               testRoleID,
		Name:             "Vault member",
		Description:      "d",
		IsBuiltIn:        true,
		AssignableScopes: []AdministrativeRoleScopeType{AdministrativeRoleScopeWorkspace, AdministrativeRoleScopeVault, AdministrativeRoleScopeBusinessUnit},
		Permissions:      []AdministrativePermission{"VaultsContentView"},
	}, roles[0])
}

func TestAdministrativeRolesGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/role/"+testRoleID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testRoleJSON))
	})
	mux.HandleFunc("/api/v2/role/missing", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"The requested resource was not found."}`))
	})

	client := newTestClient(t, mux)

	role, err := client.AdministrativeRoles.Get(testRoleID)
	require.NoError(t, err)
	assert.Equal(t, "Vault member", role.Name)

	_, err = client.AdministrativeRoles.Get("missing")
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)
}

func TestAdministrativeRolesGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError error
	}{
		{name: "Found", data: testRoleJSON},
		{name: "NotFound", data: "", wantError: ErrAdministrativeRoleNotFound},
		{name: "Multiple", data: testRoleJSON + "," + testRoleJSON, wantError: ErrMultipleAdministrativeRolesFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/role", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "Vault member", r.URL.Query().Get("name"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":2,"totalPage":1,"data":[` + tt.data + `]}`))
			})

			client := newTestClient(t, mux)

			role, err := client.AdministrativeRoles.GetByName("Vault member")
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testRoleID, role.Id)
		})
	}
}
