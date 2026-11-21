package dvls

import (
	"context"
	"fmt"
	"net/url"
)

const userGroupEndpoint = "user-group"

var ErrUserGroupNotFound = fmt.Errorf("user group not found")
var ErrMultipleUserGroupsFound = fmt.Errorf("multiple user groups found")

type UserGroups service

// UserGroupType is the source of a user group.
type UserGroupType string

const (
	UserGroupTypeCustom          UserGroupType = "Custom"
	UserGroupTypeActiveDirectory UserGroupType = "ActiveDirectory"
	UserGroupTypeEntraId         UserGroupType = "EntraId"
	UserGroupTypeOkta            UserGroupType = "Okta"
	UserGroupTypePingOne         UserGroupType = "PingOne"
	UserGroupTypeKeycloak        UserGroupType = "Keycloak"
)

// UserGroup represents a DVLS user group. IsWorkspaceOwner reports whether the group makes its members workspace
// owners.
type UserGroup struct {
	Id               string        `json:"id"`
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Type             UserGroupType `json:"type"`
	IsWorkspaceOwner bool          `json:"isWorkspaceOwner"`
	IsEnabled        bool          `json:"isEnabled"`
}

// List returns all user groups.
func (c *UserGroups) List() ([]UserGroup, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all user groups.
// The provided context can be used to cancel the request.
func (c *UserGroups) ListWithContext(ctx context.Context) ([]UserGroup, error) {
	return c.list(ctx, nil)
}

func (c *UserGroups) list(ctx context.Context, query url.Values) ([]UserGroup, error) {
	groups, err := listAll[UserGroup](ctx, c.client, query, userGroupEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching user groups: %w", err)
	}

	return groups, nil
}

// Get returns a single user group based on userGroupId.
func (c *UserGroups) Get(userGroupId string) (UserGroup, error) {
	return c.GetWithContext(context.Background(), userGroupId)
}

// GetWithContext returns a single user group based on userGroupId.
// The provided context can be used to cancel the request.
func (c *UserGroups) GetWithContext(ctx context.Context, userGroupId string) (UserGroup, error) {
	reqUrl, err := c.client.v2Url(nil, userGroupEndpoint, userGroupId)
	if err != nil {
		return UserGroup{}, fmt.Errorf("failed to build user group url: %w", err)
	}

	var group UserGroup
	err = c.client.getJSON(ctx, reqUrl, &group)
	if err != nil {
		return UserGroup{}, fmt.Errorf("error while fetching user group: %w", err)
	}

	return group, nil
}

// GetByName returns a single user group based on name. The name matches ignoring case.
// Returns ErrUserGroupNotFound if no user group is found.
// Returns ErrMultipleUserGroupsFound if more than one user group matches the name.
func (c *UserGroups) GetByName(name string) (UserGroup, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single user group based on name. The name matches ignoring case.
// Returns ErrUserGroupNotFound if no user group is found.
// Returns ErrMultipleUserGroupsFound if more than one user group matches the name.
// The provided context can be used to cancel the request.
func (c *UserGroups) GetByNameWithContext(ctx context.Context, name string) (UserGroup, error) {
	groups, err := c.list(ctx, url.Values{"name": {name}})
	if err != nil {
		return UserGroup{}, err
	}

	return singleByName(groups, name, func(g UserGroup) string { return g.Name }, ErrUserGroupNotFound, ErrMultipleUserGroupsFound)
}
