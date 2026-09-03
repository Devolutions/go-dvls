package dvls

import (
	"context"
	"fmt"
)

const userGroupListEndpoint = "/api/security/roles/basic"

var ErrUserGroupNotFound = fmt.Errorf("user group not found")
var ErrMultipleUserGroupsFound = fmt.Errorf("multiple user groups found")

type UserGroups service

// UserGroup represents a DVLS user group.
type UserGroup struct {
	Id              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	IsAdministrator bool          `json:"isAdministrator"`
	Type            UserGroupType `json:"roleType"`
}

// List returns all user groups.
func (c *UserGroups) List() ([]UserGroup, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all user groups.
// The provided context can be used to cancel the request.
func (c *UserGroups) ListWithContext(ctx context.Context) ([]UserGroup, error) {
	return fetchDataList[UserGroup](ctx, c.client, userGroupListEndpoint, "user groups")
}

// GetByName returns a single user group based on name.
// Returns ErrUserGroupNotFound if no user group is found.
// Returns ErrMultipleUserGroupsFound if more than one user group matches the name.
func (c *UserGroups) GetByName(name string) (UserGroup, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single user group based on name.
// Returns ErrUserGroupNotFound if no user group is found.
// Returns ErrMultipleUserGroupsFound if more than one user group matches the name.
// The provided context can be used to cancel the request.
func (c *UserGroups) GetByNameWithContext(ctx context.Context, name string) (UserGroup, error) {
	groups, err := c.ListWithContext(ctx)
	if err != nil {
		return UserGroup{}, err
	}

	return singleByName(groups, name, func(g UserGroup) string { return g.Name }, ErrUserGroupNotFound, ErrMultipleUserGroupsFound)
}
