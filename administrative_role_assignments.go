package dvls

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type AdministrativeRoleAssignments service

const administrativeRoleAssignmentsEndpoint string = "role-assignment"

var ErrAdministrativeRoleAssignmentAlreadyExists = fmt.Errorf("administrative role assignment already exists")

// AdministrativeRoleAssignment represents a role granted to one assignee at one scope. ScopeId is empty for the
// Workspace scope.
type AdministrativeRoleAssignment struct {
	Id           string                         `json:"id"`
	RoleId       string                         `json:"roleId"`
	RoleName     string                         `json:"roleName"`
	AssigneeId   string                         `json:"assigneeId"`
	AssigneeName string                         `json:"assigneeName"`
	AssigneeType AdministrativeRoleAssigneeType `json:"assigneeType"`
	ScopeType    AdministrativeRoleScopeType    `json:"scopeType"`
	ScopeId      string                         `json:"scopeId,omitempty"`
	ScopeName    string                         `json:"scopeName"`
}

// AdministrativeRoleAssignmentRequest represents a role assignment to create. ScopeId must be empty when ScopeType
// is AdministrativeRoleScopeWorkspace, and the vault or business unit id otherwise.
type AdministrativeRoleAssignmentRequest struct {
	RoleId     string                      `json:"roleId"`
	AssigneeId string                      `json:"assigneeId"`
	ScopeType  AdministrativeRoleScopeType `json:"scopeType"`
	ScopeId    string                      `json:"scopeId,omitempty"`
}

// AdministrativeRoleAssignmentFilter represents the optional filters of List. An empty field is not applied, and
// the zero value lists every assignment visible to the caller. ScopeId needs a ScopeType other than Workspace. A
// filter naming a role, assignee or scope that does not exist is rejected with an HTTP 400.
type AdministrativeRoleAssignmentFilter struct {
	RoleId     string
	AssigneeId string
	ScopeType  AdministrativeRoleScopeType
	ScopeId    string
}

// New creates a permanent role assignment and returns it.
// Returns ErrAdministrativeRoleAssignmentAlreadyExists if the same role is already assigned to the assignee at the
// scope.
func (c *AdministrativeRoleAssignments) New(request AdministrativeRoleAssignmentRequest) (AdministrativeRoleAssignment, error) {
	return c.NewWithContext(context.Background(), request)
}

// NewWithContext creates a permanent role assignment and returns it.
// Returns ErrAdministrativeRoleAssignmentAlreadyExists if the same role is already assigned to the assignee at the
// scope.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) NewWithContext(ctx context.Context, request AdministrativeRoleAssignmentRequest) (AdministrativeRoleAssignment, error) {
	reqUrl, err := c.client.v2Url(nil, administrativeRoleAssignmentsEndpoint)
	if err != nil {
		return AdministrativeRoleAssignment{}, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	var assignment AdministrativeRoleAssignment
	err = c.client.sendJSON(ctx, http.MethodPost, reqUrl, request, &assignment)
	if err != nil {
		if hasStatusCode(err, http.StatusConflict) {
			return AdministrativeRoleAssignment{}, fmt.Errorf("%w: %w", ErrAdministrativeRoleAssignmentAlreadyExists, err)
		}

		return AdministrativeRoleAssignment{}, fmt.Errorf("error while creating role assignment: %w", err)
	}

	return assignment, nil
}

// List returns the role assignments matching the filter.
func (c *AdministrativeRoleAssignments) List(filter AdministrativeRoleAssignmentFilter) ([]AdministrativeRoleAssignment, error) {
	return c.ListWithContext(context.Background(), filter)
}

// ListWithContext returns the role assignments matching the filter.
// This function handles pagination automatically and returns all assignments across all pages.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) ListWithContext(ctx context.Context, filter AdministrativeRoleAssignmentFilter) ([]AdministrativeRoleAssignment, error) {
	query := url.Values{}
	if filter.RoleId != "" {
		query.Set("roleId", filter.RoleId)
	}
	if filter.AssigneeId != "" {
		query.Set("assigneeId", filter.AssigneeId)
	}
	if filter.ScopeType != "" {
		query.Set("scopeType", string(filter.ScopeType))
	}
	if filter.ScopeId != "" {
		query.Set("scopeId", filter.ScopeId)
	}

	assignments, err := listAll[AdministrativeRoleAssignment](ctx, c.client, query, administrativeRoleAssignmentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching role assignments: %w", err)
	}

	return assignments, nil
}

// Get returns a single role assignment based on assignmentId.
func (c *AdministrativeRoleAssignments) Get(assignmentId string) (AdministrativeRoleAssignment, error) {
	return c.GetWithContext(context.Background(), assignmentId)
}

// GetWithContext returns a single role assignment based on assignmentId.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) GetWithContext(ctx context.Context, assignmentId string) (AdministrativeRoleAssignment, error) {
	reqUrl, err := c.client.v2Url(nil, administrativeRoleAssignmentsEndpoint, assignmentId)
	if err != nil {
		return AdministrativeRoleAssignment{}, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	var assignment AdministrativeRoleAssignment
	err = c.client.getJSON(ctx, reqUrl, &assignment)
	if err != nil {
		return AdministrativeRoleAssignment{}, fmt.Errorf("error while fetching role assignment: %w", err)
	}

	return assignment, nil
}

// Delete removes a single role assignment based on assignmentId.
func (c *AdministrativeRoleAssignments) Delete(assignmentId string) error {
	return c.DeleteWithContext(context.Background(), assignmentId)
}

// DeleteWithContext removes a single role assignment based on assignmentId.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) DeleteWithContext(ctx context.Context, assignmentId string) error {
	reqUrl, err := c.client.v2Url(nil, administrativeRoleAssignmentsEndpoint, assignmentId)
	if err != nil {
		return fmt.Errorf("failed to build role assignment url: %w", err)
	}

	err = c.client.sendJSON(ctx, http.MethodDelete, reqUrl, nil, nil)
	if err != nil {
		return fmt.Errorf("error while deleting role assignment: %w", err)
	}

	return nil
}
