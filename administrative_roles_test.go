//go:build integration

package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requireRoleAssignmentSupport skips the test when the server does not expose
// the administrative role endpoints (DVLS 2026.3+), and returns the roles
// fetched by the probe.
func requireRoleAssignmentSupport(t *testing.T) []AdministrativeRole {
	t.Helper()

	roles, err := testClient.AdministrativeRoles.List()
	if IsNotFound(err) {
		t.Skip("administrative roles require DVLS 2026.3 or later")
	}
	require.NoError(t, err)

	return roles
}

func Test_AdministrativeRoles(t *testing.T) {
	roles := requireRoleAssignmentSupport(t)
	require.NotEmpty(t, roles)

	var vaultUser AdministrativeRole
	for _, role := range roles {
		if role.Id == BuiltinRoleVaultUserId {
			vaultUser = role
		}
	}
	require.NotEmpty(t, vaultUser.Id, "expected built-in Vault User role in list")
	assert.True(t, vaultUser.IsBuiltIn)
	assert.Contains(t, vaultUser.Permissions, AdministrativePermissionRepositoriesContentView)
	assert.Contains(t, vaultUser.SupportedScopes, AdministrativeRoleScopeVault)

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
	requireRoleAssignmentSupport(t)

	vault := createTestVault(t, "role-assignments")

	admins, err := testClient.AdministrativeRoleAssignments.GetMembers(BuiltinRoleBuiltinAdministratorId, AdministrativeRoleScopeGlobal, "")
	require.NoError(t, err)
	require.NotEmpty(t, admins, "expected at least one Built-in administrator member")
	assigneeId := admins[0].AssigneeId

	scopeResourceId := vault.Id
	err = testClient.AdministrativeRoleAssignments.AddMember(AdministrativeRoleMemberRequest{
		AdministrativeRoleId: BuiltinRoleVaultUserId,
		AssigneeId:           assigneeId,
		ScopeType:            AdministrativeRoleScopeVault,
		ScopeResourceId:      &scopeResourceId,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		testClient.AdministrativeRoleAssignments.DeleteScope(BuiltinRoleVaultUserId, AdministrativeRoleScopeVault, vault.Id)
	})

	members, err := testClient.AdministrativeRoleAssignments.GetMembers(BuiltinRoleVaultUserId, AdministrativeRoleScopeVault, vault.Id)
	require.NoError(t, err)
	var foundMember bool
	for _, member := range members {
		if member.AssigneeId == assigneeId {
			foundMember = true
		}
	}
	assert.True(t, foundMember, "expected the assignee in the Vault User members")

	scopeType := AdministrativeRoleScopeVault
	assignments, err := testClient.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{
		RoleId:    BuiltinRoleVaultUserId,
		ScopeType: &scopeType,
		ScopeId:   vault.Id,
	})
	require.NoError(t, err)
	require.Len(t, assignments, 1)
	assert.Equal(t, BuiltinRoleVaultUserId, assignments[0].AdministrativeRoleId)
	assert.Equal(t, vault.Id, assignments[0].ScopeResourceId)
	assert.True(t, assignments[0].IsBuiltIn)

	byAssignee, err := testClient.AdministrativeRoleAssignments.ListByAssignee(assigneeId, false)
	require.NoError(t, err)
	var foundAssignment bool
	for _, assignment := range byAssignee {
		if assignment.AdministrativeRoleId == BuiltinRoleVaultUserId && assignment.ScopeResourceId == vault.Id {
			foundAssignment = true
		}
	}
	assert.True(t, foundAssignment, "expected the created assignment in the assignee's assignments")

	err = testClient.AdministrativeRoleAssignments.DeleteScope(BuiltinRoleVaultUserId, AdministrativeRoleScopeVault, vault.Id)
	require.NoError(t, err)

	members, err = testClient.AdministrativeRoleAssignments.GetMembers(BuiltinRoleVaultUserId, AdministrativeRoleScopeVault, vault.Id)
	require.NoError(t, err)
	assert.Empty(t, members)
}
