package dvls

//go:generate stringer -type=AdministrativeRoleScopeType -trimprefix AdministrativeRoleScope
type AdministrativeRoleScopeType int

const (
	AdministrativeRoleScopeGlobal AdministrativeRoleScopeType = iota
	AdministrativeRoleScopeOrganizationalUnit
	AdministrativeRoleScopeVault
	AdministrativeRoleScopeGateway
	AdministrativeRoleScopePamProvider
)

//go:generate stringer -type=AdministrativeRoleMemberAction -trimprefix AdministrativeRoleMemberAction
type AdministrativeRoleMemberAction int

const (
	AdministrativeRoleMemberActionAdd AdministrativeRoleMemberAction = iota
	AdministrativeRoleMemberActionDelete
)

//go:generate stringer -type=AdministrativeRoleAssigneeType -trimprefix AdministrativeRoleAssignee
type AdministrativeRoleAssigneeType int

const (
	AdministrativeRoleAssigneeUser AdministrativeRoleAssigneeType = iota
	AdministrativeRoleAssigneeApplication
	AdministrativeRoleAssigneeUserGroup
)

// Built-in administrative role IDs.
const (
	BuiltinRoleBuiltinAdministratorId           string = "00000000-0000-0000-0000-000000000001"
	BuiltinRoleUsersAdministratorId             string = "00000000-0000-0000-0000-000000000002"
	BuiltinRoleVaultsAdministratorId            string = "00000000-0000-0000-0000-000000000004"
	BuiltinRolePamAdministratorId               string = "00000000-0000-0000-0000-000000000005"
	BuiltinRoleGatewayAdministratorId           string = "00000000-0000-0000-0000-00000000000a"
	BuiltinRoleWorkspaceSettingsAdministratorId string = "00000000-0000-0000-0000-00000000000b"
	BuiltinRoleWorkspaceAdministratorId         string = "00000000-0000-0000-0000-00000000000d"
	BuiltinRoleVaultOwnerId                     string = "00000000-0000-0000-0000-00000000000e"
	BuiltinRoleVaultUserId                      string = "00000000-0000-0000-0000-00000000000f"
	BuiltinRolePamProviderAdministratorId       string = "00000000-0000-0000-0000-000000000011"
	BuiltinRoleUserGroupsAdministratorId        string = "00000000-0000-0000-0000-000000000013"
	BuiltinRoleLicensesAdministratorId          string = "00000000-0000-0000-0000-000000000014"
	BuiltinRoleEntryTemplatesAdministratorId    string = "00000000-0000-0000-0000-000000000015"
	BuiltinRolePasswordPoliciesAdministratorId  string = "00000000-0000-0000-0000-000000000016"
	BuiltinRoleSystemImagesAdministratorId      string = "00000000-0000-0000-0000-000000000017"
	BuiltinRolePamAccountCreatorId              string = "00000000-0000-0000-0000-000000000018"
	BuiltinRoleWorkspaceLogsViewerId            string = "00000000-0000-0000-0000-000000000019"
)
