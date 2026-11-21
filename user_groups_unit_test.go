package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserGroupJSON = `{"id":"00000000-0000-0000-eeee-000000000001","name":"Everyone","description":"","isEnabled":true,"isWorkspaceOwner":false,"type":"Custom"}`

func TestUserGroupsListAndGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user-group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testUserGroupJSON + `]}`))
	})
	mux.HandleFunc("/api/v2/user-group/00000000-0000-0000-eeee-000000000001", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testUserGroupJSON))
	})

	client := newTestClient(t, mux)

	want := UserGroup{Id: "00000000-0000-0000-eeee-000000000001", Name: "Everyone", Type: UserGroupTypeCustom, IsEnabled: true}

	groups, err := client.UserGroups.List()
	require.NoError(t, err)
	assert.Equal(t, []UserGroup{want}, groups)

	group, err := client.UserGroups.Get(want.Id)
	require.NoError(t, err)
	assert.Equal(t, want, group)
}

func TestUserGroupsGetByName(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError error
	}{
		{name: "Found", data: testUserGroupJSON},
		{name: "NotFound", data: "", wantError: ErrUserGroupNotFound},
		{name: "Multiple", data: testUserGroupJSON + "," + testUserGroupJSON, wantError: ErrMultipleUserGroupsFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/user-group", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "Everyone", r.URL.Query().Get("name"))
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + tt.data + `]}`))
			})

			client := newTestClient(t, mux)

			group, err := client.UserGroups.GetByName("Everyone")
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "Everyone", group.Name)
		})
	}
}
