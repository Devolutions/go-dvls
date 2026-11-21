//go:build integration

package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AdministrativeRoles(t *testing.T) {
	roles, err := testClient.AdministrativeRoles.List()
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	var vaultUser AdministrativeRole
	for _, role := range roles {
		if role.Id == BuiltinRoleVaultUserId {
			vaultUser = role
		}
	}
	require.NotEmpty(t, vaultUser.Id, "expected built-in Vault member role in list")
	assert.True(t, vaultUser.IsBuiltIn)
	assert.Contains(t, vaultUser.Permissions, AdministrativePermission("VaultsContentView"))
	assert.Contains(t, vaultUser.AssignableScopes, AdministrativeRoleScopeVault)

	fetched, err := testClient.AdministrativeRoles.Get(BuiltinRoleVaultUserId)
	require.NoError(t, err)
	assert.Equal(t, vaultUser.Name, fetched.Name)

	byName, err := testClient.AdministrativeRoles.GetByName(vaultUser.Name)
	require.NoError(t, err)
	assert.Equal(t, BuiltinRoleVaultUserId, byName.Id)

	_, err = testClient.AdministrativeRoles.Get("00000000-0000-0000-0000-0000000000ff")
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)

	_, err = testClient.AdministrativeRoles.GetByName("go-dvls-nonexistent-role")
	assert.ErrorIs(t, err, ErrAdministrativeRoleNotFound)
}

func Test_AdministrativeRoleAssignments(t *testing.T) {
	vault := createTestVault(t, "role-assignments")

	owners, err := testClient.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{
		RoleId:    BuiltinRoleBuiltinAdministratorId,
		ScopeType: AdministrativeRoleScopeWorkspace,
	})
	require.NoError(t, err)
	require.NotEmpty(t, owners, "expected at least one workspace owner")

	groups, err := testClient.UserGroups.List()
	require.NoError(t, err)
	require.NotEmpty(t, groups)
	assigneeId := groups[0].Id

	created, err := testClient.AdministrativeRoleAssignments.New(AdministrativeRoleAssignmentRequest{
		RoleId:     BuiltinRoleVaultUserId,
		AssigneeId: assigneeId,
		ScopeType:  AdministrativeRoleScopeVault,
		ScopeId:    vault.Id,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = testClient.AdministrativeRoleAssignments.Delete(created.Id) })
	assert.Equal(t, AdministrativeRoleAssigneeUserGroup, created.AssigneeType)
	assert.Equal(t, vault.Id, created.ScopeId)

	fetched, err := testClient.AdministrativeRoleAssignments.Get(created.Id)
	require.NoError(t, err)
	assert.Equal(t, created, fetched)

	_, err = testClient.AdministrativeRoleAssignments.New(AdministrativeRoleAssignmentRequest{
		RoleId:     BuiltinRoleVaultUserId,
		AssigneeId: assigneeId,
		ScopeType:  AdministrativeRoleScopeVault,
		ScopeId:    vault.Id,
	})
	assert.ErrorIs(t, err, ErrAdministrativeRoleAssignmentAlreadyExists)

	assignments, err := testClient.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{
		RoleId:    BuiltinRoleVaultUserId,
		ScopeType: AdministrativeRoleScopeVault,
		ScopeId:   vault.Id,
	})
	require.NoError(t, err)
	assert.Contains(t, assignments, created)

	byAssignee, err := testClient.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{AssigneeId: assigneeId})
	require.NoError(t, err)
	assert.Contains(t, byAssignee, created)

	err = testClient.AdministrativeRoleAssignments.Delete(created.Id)
	require.NoError(t, err)

	_, err = testClient.AdministrativeRoleAssignments.Get(created.Id)
	assert.True(t, IsNotFound(err))
}
