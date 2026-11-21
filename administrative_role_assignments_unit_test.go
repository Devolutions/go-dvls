package dvls

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAssignmentID   = "99999999-8888-7777-6666-555555555555"
	testAssignmentJSON = `{"id":"` + testAssignmentID + `","roleId":"` + testRoleID + `","roleName":"Vault member","assigneeId":"` + testAssigneeID + `","assigneeName":"app","assigneeType":"ApplicationIdentity","scopeType":"Vault","scopeId":"` + testVaultID + `","scopeName":"Prod"}`
)

func TestRoleAssignmentsNew(t *testing.T) {
	tests := []struct {
		name     string
		request  AdministrativeRoleAssignmentRequest
		wantBody string
	}{
		{
			name:     "Vault",
			request:  AdministrativeRoleAssignmentRequest{RoleId: testRoleID, AssigneeId: testAssigneeID, ScopeType: AdministrativeRoleScopeVault, ScopeId: testVaultID},
			wantBody: `{"roleId":"` + testRoleID + `","assigneeId":"` + testAssigneeID + `","scopeType":"Vault","scopeId":"` + testVaultID + `"}`,
		},
		{
			name:     "WorkspaceOmitsScopeId",
			request:  AdministrativeRoleAssignmentRequest{RoleId: testRoleID, AssigneeId: testAssigneeID, ScopeType: AdministrativeRoleScopeWorkspace},
			wantBody: `{"roleId":"` + testRoleID + `","assigneeId":"` + testAssigneeID + `","scopeType":"Workspace"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/role-assignment", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var body json.RawMessage
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.JSONEq(t, tt.wantBody, string(body))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(testAssignmentJSON))
			})

			client := newTestClient(t, mux)

			assignment, err := client.AdministrativeRoleAssignments.New(tt.request)
			require.NoError(t, err)
			assert.Equal(t, AdministrativeRoleAssignment{
				Id:           testAssignmentID,
				RoleId:       testRoleID,
				RoleName:     "Vault member",
				AssigneeId:   testAssigneeID,
				AssigneeName: "app",
				AssigneeType: AdministrativeRoleAssigneeApplicationIdentity,
				ScopeType:    AdministrativeRoleScopeVault,
				ScopeId:      testVaultID,
				ScopeName:    "Prod",
			}, assignment)
		})
	}
}

func TestRoleAssignmentsList_Filters(t *testing.T) {
	tests := []struct {
		name      string
		filter    AdministrativeRoleAssignmentFilter
		wantQuery url.Values
	}{
		{name: "None", filter: AdministrativeRoleAssignmentFilter{}, wantQuery: url.Values{}},
		{
			name:      "RoleAtScope",
			filter:    AdministrativeRoleAssignmentFilter{RoleId: testRoleID, ScopeType: AdministrativeRoleScopeVault, ScopeId: testVaultID},
			wantQuery: url.Values{"roleId": {testRoleID}, "scopeType": {"Vault"}, "scopeId": {testVaultID}},
		},
		{name: "Assignee", filter: AdministrativeRoleAssignmentFilter{AssigneeId: testAssigneeID}, wantQuery: url.Values{"assigneeId": {testAssigneeID}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/api/v2/role-assignment", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				query := r.URL.Query()
				query.Del("pageSize")
				query.Del("pageNumber")
				assert.Equal(t, tt.wantQuery, query)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":1,"totalPage":1,"data":[` + testAssignmentJSON + `]}`))
			})

			client := newTestClient(t, mux)

			assignments, err := client.AdministrativeRoleAssignments.List(tt.filter)
			require.NoError(t, err)
			require.Len(t, assignments, 1)
			assert.Equal(t, testAssignmentID, assignments[0].Id)
		})
	}
}

func TestRoleAssignmentsGetAndDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/role-assignment/"+testAssignmentID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte(testAssignmentJSON))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, mux)

	assignment, err := client.AdministrativeRoleAssignments.Get(testAssignmentID)
	require.NoError(t, err)
	assert.Equal(t, AdministrativeRoleScopeVault, assignment.ScopeType)

	require.NoError(t, client.AdministrativeRoleAssignments.Delete(testAssignmentID))
}

func TestRoleAssignmentsNew_Duplicate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/role-assignment", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"message":"The role assignment already exists."}`))
	})

	client := newTestClient(t, mux)

	_, err := client.AdministrativeRoleAssignments.New(AdministrativeRoleAssignmentRequest{RoleId: testRoleID, AssigneeId: testAssigneeID, ScopeType: AdministrativeRoleScopeWorkspace})
	assert.ErrorIs(t, err, ErrAdministrativeRoleAssignmentAlreadyExists)
	assert.True(t, hasStatusCode(err, http.StatusConflict), "the RequestError stays in the chain")
}
