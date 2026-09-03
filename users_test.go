//go:build integration

package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Users(t *testing.T) {
	users, err := testClient.Users.List()
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, u := range users {
		assert.NotEmpty(t, u.Id)
		assert.NotEmpty(t, u.Name)
	}

	user, err := testClient.Users.GetByName(users[0].Name)
	require.NoError(t, err)
	assert.Equal(t, users[0].Id, user.Id)

	_, err = testClient.Users.GetByName("go-dvls-nonexistent-user")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func Test_Applications(t *testing.T) {
	apps, err := testClient.Users.ListApplications()
	require.NoError(t, err)

	require.NotEmpty(t, apps)

	for _, a := range apps {
		assert.NotEmpty(t, a.Id)
		assert.Equal(t, UserAuthenticationApplication, a.AuthenticationType)
	}

	app, err := testClient.Users.GetApplicationByName(apps[0].Name)
	require.NoError(t, err)
	assert.Equal(t, apps[0].Id, app.Id)

	_, err = testClient.Users.GetApplicationByName("go-dvls-nonexistent-app")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func Test_UserGroups(t *testing.T) {
	groups, err := testClient.UserGroups.List()
	require.NoError(t, err)
	require.NotEmpty(t, groups)

	for _, g := range groups {
		assert.NotEmpty(t, g.Id)
		assert.NotEmpty(t, g.Name)
	}

	group, err := testClient.UserGroups.GetByName(groups[0].Name)
	require.NoError(t, err)
	assert.Equal(t, groups[0].Id, group.Id)

	_, err = testClient.UserGroups.GetByName("go-dvls-nonexistent-group")
	assert.ErrorIs(t, err, ErrUserGroupNotFound)
}
