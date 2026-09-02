package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/users/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + testUserID + `","name":"admin","fullName":"Admin User","email":"admin@example.com","authenticationType":0,"isAdministrator":true,"isEnabled":true,"userGroups":["Ops"]},
			{"id":"` + testAssigneeID + `","name":"jdoe","fullName":"John Doe","authenticationType":8,"isEnabled":true}
		]}`))
	})

	client := newTestClient(t, mux)

	users, err := client.Users.List()
	require.NoError(t, err)
	require.Len(t, users, 2)
	assert.Equal(t, testUserID, users[0].Id)
	assert.Equal(t, "admin", users[0].Name)
	assert.Equal(t, UserAuthenticationBuiltin, users[0].AuthenticationType)
	assert.True(t, users[0].IsAdministrator)
	assert.Equal(t, []string{"Ops"}, users[0].UserGroups)
	assert.Equal(t, UserAuthenticationAzureAD, users[1].AuthenticationType)
}

func TestUsersList_ResultError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/users/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":2,"message":"AccessDenied"}`))
	})

	client := newTestClient(t, mux)

	_, err := client.Users.List()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AccessDenied")
}

func TestUsersGetByName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/users/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + testUserID + `","name":"admin"},
			{"id":"` + testAssigneeID + `","name":"dup"},
			{"id":"` + testRoleID + `","name":"dup"}
		]}`))
	})

	client := newTestClient(t, mux)

	user, err := client.Users.GetByName("admin")
	require.NoError(t, err)
	assert.Equal(t, testUserID, user.Id)

	_, err = client.Users.GetByName("dup")
	assert.ErrorIs(t, err, ErrMultipleUsersFound)

	_, err = client.Users.GetByName("nobody")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUsersApplications(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/security/application/users/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"` + testAssigneeID + `","name":"my-app-key","fullName":"CI","authenticationType":9,"isEnabled":true}
		]}`))
	})

	client := newTestClient(t, mux)

	apps, err := client.Users.ListApplications()
	require.NoError(t, err)
	require.Len(t, apps, 1)
	assert.Equal(t, testAssigneeID, apps[0].Id)
	assert.Equal(t, UserAuthenticationApplication, apps[0].AuthenticationType)

	app, err := client.Users.GetApplicationByName("my-app-key")
	require.NoError(t, err)
	assert.Equal(t, testAssigneeID, app.Id)

	_, err = client.Users.GetApplicationByName("missing")
	assert.ErrorIs(t, err, ErrUserNotFound)
}
