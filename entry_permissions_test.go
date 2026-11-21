//go:build integration

package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EntryPermissions(t *testing.T) {
	vault := createTestVault(t, "entry-permissions")

	folderId, err := testClient.Entries.Folder.New(Entry{
		VaultId: vault.Id,
		Name:    "Permissions Folder",
		Type:    EntryFolderType,
		SubType: EntryFolderSubTypeFolder,
		Data:    &EntryFolderData{},
	})
	require.NoError(t, err)

	security, err := testClient.Entries.Permissions.Get(vault.Id, folderId)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideDefault, security.RoleOverride)

	users, err := testClient.Users.List()
	require.NoError(t, err)
	require.NotEmpty(t, users)

	custom := EntrySecurity{
		RoleOverride: SecurityRoleOverrideCustom,
		Permissions: []EntryPermission{
			{Override: SecurityRoleOverrideCustom, Right: SecurityRoleRightView, Roles: []string{users[0].Id}},
			{Override: SecurityRoleOverrideNever, Right: SecurityRoleRightDelete, Roles: []string{}},
		},
	}
	err = testClient.Entries.Permissions.Set(vault.Id, folderId, custom)
	require.NoError(t, err)

	security, err = testClient.Entries.Permissions.Get(vault.Id, folderId)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideCustom, security.RoleOverride)
	assert.ElementsMatch(t, custom.Permissions, security.Permissions)

	err = testClient.Entries.Permissions.Set(vault.Id, folderId, EntrySecurity{})
	require.NoError(t, err)

	security, err = testClient.Entries.Permissions.Get(vault.Id, folderId)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideDefault, security.RoleOverride)
	assert.Empty(t, security.Permissions)

	_, err = testClient.Entries.Permissions.Get(vault.Id, "00000000-0000-0000-0000-0000000000ff")
	assert.ErrorIs(t, err, ErrEntryNotFound)

	folder, err := testClient.Entries.Folder.GetById(vault.Id, folderId)
	require.NoError(t, err)
	assert.Equal(t, "Permissions Folder", folder.Name)
}
