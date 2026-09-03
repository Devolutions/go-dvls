package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type AdministrativeRoleAssignments service

const administrativeRoleAssignmentsEndpoint string = "/api/v3/administrative-role-assignments"

// AdministrativeRoleAssignment represents a role assigned at a given scope
// along with its members.
type AdministrativeRoleAssignment struct {
	AdministrativeRoleId string                               `json:"administrativeRoleId"`
	RoleName             string                               `json:"roleName"`
	RoleDescription      string                               `json:"roleDescription"`
	ScopeType            AdministrativeRoleScopeType          `json:"scopeType"`
	ScopeResourceId      string                               `json:"scopeResourceId,omitempty"`
	ScopeDisplayName     string                               `json:"scopeDisplayName"`
	IsAssignable         bool                                 `json:"isAssignable"`
	IsBuiltIn            bool                                 `json:"isBuiltin"`
	IsPam                bool                                 `json:"isPam"`
	IsPrivileged         bool                                 `json:"isPrivileged"`
	Members              []AdministrativeRoleAssignmentMember `json:"members"`
}

// AdministrativeRoleAssignmentMember represents a member summary in an
// AdministrativeRoleAssignment.
type AdministrativeRoleAssignmentMember struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatarUrl"`
	Type      int    `json:"type"`
	UserType  int    `json:"userType"`
}

// AdministrativeRoleMember represents a single role assignment (one assignee
// on one role at one scope). Id is the assignment id used by RemoveMember.
type AdministrativeRoleMember struct {
	Id                   string                         `json:"id"`
	AdministrativeRoleId string                         `json:"administrativeRoleId"`
	AssigneeId           string                         `json:"assigneeId"`
	AssigneeName         string                         `json:"assigneeName"`
	AssigneeIconUrl      string                         `json:"assigneeIconUrl"`
	AssigneeType         AdministrativeRoleAssigneeType `json:"assigneeType"`
}

// UserAdministrativeRoleAssignment represents a role assignment as resolved
// for a specific assignee, including assignments inherited through user
// group membership.
type UserAdministrativeRoleAssignment struct {
	AssignmentId         string                      `json:"assignmentId"`
	AdministrativeRoleId string                      `json:"administrativeRoleId"`
	RoleName             string                      `json:"roleName"`
	RoleDescription      string                      `json:"roleDescription"`
	Permissions          []AdministrativePermission  `json:"permissions"`
	IsPam                bool                        `json:"isPam"`
	IsPrivileged         bool                        `json:"isPrivileged"`
	ScopeType            AdministrativeRoleScopeType `json:"scopeType"`
	ScopeResourceId      string                      `json:"scopeResourceId,omitempty"`
	ScopeDisplayName     string                      `json:"scopeDisplayName"`
	IsPermanent          bool                        `json:"isPermanent"`
	IsRemovable          bool                        `json:"isRemovable"`
	StartDate            *ServerTime                 `json:"startDate,omitempty"`
	EndDate              *ServerTime                 `json:"endDate,omitempty"`
	GroupId              string                      `json:"groupId,omitempty"`
	GroupName            string                      `json:"groupName,omitempty"`
}

// AdministrativeRoleMemberRequest represents a role assignment to create.
// ScopeResourceId must be nil when ScopeType is AdministrativeRoleScopeGlobal,
// and the target resource id (vault, gateway, PAM provider or business unit)
// otherwise.
type AdministrativeRoleMemberRequest struct {
	AdministrativeRoleId string                      `json:"administrativeRoleId"`
	AssigneeId           string                      `json:"assigneeId"`
	ScopeType            AdministrativeRoleScopeType `json:"scopeType"`
	ScopeResourceId      *string                     `json:"scopeResourceId,omitempty"`
}

// AdministrativeRoleMemberUpdate represents a bulk role assignment change.
type AdministrativeRoleMemberUpdate struct {
	AdministrativeRoleMemberRequest
	Action AdministrativeRoleMemberAction `json:"action"`
}

// AdministrativeRoleAssignmentFilter represents the optional filters of List.
// The zero value lists every assignment.
type AdministrativeRoleAssignmentFilter struct {
	Filter    string
	RoleId    string
	ScopeType *AdministrativeRoleScopeType
	ScopeId   string
	MemberId  string
}

// List returns all role assignments matching the filter. Requires DVLS 2026.3 or later.
func (c *AdministrativeRoleAssignments) List(filter AdministrativeRoleAssignmentFilter) ([]AdministrativeRoleAssignment, error) {
	return c.ListWithContext(context.Background(), filter)
}

// ListWithContext returns all role assignments matching the filter. Requires DVLS 2026.3 or later.
// This function handles pagination automatically and returns all assignments across all pages.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) ListWithContext(ctx context.Context, filter AdministrativeRoleAssignmentFilter) ([]AdministrativeRoleAssignment, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role assignment url: %w", err)
	}

	baseQuery := parsedUrl.Query()
	baseQuery.Set("pageSize", strconv.Itoa(listPageSize))
	if filter.Filter != "" {
		baseQuery.Set("filter", filter.Filter)
	}
	if filter.RoleId != "" {
		baseQuery.Set("roleId", filter.RoleId)
	}
	if filter.ScopeType != nil {
		baseQuery.Set("scopeType", strconv.Itoa(int(*filter.ScopeType)))
	}
	if filter.ScopeId != "" {
		baseQuery.Set("scopeId", filter.ScopeId)
	}
	if filter.MemberId != "" {
		baseQuery.Set("memberId", filter.MemberId)
	}

	var allAssignments []AdministrativeRoleAssignment
	err = fetchAllPages(func(pageNumber int) (pagedResponse, int, error) {
		baseQuery.Set("pageNumber", strconv.Itoa(pageNumber))
		parsedUrl.RawQuery = baseQuery.Encode()

		resp, err := c.client.RequestWithContext(ctx, parsedUrl.String(), http.MethodGet, nil, RequestOptions{RawBody: true})
		if err != nil {
			return pagedResponse{}, 0, fmt.Errorf("error while fetching role assignments (page %d): %w", pageNumber, err)
		}

		items, paged, err := decodeV3PagedResponse(resp.Response)
		if err != nil {
			return pagedResponse{}, 0, fmt.Errorf("error while fetching role assignments (page %d): %w", pageNumber, err)
		}

		var assignments []AdministrativeRoleAssignment
		if len(items) > 0 {
			if err := json.Unmarshal(items, &assignments); err != nil {
				return pagedResponse{}, 0, fmt.Errorf("failed to unmarshal response body (page %d): %w", pageNumber, err)
			}
		}

		allAssignments = append(allAssignments, assignments...)

		return paged, len(assignments), nil
	})
	if err != nil {
		return nil, err
	}

	return allAssignments, nil
}

// GetMembers returns the members of a role at a given scope. scopeResourceId
// must be empty when scopeType is AdministrativeRoleScopeGlobal.
// Returns an empty list when the role has no members at the scope or does not exist.
func (c *AdministrativeRoleAssignments) GetMembers(roleId string, scopeType AdministrativeRoleScopeType, scopeResourceId string) ([]AdministrativeRoleMember, error) {
	return c.GetMembersWithContext(context.Background(), roleId, scopeType, scopeResourceId)
}

// GetMembersWithContext returns the members of a role at a given scope. scopeResourceId
// must be empty when scopeType is AdministrativeRoleScopeGlobal.
// Returns an empty list when the role has no members at the scope or does not exist.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) GetMembersWithContext(ctx context.Context, roleId string, scopeType AdministrativeRoleScopeType, scopeResourceId string) ([]AdministrativeRoleMember, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, roleId, "members")
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	return c.fetchMembers(ctx, reqUrl, scopeType, scopeResourceId)
}

// ListScopeMembers returns every role assignment at a given scope, across all
// roles. scopeResourceId must be empty when scopeType is AdministrativeRoleScopeGlobal.
func (c *AdministrativeRoleAssignments) ListScopeMembers(scopeType AdministrativeRoleScopeType, scopeResourceId string) ([]AdministrativeRoleMember, error) {
	return c.ListScopeMembersWithContext(context.Background(), scopeType, scopeResourceId)
}

// ListScopeMembersWithContext returns every role assignment at a given scope, across all
// roles. scopeResourceId must be empty when scopeType is AdministrativeRoleScopeGlobal.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) ListScopeMembersWithContext(ctx context.Context, scopeType AdministrativeRoleScopeType, scopeResourceId string) ([]AdministrativeRoleMember, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "members")
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	return c.fetchMembers(ctx, reqUrl, scopeType, scopeResourceId)
}

func (c *AdministrativeRoleAssignments) fetchMembers(ctx context.Context, reqUrl string, scopeType AdministrativeRoleScopeType, scopeResourceId string) ([]AdministrativeRoleMember, error) {
	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role assignment url: %w", err)
	}

	query := parsedUrl.Query()
	query.Set("scopeType", strconv.Itoa(int(scopeType)))
	if scopeResourceId != "" {
		query.Set("scopeResourceId", scopeResourceId)
	}
	parsedUrl.RawQuery = query.Encode()

	resp, err := c.client.RequestWithContext(ctx, parsedUrl.String(), http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching role assignment members: %w", err)
	}

	return decodeV3[[]AdministrativeRoleMember](resp.Response)
}

// ListByAssignee returns the role assignments of a user, application or user
// group. When includeIndirect is true, assignments inherited through user
// group membership are included.
// Returns an empty list when the assignee has no assignments or does not exist.
func (c *AdministrativeRoleAssignments) ListByAssignee(assigneeId string, includeIndirect bool) ([]UserAdministrativeRoleAssignment, error) {
	return c.ListByAssigneeWithContext(context.Background(), assigneeId, includeIndirect)
}

// ListByAssigneeWithContext returns the role assignments of a user, application or user
// group. When includeIndirect is true, assignments inherited through user
// group membership are included.
// Returns an empty list when the assignee has no assignments or does not exist.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) ListByAssigneeWithContext(ctx context.Context, assigneeId string, includeIndirect bool) ([]UserAdministrativeRoleAssignment, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "by-assignee", assigneeId)
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse role assignment url: %w", err)
	}

	query := parsedUrl.Query()
	query.Set("includeIndirect", strconv.FormatBool(includeIndirect))
	parsedUrl.RawQuery = query.Encode()

	resp, err := c.client.RequestWithContext(ctx, parsedUrl.String(), http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching role assignments: %w", err)
	}

	return decodeV3[[]UserAdministrativeRoleAssignment](resp.Response)
}

// Me returns the role assignments of the authenticated user.
func (c *AdministrativeRoleAssignments) Me() ([]UserAdministrativeRoleAssignment, error) {
	return c.MeWithContext(context.Background())
}

// MeWithContext returns the role assignments of the authenticated user.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) MeWithContext(ctx context.Context) ([]UserAdministrativeRoleAssignment, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "me")
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching role assignments: %w", err)
	}

	return decodeV3[[]UserAdministrativeRoleAssignment](resp.Response)
}

// AvailableRoles returns the role definitions available for assignment.
func (c *AdministrativeRoleAssignments) AvailableRoles() ([]AdministrativeRole, error) {
	return c.AvailableRolesWithContext(context.Background())
}

// AvailableRolesWithContext returns the role definitions available for assignment.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) AvailableRolesWithContext(ctx context.Context) ([]AdministrativeRole, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "available-roles")
	if err != nil {
		return nil, fmt.Errorf("failed to build role assignment url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching available roles: %w", err)
	}

	return decodeV3[[]AdministrativeRole](resp.Response)
}

// AddMember assigns a role to a user, application or user group at a given scope.
func (c *AdministrativeRoleAssignments) AddMember(request AdministrativeRoleMemberRequest) error {
	return c.AddMemberWithContext(context.Background(), request)
}

// AddMemberWithContext assigns a role to a user, application or user group at a given scope.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) AddMemberWithContext(ctx context.Context, request AdministrativeRoleMemberRequest) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "members")
	if err != nil {
		return fmt.Errorf("failed to build role assignment url: %w", err)
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(requestJson), RequestOptions{RawBody: true})
	if err != nil {
		return fmt.Errorf("error while adding role assignment member: %w", err)
	}

	return checkV3SaveResponse(resp)
}

// UpdateMembers applies a batch of role assignment additions and removals.
func (c *AdministrativeRoleAssignments) UpdateMembers(updates []AdministrativeRoleMemberUpdate) error {
	return c.UpdateMembersWithContext(context.Background(), updates)
}

// UpdateMembersWithContext applies a batch of role assignment additions and removals.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) UpdateMembersWithContext(ctx context.Context, updates []AdministrativeRoleMemberUpdate) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "members", "bulk")
	if err != nil {
		return fmt.Errorf("failed to build role assignment url: %w", err)
	}

	updatesJson, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(updatesJson), RequestOptions{RawBody: true})
	if err != nil {
		return fmt.Errorf("error while updating role assignment members: %w", err)
	}

	return checkV3SaveResponse(resp)
}

// RemoveMember deletes a single role assignment based on assignmentId
// (AdministrativeRoleMember.Id).
func (c *AdministrativeRoleAssignments) RemoveMember(assignmentId string) error {
	return c.RemoveMemberWithContext(context.Background(), assignmentId)
}

// RemoveMemberWithContext deletes a single role assignment based on assignmentId
// (AdministrativeRoleMember.Id).
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) RemoveMemberWithContext(ctx context.Context, assignmentId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "members", assignmentId)
	if err != nil {
		return fmt.Errorf("failed to build role assignment url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodDelete, nil, RequestOptions{RawBody: true})
	if err != nil {
		return fmt.Errorf("error while removing role assignment member: %w", err)
	}

	return checkV3SaveResponse(resp)
}

// DeleteScope removes every member of a role at a given scope. scopeResourceId
// must be empty when scopeType is AdministrativeRoleScopeGlobal.
func (c *AdministrativeRoleAssignments) DeleteScope(roleId string, scopeType AdministrativeRoleScopeType, scopeResourceId string) error {
	return c.DeleteScopeWithContext(context.Background(), roleId, scopeType, scopeResourceId)
}

// DeleteScopeWithContext removes every member of a role at a given scope. scopeResourceId
// must be empty when scopeType is AdministrativeRoleScopeGlobal.
// The provided context can be used to cancel the request.
func (c *AdministrativeRoleAssignments) DeleteScopeWithContext(ctx context.Context, roleId string, scopeType AdministrativeRoleScopeType, scopeResourceId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, administrativeRoleAssignmentsEndpoint, "scope")
	if err != nil {
		return fmt.Errorf("failed to build role assignment url: %w", err)
	}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return fmt.Errorf("failed to parse role assignment url: %w", err)
	}

	query := parsedUrl.Query()
	query.Set("roleId", roleId)
	query.Set("scopeType", strconv.Itoa(int(scopeType)))
	if scopeResourceId != "" {
		query.Set("scopeResourceId", scopeResourceId)
	}
	parsedUrl.RawQuery = query.Encode()

	resp, err := c.client.RequestWithContext(ctx, parsedUrl.String(), http.MethodDelete, nil, RequestOptions{RawBody: true})
	if err != nil {
		return fmt.Errorf("error while deleting role assignment scope: %w", err)
	}

	return checkV3SaveResponse(resp)
}
