package dvls

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryPermissionsGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/connections/partial/"+testEntryID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{"id":"` + testEntryID + `","name":"Entry","connectionType":26,
			"security":{"roleOverride":1,"viewOverride":5,"viewRoles":["role-1"],"permissions":[{"override":1,"right":4,"roles":["role-2"]}]}}}`))
	})

	client := newTestClient(t, mux)

	security, err := client.Entries.Permissions.Get(testEntryID)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideCustom, security.RoleOverride)
	assert.Equal(t, SecurityRoleOverrideCustomInherited, security.ViewOverride)
	assert.Equal(t, []string{"role-1"}, security.ViewRoles)
	require.Len(t, security.Permissions, 1)
	assert.Equal(t, SecurityRoleRightEdit, security.Permissions[0].Right)
	assert.Equal(t, SecurityRoleOverrideCustom, security.Permissions[0].Override)
	assert.Equal(t, []string{"role-2"}, security.Permissions[0].Roles)
}

func TestEntryPermissionsGet_NoSecurityNode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/connections/partial/"+testEntryID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{"id":"` + testEntryID + `","name":"Entry","connectionType":26}}`))
	})

	client := newTestClient(t, mux)

	security, err := client.Entries.Permissions.Get(testEntryID)
	require.NoError(t, err)
	assert.Equal(t, EntrySecurity{}, security)
}

func TestEntryPermissionsSet_RoundTripPreservesUnknownFields(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/connections/partial/"+testEntryID, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{
			"id":"` + testEntryID + `","name":"Entry","connectionType":26,"repositoryId":"` + testVaultID + `",
			"events":{"openCommentPrompt":true},"data":"{\"nested\":true}","futureField":"abc",
			"security":{"roleOverride":0,"checkOutMode":2,"allowOffline":1,"passwordComplexityId":"pc-1"}}}`))
	})

	var savedBody []byte
	mux.HandleFunc("/api/connections/partial/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		savedBody = body

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1}`))
	})

	client := newTestClient(t, mux)

	err := client.Entries.Permissions.Set(testEntryID, EntrySecurity{
		RoleOverride: SecurityRoleOverrideCustom,
		ViewOverride: SecurityRoleOverrideCustom,
		ViewRoles:    []string{"role-1"},
		Permissions: []EntryPermission{
			{Override: SecurityRoleOverrideCustom, Right: SecurityRoleRightEdit, Roles: []string{"role-2"}},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, savedBody)

	var saved map[string]any
	require.NoError(t, json.Unmarshal(savedBody, &saved))

	assert.NotContains(t, saved, "result")
	assert.Equal(t, testEntryID, saved["id"])
	assert.Equal(t, map[string]any{"openCommentPrompt": true}, saved["events"])
	assert.Equal(t, `{"nested":true}`, saved["data"])
	assert.Equal(t, "abc", saved["futureField"])
	assert.Equal(t, testVaultID, saved["repositoryId"])

	security, ok := saved["security"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(2), security["checkOutMode"])
	assert.Equal(t, float64(1), security["allowOffline"])
	assert.Equal(t, "pc-1", security["passwordComplexityId"])
	assert.Equal(t, float64(1), security["roleOverride"])
	assert.Equal(t, float64(1), security["viewOverride"])
	assert.Equal(t, []any{"role-1"}, security["viewRoles"])
	permissions, ok := security["permissions"].([]any)
	require.True(t, ok)
	require.Len(t, permissions, 1)
	assert.Equal(t, map[string]any{"override": float64(1), "right": float64(4), "roles": []any{"role-2"}}, permissions[0])
}

func TestEntryPermissionsSet_FolderGroupTrimmed(t *testing.T) {
	cases := []struct {
		name          string
		entryName     string
		group         string
		expectedGroup string
	}{
		{name: "NestedFolder", entryName: "Sub", group: `Parent\Sub`, expectedGroup: "Parent"},
		{name: "RootFolder", entryName: "Root", group: "Root", expectedGroup: ""},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/connections/partial/"+testEntryID, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				entryJson, err := json.Marshal(map[string]any{
					"id":             testEntryID,
					"name":           testCase.entryName,
					"group":          testCase.group,
					"connectionType": 25,
				})
				require.NoError(t, err)
				w.Write([]byte(`{"result":1,"data":` + string(entryJson) + `}`))
			})

			mux.HandleFunc("/api/connections/partial/save", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var saved map[string]any
				require.NoError(t, json.Unmarshal(body, &saved))
				assert.Equal(t, testCase.expectedGroup, saved["group"])

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"result":1}`))
			})

			client := newTestClient(t, mux)

			err := client.Entries.Permissions.Set(testEntryID, EntrySecurity{RoleOverride: SecurityRoleOverrideEveryone, ViewOverride: SecurityRoleOverrideEveryone})
			require.NoError(t, err)
		})
	}
}

func TestEntryPermissionsSet_ValidationErrors(t *testing.T) {
	client := newTestClient(t, http.NewServeMux())

	err := client.Entries.Permissions.Set("", EntrySecurity{})
	assert.ErrorContains(t, err, "entry id is required")

	err = client.Entries.Permissions.Set(testEntryID, EntrySecurity{
		RoleOverride: SecurityRoleOverrideCustom,
		Permissions:  []EntryPermission{{Right: SecurityRoleRightView}},
	})
	assert.ErrorContains(t, err, "View right")

	err = client.Entries.Permissions.Set(testEntryID, EntrySecurity{
		RoleOverride: SecurityRoleOverrideCustom,
		Permissions: []EntryPermission{
			{Right: SecurityRoleRightEdit},
			{Right: SecurityRoleRightEdit},
		},
	})
	assert.ErrorContains(t, err, "duplicate permission right")

	err = client.Entries.Permissions.Set(testEntryID, EntrySecurity{
		RoleOverride: SecurityRoleOverrideDefault,
		Permissions:  []EntryPermission{{Right: SecurityRoleRightEdit}},
	})
	assert.ErrorContains(t, err, "RoleOverride")
}

func TestEntryPermissionsSet_SaveResultError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/connections/partial/"+testEntryID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{"id":"` + testEntryID + `","name":"Entry","connectionType":26}}`))
	})
	mux.HandleFunc("/api/connections/partial/save", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":2,"message":"access denied"}`))
	})

	client := newTestClient(t, mux)

	err := client.Entries.Permissions.Set(testEntryID, EntrySecurity{RoleOverride: SecurityRoleOverrideEveryone, ViewOverride: SecurityRoleOverrideEveryone})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AccessDenied")
}

func TestSecurityRoleEnumValues(t *testing.T) {
	assert.EqualValues(t, 0, SecurityRoleRightView)
	assert.EqualValues(t, 4, SecurityRoleRightEdit)
	assert.EqualValues(t, 22, SecurityRoleRightExecute)
	assert.EqualValues(t, 37, SecurityRoleRightEditSessionRecordingConfiguration)
	assert.EqualValues(t, 5, SecurityRoleOverrideCustomInherited)
}
