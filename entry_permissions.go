package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type EntryPermissionsService service

// EntryPermission represents the principals granted a single right on an
// entry or folder. Roles contains principal ids (users, user groups and
// applications share the same list).
type EntryPermission struct {
	Override SecurityRoleOverride `json:"override"`
	Right    SecurityRoleRight    `json:"right"`
	Roles    []string             `json:"roles"`
}

// EntrySecurity represents the permission block of an entry or folder.
// The View right is not part of Permissions: it is controlled by ViewOverride
// and ViewRoles. When RoleOverride is not Custom or CustomInherited, the
// server forces ViewOverride to RoleOverride and clears ViewRoles, so a Get
// following a Set can return normalized values.
type EntrySecurity struct {
	RoleOverride SecurityRoleOverride `json:"roleOverride"`
	ViewOverride SecurityRoleOverride `json:"viewOverride"`
	ViewRoles    []string             `json:"viewRoles"`
	Permissions  []EntryPermission    `json:"permissions"`
}

// Get returns the EntrySecurity of the entry or folder specified by entryId.
// Returns ErrEntryNotFound if no entry is found.
func (c *EntryPermissionsService) Get(entryId string) (EntrySecurity, error) {
	return c.GetWithContext(context.Background(), entryId)
}

// GetWithContext returns the EntrySecurity of the entry or folder specified by entryId.
// Returns ErrEntryNotFound if no entry is found.
// The provided context can be used to cancel the request.
func (c *EntryPermissionsService) GetWithContext(ctx context.Context, entryId string) (EntrySecurity, error) {
	connectionJson, err := c.getPartialConnection(ctx, entryId)
	if err != nil {
		return EntrySecurity{}, err
	}

	var connection struct {
		Security json.RawMessage `json:"security"`
	}
	err = json.Unmarshal(connectionJson, &connection)
	if err != nil {
		return EntrySecurity{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if !jsonValuePresent(connection.Security) {
		return EntrySecurity{}, nil
	}

	var security EntrySecurity
	err = json.Unmarshal(connection.Security, &security)
	if err != nil {
		return EntrySecurity{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return security, nil
}

// Set replaces the permission block (RoleOverride, ViewOverride, ViewRoles and
// Permissions) of the entry or folder specified by entryId. Every other entry
// field and security setting is preserved.
// Set fetches the entry and saves it back in two requests; a concurrent
// modification of the entry between the two requests is overwritten.
func (c *EntryPermissionsService) Set(entryId string, security EntrySecurity) error {
	return c.SetWithContext(context.Background(), entryId, security)
}

// SetWithContext replaces the permission block (RoleOverride, ViewOverride, ViewRoles and
// Permissions) of the entry or folder specified by entryId. Every other entry
// field and security setting is preserved.
// SetWithContext fetches the entry and saves it back in two requests; a concurrent
// modification of the entry between the two requests is overwritten.
// The provided context can be used to cancel the request.
func (c *EntryPermissionsService) SetWithContext(ctx context.Context, entryId string, security EntrySecurity) error {
	err := validateEntrySecurity(entryId, security)
	if err != nil {
		return err
	}

	connectionJson, err := c.getPartialConnection(ctx, entryId)
	if err != nil {
		return err
	}

	connection := map[string]json.RawMessage{}
	err = json.Unmarshal(connectionJson, &connection)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	securityNode := map[string]json.RawMessage{}
	if existing, ok := connection["security"]; ok && jsonValuePresent(existing) {
		err = json.Unmarshal(existing, &securityNode)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	err = setSecurityPermissionFields(securityNode, security)
	if err != nil {
		return err
	}

	securityJson, err := json.Marshal(securityNode)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}
	connection["security"] = securityJson

	err = fixupFolderGroup(connection)
	if err != nil {
		return err
	}

	saveJson, err := json.Marshal(connection)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPut, bytes.NewBuffer(saveJson))
	if err != nil {
		return fmt.Errorf("error while saving entry permissions: %w", err)
	}

	return resp.CheckRespSaveResult()
}

// getPartialConnection fetches the raw partial connection JSON. The save
// endpoint overwrites the whole entry, so every field must be preserved byte
// for byte; never decode the connection into a typed struct before saving.
func (c *EntryPermissionsService) getPartialConnection(ctx context.Context, entryId string) (json.RawMessage, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return nil, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		if IsNotFound(err) {
			return nil, ErrEntryNotFound
		}

		return nil, fmt.Errorf("error while fetching entry: %w", err)
	}

	if SaveResult(resp.Result) == SaveResultNotFound {
		return nil, ErrEntryNotFound
	}
	if err = resp.CheckRespSaveResult(); err != nil {
		return nil, err
	}

	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	err = json.Unmarshal(resp.Response, &envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return envelope.Data, nil
}

func validateEntrySecurity(entryId string, security EntrySecurity) error {
	if entryId == "" {
		return fmt.Errorf("entry id is required")
	}

	seen := map[SecurityRoleRight]struct{}{}
	for _, permission := range security.Permissions {
		if permission.Right == SecurityRoleRightView {
			return fmt.Errorf("the View right cannot be set through Permissions, use ViewOverride and ViewRoles")
		}

		if _, ok := seen[permission.Right]; ok {
			return fmt.Errorf("duplicate permission right %s", permission.Right)
		}
		seen[permission.Right] = struct{}{}
	}

	customOverride := security.RoleOverride == SecurityRoleOverrideCustom || security.RoleOverride == SecurityRoleOverrideCustomInherited
	if !customOverride && (len(security.Permissions) > 0 || len(security.ViewRoles) > 0) {
		return fmt.Errorf("Permissions and ViewRoles require RoleOverride to be Custom or CustomInherited")
	}

	return nil
}

func setSecurityPermissionFields(securityNode map[string]json.RawMessage, security EntrySecurity) error {
	normalized := security
	if normalized.ViewRoles == nil {
		normalized.ViewRoles = []string{}
	}

	normalized.Permissions = make([]EntryPermission, len(security.Permissions))
	copy(normalized.Permissions, security.Permissions)
	for i := range normalized.Permissions {
		if normalized.Permissions[i].Roles == nil {
			normalized.Permissions[i].Roles = []string{}
		}
	}

	normalizedJson, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	fields := map[string]json.RawMessage{}
	err = json.Unmarshal(normalizedJson, &fields)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	for field, value := range fields {
		securityNode[field] = value
	}

	return nil
}

// fixupFolderGroup rewrites the group field of folders to the parent path
// before saving. The partial connection endpoint returns a folder's group as
// its own full path; saving it unchanged would move the folder into itself.
func fixupFolderGroup(connection map[string]json.RawMessage) error {
	var connectionType int
	if raw, ok := connection["connectionType"]; ok {
		if err := json.Unmarshal(raw, &connectionType); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	var connectionSubType string
	if raw, ok := connection["connectionSubType"]; ok {
		if err := json.Unmarshal(raw, &connectionSubType); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	isFolder := connectionType == int(ServerConnectionGroup) ||
		(connectionType == int(ServerConnectionPAM) && strings.EqualFold(connectionSubType, "Group"))
	if !isFolder {
		return nil
	}

	var name, group string
	if raw, ok := connection["name"]; ok {
		if err := json.Unmarshal(raw, &name); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}
	if raw, ok := connection["group"]; ok {
		if err := json.Unmarshal(raw, &group); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	if group == name {
		group = ""
	} else {
		group = strings.TrimSuffix(group, "\\"+name)
	}
	group = strings.TrimRight(group, "\\")

	groupJson, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}
	connection["group"] = groupJson

	return nil
}
