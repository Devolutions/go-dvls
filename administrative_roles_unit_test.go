package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdministrativeRolesList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + BuiltinRoleVaultUserId + `","name":"Vault User","description":"Has access to the content of the vault","isBuiltIn":true,"permissions":[332],"supportedScopes":[0,2]},
			{"id":"` + testRoleID + `","name":"Custom Role","description":"","permissions":[100,101],"supportedScopes":[0],"isUsed":true}
		]}`))
	})

	client := newTestClient(t, mux)

	roles, err := client.AdministrativeRoles.List()
	require.NoError(t, err)
	require.Len(t, roles, 2)
	assert.Equal(t, BuiltinRoleVaultUserId, roles[0].Id)
	assert.True(t, roles[0].IsBuiltIn)
	assert.Equal(t, []AdministrativePermission{AdministrativePermissionRepositoriesContentView}, roles[0].Permissions)
	assert.Equal(t, []AdministrativeRoleScopeType{AdministrativeRoleScopeGlobal, AdministrativeRoleScopeVault}, roles[0].SupportedScopes)
	assert.Equal(t, []AdministrativePermission{AdministrativePermissionUsersView, AdministrativePermissionUsersAdd}, roles[1].Permissions)
	assert.True(t, roles[1].IsUsed)
}

func TestAdministrativeRolesGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles/"+testRoleID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{"id":"` + testRoleID + `","name":"Custom Role","description":"desc","permissions":[332],"supportedScopes":[2],"requiredPermissionOnScope":332}}`))
	})

	client := newTestClient(t, mux)

	role, err := client.AdministrativeRoles.Get(testRoleID)
	require.NoError(t, err)
	assert.Equal(t, "Custom Role", role.Name)
	require.NotNil(t, role.RequiredPermissionOnScope)
	assert.Equal(t, AdministrativePermissionRepositoriesContentView, *role.RequiredPermissionOnScope)
}

func TestAdministrativeRolesGet_NotFoundStatusCode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles/"+testRoleID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
	})

	client := newTestClient(t, mux)

	_, err := client.AdministrativeRoles.Get(testRoleID)
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)
}

func TestAdministrativeRolesGet_NotFoundResult(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles/"+testRoleID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":6}`))
	})

	client := newTestClient(t, mux)

	_, err := client.AdministrativeRoles.Get(testRoleID)
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)
}

func TestAdministrativeRolesGetByName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"role-1","name":"Alpha","description":"","permissions":[],"supportedScopes":[0]},
			{"id":"role-2","name":"Beta","description":"","permissions":[],"supportedScopes":[0]}
		]}`))
	})

	client := newTestClient(t, mux)

	role, err := client.AdministrativeRoles.GetByName("Beta")
	require.NoError(t, err)
	assert.Equal(t, "role-2", role.Id)
}

func TestAdministrativeRolesGetByName_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[]}`))
	})

	client := newTestClient(t, mux)

	_, err := client.AdministrativeRoles.GetByName("NonExistent")
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)
}

func TestAdministrativeRolesGetByName_Multiple(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-roles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"role-1","name":"Dup","description":"","permissions":[],"supportedScopes":[0]},
			{"id":"role-2","name":"Dup","description":"","permissions":[],"supportedScopes":[0]}
		]}`))
	})

	client := newTestClient(t, mux)

	_, err := client.AdministrativeRoles.GetByName("Dup")
	assert.ErrorIs(t, err, ErrMultipleAdministrativeRolesFound)
}

func TestAdministrativeEnumValues(t *testing.T) {
	assert.EqualValues(t, 2, AdministrativeRoleScopeVault)
	assert.EqualValues(t, 332, AdministrativePermissionRepositoriesContentView)
	assert.EqualValues(t, 205, AdministrativePermissionAdministrativeRolesAssignmentsManage)
}
