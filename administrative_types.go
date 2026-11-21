package dvls

// AdministrativePermission is the name of a permission a role grants (e.g. VaultsContentView).
type AdministrativePermission string

// AdministrativeRoleScopeType is the scope of a role assignment.
type AdministrativeRoleScopeType string

const (
	AdministrativeRoleScopeWorkspace    AdministrativeRoleScopeType = "Workspace"
	AdministrativeRoleScopeVault        AdministrativeRoleScopeType = "Vault"
	AdministrativeRoleScopeBusinessUnit AdministrativeRoleScopeType = "BusinessUnit"
)

// AdministrativeRoleAssigneeType is the kind of principal holding a role assignment.
type AdministrativeRoleAssigneeType string

const (
	AdministrativeRoleAssigneeUser                AdministrativeRoleAssigneeType = "User"
	AdministrativeRoleAssigneeUserGroup           AdministrativeRoleAssigneeType = "UserGroup"
	AdministrativeRoleAssigneeApplicationIdentity AdministrativeRoleAssigneeType = "ApplicationIdentity"
)

// Built-in administrative role IDs.
const (
	BuiltinRoleBuiltinAdministratorId          string = "00000000-0000-0000-0000-000000000001"
	BuiltinRoleUsersAdministratorId            string = "00000000-0000-0000-0000-000000000002"
	BuiltinRoleVaultsAdministratorId           string = "00000000-0000-0000-0000-000000000004"
	BuiltinRoleWorkspaceAdministratorId        string = "00000000-0000-0000-0000-00000000000d"
	BuiltinRoleVaultOwnerId                    string = "00000000-0000-0000-0000-00000000000e"
	BuiltinRoleVaultUserId                     string = "00000000-0000-0000-0000-00000000000f"
	BuiltinRoleLicensesAdministratorId         string = "00000000-0000-0000-0000-000000000014"
	BuiltinRoleEntryTemplatesAdministratorId   string = "00000000-0000-0000-0000-000000000015"
	BuiltinRolePasswordPoliciesAdministratorId string = "00000000-0000-0000-0000-000000000016"
	BuiltinRoleSystemImagesAdministratorId     string = "00000000-0000-0000-0000-000000000017"
	BuiltinRoleWorkspaceLogsViewerId           string = "00000000-0000-0000-0000-000000000019"
)
