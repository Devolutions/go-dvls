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

	security, err := testClient.Entries.Permissions.Get(folderId)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideDefault, security.RoleOverride)

	err = testClient.Entries.Permissions.Set(folderId, EntrySecurity{
		RoleOverride: SecurityRoleOverrideEveryone,
		ViewOverride: SecurityRoleOverrideEveryone,
	})
	require.NoError(t, err)

	security, err = testClient.Entries.Permissions.Get(folderId)
	require.NoError(t, err)
	assert.Equal(t, SecurityRoleOverrideEveryone, security.RoleOverride)

	folder, err := testClient.Entries.Folder.GetById(vault.Id, folderId)
	require.NoError(t, err)
	assert.Equal(t, "Permissions Folder", folder.Name)
}
