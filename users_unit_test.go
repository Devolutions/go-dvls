package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUserJSON        = `{"id":"` + testUserID + `","username":"jdoe","fullName":"John Doe","email":"j@x","identityProvider":"EntraId","isEnabled":true,"isWorkspaceOwner":false}`
	testApplicationJSON = `{"id":"` + testAssigneeID + `","name":"CI pipeline","appKey":"ci","isEnabled":true,"isWorkspaceOwner":true}`
)

func TestUsersListAndGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testUserJSON + `]}`))
	})
	mux.HandleFunc("/api/v2/user/"+testUserID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testUserJSON))
	})

	client := newTestClient(t, mux)

	want := User{Id: testUserID, Name: "jdoe", FullName: "John Doe", Email: "j@x", IdentityProvider: IdentityProviderEntraId, IsEnabled: true}

	users, err := client.Users.List()
	require.NoError(t, err)
	assert.Equal(t, []User{want}, users)

	user, err := client.Users.Get(testUserID)
	require.NoError(t, err)
	assert.Equal(t, want, user)
}

func TestUsersGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError error
	}{
		{name: "Found", data: testUserJSON},
		{name: "NotFound", data: "", wantError: ErrUserNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/user", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "jdoe", r.URL.Query().Get("username"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + tt.data + `]}`))
			})

			client := newTestClient(t, mux)

			user, err := client.Users.GetByName("jdoe")
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testUserID, user.Id)
		})
	}
}

func TestUsersApplications(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/application-identity", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Has("appKey") {
			assert.Equal(t, "ci", r.URL.Query().Get("appKey"))
		}
		assert.False(t, r.URL.Query().Has("name"), "an application identity is looked up by its unique appKey, not its display name")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testApplicationJSON + `]}`))
	})
	mux.HandleFunc("/api/v2/application-identity/"+testAssigneeID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testApplicationJSON))
	})

	client := newTestClient(t, mux)

	want := ApplicationIdentity{Id: testAssigneeID, Name: "CI pipeline", AppKey: "ci", IsEnabled: true, IsWorkspaceOwner: true}

	apps, err := client.Users.ListApplications()
	require.NoError(t, err)
	assert.Equal(t, []ApplicationIdentity{want}, apps)

	app, err := client.Users.GetApplication(testAssigneeID)
	require.NoError(t, err)
	assert.Equal(t, want, app)

	app, err = client.Users.GetApplicationByName("ci")
	require.NoError(t, err)
	assert.Equal(t, want, app)
}
