package dvls

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type AdministrativeRoles service

// AdministrativeRole represents a DVLS administrative role definition.
type AdministrativeRole struct {
	Id                        string                        `json:"id,omitempty"`
	Name                      string                        `json:"name"`
	Description               string                        `json:"description"`
	Permissions               []AdministrativePermission    `json:"permissions"`
	SupportedScopes           []AdministrativeRoleScopeType `json:"supportedScopes"`
	RequiredPermissionOnScope *AdministrativePermission     `json:"requiredPermissionOnScope,omitempty"`
	IsAssignable              bool                          `json:"isAssignable,omitempty"`
	IsBuiltIn                 bool                          `json:"isBuiltIn,omitempty"`
	IsPam                     bool                          `json:"isPam,omitempty"`
	IsPrivileged              bool                          `json:"isPrivileged,omitempty"`
	IsUsed                    bool                          `json:"isUsed,omitempty"`
}

const administrativeRolesEndpoint string = "/api/v3/administrative-roles"

var ErrAdministrativeRoleNotFound = fmt.Errorf("administrative role not found")
var ErrMultipleAdministrativeRolesFound = fmt.Errorf("multiple administrative roles found")

// List returns all administrative roles. Requires DVLS 2026.3 or later.
func (c *AdministrativeRoles) List() ([]AdministrativeRole, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all administrative roles. Requires DVLS 2026.3 or later.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoles) ListWithContext(ctx context.Context) ([]AdministrativeRole, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRolesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build administrative role url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching administrative roles: %w", err)
	}

	return decodeV3[[]AdministrativeRole](resp.Response)
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
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRolesEndpoint, roleId)
	if err != nil {
		return AdministrativeRole{}, fmt.Errorf("failed to build administrative role url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		if IsNotFound(err) {
			return AdministrativeRole{}, ErrAdministrativeRoleNotFound
		}

		return AdministrativeRole{}, fmt.Errorf("error while fetching administrative role: %w", err)
	}

	role, err := decodeV3[AdministrativeRole](resp.Response)
	if err != nil {
		if isV3ResultNotFound(err) {
			return AdministrativeRole{}, ErrAdministrativeRoleNotFound
		}

		return AdministrativeRole{}, err
	}

	return role, nil
}

// GetByName returns a single AdministrativeRole based on name.
// Returns ErrAdministrativeRoleNotFound if no role is found.
// Returns ErrMultipleAdministrativeRolesFound if more than one role matches the name.
func (c *AdministrativeRoles) GetByName(name string) (AdministrativeRole, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single AdministrativeRole based on name.
// Returns ErrAdministrativeRoleNotFound if no role is found.
// Returns ErrMultipleAdministrativeRolesFound if more than one role matches the name.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoles) GetByNameWithContext(ctx context.Context, name string) (AdministrativeRole, error) {
	roles, err := c.ListWithContext(ctx)
	if err != nil {
		return AdministrativeRole{}, err
	}

	var matches []AdministrativeRole
	for _, r := range roles {
		if r.Name == name {
			matches = append(matches, r)
		}
	}

	if len(matches) == 0 {
		return AdministrativeRole{}, ErrAdministrativeRoleNotFound
	}

	if len(matches) > 1 {
		return AdministrativeRole{}, ErrMultipleAdministrativeRolesFound
	}

	return matches[0], nil
}
