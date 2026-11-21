package dvls

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryPermissionsGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath+"/permission", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"mode":"Custom","permissions":[{"mode":"Custom","right":"View","assigneeIds":["` + testUserID + `"]},{"mode":"Never","right":"Delete"}]}`))
	})

	client := newTestClient(t, mux)

	security, err := client.Entries.Permissions.Get(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, EntrySecurity{
		RoleOverride: SecurityRoleOverrideCustom,
		Permissions: []EntryPermission{
			{Override: SecurityRoleOverrideCustom, Right: SecurityRoleRightView, Roles: []string{testUserID}},
			{Override: SecurityRoleOverrideNever, Right: SecurityRoleRightDelete},
		},
	}, security)
}

func TestEntryPermissionsGet_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath+"/permission", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"The requested resource was not found."}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Entries.Permissions.Get(testVaultID, testEntryID)
	assert.ErrorIs(t, err, ErrEntryNotFound)
}

func TestEntryPermissionsSet(t *testing.T) {
	tests := []struct {
		name     string
		security EntrySecurity
		wantBody string
	}{
		{
			name: "Custom",
			security: EntrySecurity{
				RoleOverride: SecurityRoleOverrideCustomInherited,
				Permissions:  []EntryPermission{{Override: SecurityRoleOverrideCustom, Right: SecurityRoleRightEdit, Roles: []string{testUserID}}},
			},
			wantBody: `{"mode":"CustomAndInherited","permissions":[{"mode":"Custom","right":"Edit","assigneeIds":["` + testUserID + `"]}]}`,
		},
		{
			name:     "ZeroValueResets",
			security: EntrySecurity{},
			wantBody: `{"mode":"Default","permissions":[]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testEntryPath+"/permission", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				var body json.RawMessage
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.JSONEq(t, tt.wantBody, string(body))
				w.Header().Set("Content-Type", "application/json")
				w.Write(body)
			})

			client := newTestClient(t, mux)

			require.NoError(t, client.Entries.Permissions.Set(testVaultID, testEntryID, tt.security))
		})
	}
}

func TestEntryPermissionsSet_Validation(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	err := client.Entries.Permissions.Set("", testEntryID, EntrySecurity{})
	assert.ErrorContains(t, err, "vault Id are required")
}
