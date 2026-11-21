package dvls

import (
	"context"
	"fmt"
	"net/http"
)

type EntryPermissionsService service

const entryPermissionPath string = "permission"

// SecurityRoleOverride is the permission mode of an entry or of one of its rights.
type SecurityRoleOverride string

const (
	SecurityRoleOverrideDefault         SecurityRoleOverride = "Default"
	SecurityRoleOverrideCustom          SecurityRoleOverride = "Custom"
	SecurityRoleOverrideInherited       SecurityRoleOverride = "Inherited"
	SecurityRoleOverrideEveryone        SecurityRoleOverride = "Everyone"
	SecurityRoleOverrideNever           SecurityRoleOverride = "Never"
	SecurityRoleOverrideCustomInherited SecurityRoleOverride = "CustomAndInherited"
)

// SecurityRoleRight is a right an entry permission sets.
type SecurityRoleRight string

const (
	SecurityRoleRightView                              SecurityRoleRight = "View"
	SecurityRoleRightViewPassword                      SecurityRoleRight = "ViewPassword"
	SecurityRoleRightAdd                               SecurityRoleRight = "Add"
	SecurityRoleRightDelete                            SecurityRoleRight = "Delete"
	SecurityRoleRightEdit                              SecurityRoleRight = "Edit"
	SecurityRoleRightEditSecurity                      SecurityRoleRight = "EditSecurity"
	SecurityRoleRightPasswordHistory                   SecurityRoleRight = "PasswordHistory"
	SecurityRoleRightEntryHistory                      SecurityRoleRight = "EntryHistory"
	SecurityRoleRightRemoteTools                       SecurityRoleRight = "RemoteTools"
	SecurityRoleRightAttachment                        SecurityRoleRight = "Attachment"
	SecurityRoleRightEditAttachment                    SecurityRoleRight = "EditAttachment"
	SecurityRoleRightInventory                         SecurityRoleRight = "Inventory"
	SecurityRoleRightDocumentation                     SecurityRoleRight = "Documentation"
	SecurityRoleRightEditDocumentation                 SecurityRoleRight = "EditDocumentation"
	SecurityRoleRightWebManagementTools                SecurityRoleRight = "WebManagementTools"
	SecurityRoleRightConsoleManagementTools            SecurityRoleRight = "ConsoleManagementTools"
	SecurityRoleRightMacroScriptTools                  SecurityRoleRight = "MacroScriptTools"
	SecurityRoleRightMacroScriptToolsEntry             SecurityRoleRight = "MacroScriptToolsEntry"
	SecurityRoleRightExecute                           SecurityRoleRight = "Execute"
	SecurityRoleRightEditInformation                   SecurityRoleRight = "EditInformation"
	SecurityRoleRightMove                              SecurityRoleRight = "Move"
	SecurityRoleRightDeleteDocumentation               SecurityRoleRight = "DeleteDocumentation"
	SecurityRoleRightViewSensitiveInformation          SecurityRoleRight = "ViewSensitiveInformation"
	SecurityRoleRightResetPassword                     SecurityRoleRight = "ResetPassword"
	SecurityRoleRightForceCheckin                      SecurityRoleRight = "ForceCheckin"
	SecurityRoleRightCheckout                          SecurityRoleRight = "Checkout"
	SecurityRoleRightReadLogs                          SecurityRoleRight = "ReadLogs"
	SecurityRoleRightEditSessionRecordingConfiguration SecurityRoleRight = "EditSessionRecordingConfiguration"
	SecurityRoleRightEditStatus                        SecurityRoleRight = "EditStatus"
	SecurityRoleRightEditDescription                   SecurityRoleRight = "EditDescription"
	SecurityRoleRightViewLogs                          SecurityRoleRight = "ViewLogs"
	SecurityRoleRightEditPassword                      SecurityRoleRight = "EditPassword"
	SecurityRoleRightViewSessionRecording              SecurityRoleRight = "ViewSessionRecording"
	SecurityRoleRightViewInformation                   SecurityRoleRight = "ViewInformation"
	SecurityRoleRightExport                            SecurityRoleRight = "Export"
	SecurityRoleRightApproveCheckoutRequest            SecurityRoleRight = "ApproveCheckoutRequest"
	SecurityRoleRightEditVPNSSHGatewayConfiguration    SecurityRoleRight = "EditVpnSshGatewayConfiguration"
)

// EntryPermission represents who holds a single right on an entry or folder. Roles contains the ids of the users,
// user groups and application identities the right names.
type EntryPermission struct {
	Override SecurityRoleOverride `json:"mode"`
	Right    SecurityRoleRight    `json:"right"`
	Roles    []string             `json:"assigneeIds,omitempty"`
}

// EntrySecurity represents the permissions an entry or folder sets, as stored. RoleOverride is the permission mode
// of the entry; Permissions sets rights one by one, the View right included, each right at most once. A right that
// is not listed takes its default.
type EntrySecurity struct {
	RoleOverride SecurityRoleOverride `json:"mode"`
	Permissions  []EntryPermission    `json:"permissions"`
}

// Get returns the EntrySecurity of the entry or folder specified by vaultId and entryId.
// Returns ErrEntryNotFound if no entry is found.
func (c *EntryPermissionsService) Get(vaultId string, entryId string) (EntrySecurity, error) {
	return c.GetWithContext(context.Background(), vaultId, entryId)
}

// GetWithContext returns the EntrySecurity of the entry or folder specified by vaultId and entryId.
// Returns ErrEntryNotFound if no entry is found.
// The provided context can be used to cancel the request.
func (c *EntryPermissionsService) GetWithContext(ctx context.Context, vaultId string, entryId string) (EntrySecurity, error) {
	reqUrl, err := c.permissionUrl(vaultId, entryId)
	if err != nil {
		return EntrySecurity{}, err
	}

	var security EntrySecurity
	err = c.client.getJSON(ctx, reqUrl, &security)
	if err != nil {
		if IsNotFound(err) {
			return EntrySecurity{}, ErrEntryNotFound
		}

		return EntrySecurity{}, fmt.Errorf("error while fetching entry permissions: %w", err)
	}

	return security, nil
}

// Set replaces the permissions of the entry or folder specified by vaultId and entryId. The zero EntrySecurity
// (mode Default, no permissions) resets the entry to its defaults.
// Returns ErrEntryNotFound if no entry is found.
func (c *EntryPermissionsService) Set(vaultId string, entryId string, security EntrySecurity) error {
	return c.SetWithContext(context.Background(), vaultId, entryId, security)
}

// SetWithContext replaces the permissions of the entry or folder specified by vaultId and entryId. The zero
// EntrySecurity (mode Default, no permissions) resets the entry to its defaults.
// Returns ErrEntryNotFound if no entry is found.
// The provided context can be used to cancel the request.
func (c *EntryPermissionsService) SetWithContext(ctx context.Context, vaultId string, entryId string, security EntrySecurity) error {
	reqUrl, err := c.permissionUrl(vaultId, entryId)
	if err != nil {
		return err
	}

	if security.RoleOverride == "" {
		security.RoleOverride = SecurityRoleOverrideDefault
	}
	if security.Permissions == nil {
		security.Permissions = []EntryPermission{}
	}

	err = c.client.sendJSON(ctx, http.MethodPut, reqUrl, security, nil)
	if err != nil {
		if IsNotFound(err) {
			return ErrEntryNotFound
		}

		return fmt.Errorf("error while saving entry permissions: %w", err)
	}

	return nil
}

func (c *EntryPermissionsService) permissionUrl(vaultId string, entryId string) (string, error) {
	if vaultId == "" || entryId == "" {
		return "", fmt.Errorf("both entry Id and vault Id are required")
	}

	return entryUrl(c.client, nil, vaultId, entryId, entryPermissionPath)
}
