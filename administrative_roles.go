package dvls

import (
	"context"
	"fmt"
	"net/url"
)

type AdministrativeRoles service

// AdministrativeRole represents a DVLS role definition. AssignableScopes is empty for a role that cannot be granted.
type AdministrativeRole struct {
	Id               string                        `json:"id,omitempty"`
	Name             string                        `json:"name"`
	Description      string                        `json:"description"`
	Permissions      []AdministrativePermission    `json:"permissions"`
	AssignableScopes []AdministrativeRoleScopeType `json:"assignableScopes"`
	IsBuiltIn        bool                          `json:"isBuiltIn,omitempty"`
}

const administrativeRolesEndpoint string = "role"

var ErrAdministrativeRoleNotFound = fmt.Errorf("administrative role not found")
var ErrMultipleAdministrativeRolesFound = fmt.Errorf("multiple administrative roles found")

// List returns all administrative roles.
func (c *AdministrativeRoles) List() ([]AdministrativeRole, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all administrative roles.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoles) ListWithContext(ctx context.Context) ([]AdministrativeRole, error) {
	return c.list(ctx, nil)
}

func (c *AdministrativeRoles) list(ctx context.Context, query url.Values) ([]AdministrativeRole, error) {
	roles, err := listAll[AdministrativeRole](ctx, c.client, query, administrativeRolesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching administrative roles: %w", err)
	}

	return roles, nil
}

// Get returns a single AdministrativeRole based on roleId.
// Returns ErrAdministrativeRoleNotFound if no role is found.
func (c *AdministrativeRoles) Get(roleId string) (AdministrativeRole, error) {
	return c.GetWithContext(context.Background(), roleId)
}

// GetWithContext returns a single AdministrativeRole based on roleId.
// Returns ErrAdministrativeRoleNotFound if no role is found.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoles) GetWithContext(ctx context.Context, roleId string) (AdministrativeRole, error) {
	reqUrl, err := c.client.v2Url(nil, administrativeRolesEndpoint, roleId)
	if err != nil {
		return AdministrativeRole{}, fmt.Errorf("failed to build administrative role url: %w", err)
	}

	var role AdministrativeRole
	err = c.client.getJSON(ctx, reqUrl, &role)
	if err != nil {
		if IsNotFound(err) {
			return AdministrativeRole{}, ErrAdministrativeRoleNotFound
		}

		return AdministrativeRole{}, fmt.Errorf("error while fetching administrative role: %w", err)
	}

	return role, nil
}

// GetByName returns a single AdministrativeRole based on name. The name matches ignoring case.
// Returns ErrAdministrativeRoleNotFound if no role is found.
// Returns ErrMultipleAdministrativeRolesFound if more than one role matches the name.
func (c *AdministrativeRoles) GetByName(name string) (AdministrativeRole, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single AdministrativeRole based on name. The name matches ignoring case.
// Returns ErrAdministrativeRoleNotFound if no role is found.
// Returns ErrMultipleAdministrativeRolesFound if more than one role matches the name.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoles) GetByNameWithContext(ctx context.Context, name string) (AdministrativeRole, error) {
	roles, err := c.list(ctx, url.Values{"name": {name}})
	if err != nil {
		return AdministrativeRole{}, err
	}

	return singleByName(roles, name, func(r AdministrativeRole) string { return r.Name }, ErrAdministrativeRoleNotFound, ErrMultipleAdministrativeRolesFound)
}
