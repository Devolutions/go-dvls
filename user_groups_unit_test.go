package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGroupsList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/roles/basic", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + testRoleID + `","name":"Ops","description":"Operations","isAdministrator":true,"roleType":1},
			{"id":"` + testAssigneeID + `","name":"Domain Admins","description":"","roleType":0}
		]}`))
	})

	client := newTestClient(t, mux)

	groups, err := client.UserGroups.List()
	require.NoError(t, err)
	require.Len(t, groups, 2)
	assert.Equal(t, testRoleID, groups[0].Id)
	assert.Equal(t, "Ops", groups[0].Name)
	assert.True(t, groups[0].IsAdministrator)
	assert.Equal(t, UserGroupTypeCustom, groups[0].Type)
	assert.Equal(t, UserGroupTypeActiveDirectory, groups[1].Type)
}

func TestUserGroupsList_ResultError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/roles/basic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":2,"message":"AccessDenied"}`))
	})

	client := newTestClient(t, mux)

	_, err := client.UserGroups.List()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AccessDenied")
}

func TestUserGroupsGetByName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/roles/basic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + testRoleID + `","name":"Ops"},
			{"id":"` + testAssigneeID + `","name":"Dup"},
			{"id":"` + testUserID + `","name":"Dup"}
		]}`))
	})

	client := newTestClient(t, mux)

	group, err := client.UserGroups.GetByName("Ops")
	require.NoError(t, err)
	assert.Equal(t, testRoleID, group.Id)

	_, err = client.UserGroups.GetByName("Dup")
	assert.ErrorIs(t, err, ErrMultipleUserGroupsFound)

	_, err = client.UserGroups.GetByName("nope")
	assert.ErrorIs(t, err, ErrUserGroupNotFound)
}
