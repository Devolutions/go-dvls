package dvls

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleAssignmentsList_PaginationWrapped(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		assert.Equal(t, "100", r.URL.Query().Get("pageSize"))
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("pageNumber") {
		case "1":
			w.Write([]byte(`{"result":1,"data":{"currentPage":1,"pageSize":1,"totalCount":2,"totalPage":2,"data":[
				{"administrativeRoleId":"` + BuiltinRoleVaultUserId + `","roleName":"Vault User","scopeType":2,"scopeResourceId":"` + testVaultID + `","scopeDisplayName":"Test Vault","isBuiltin":true,"isAssignable":true,"members":[{"id":"` + testAssigneeID + `","name":"someuser","type":0,"userType":0}]}
			]}}`))
		case "2":
			w.Write([]byte(`{"result":1,"data":{"currentPage":2,"pageSize":1,"totalCount":2,"totalPage":2,"data":[
				{"administrativeRoleId":"` + testRoleID + `","roleName":"Custom Role","scopeType":0,"scopeDisplayName":"Global","members":[]}
			]}}`))
		default:
			t.Errorf("unexpected page %q", r.URL.Query().Get("pageNumber"))
		}
	})

	client := newTestClient(t, mux)

	assignments, err := client.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{})
	require.NoError(t, err)
	require.Len(t, assignments, 2)
	assert.Equal(t, BuiltinRoleVaultUserId, assignments[0].AdministrativeRoleId)
	assert.True(t, assignments[0].IsBuiltIn)
	assert.Equal(t, AdministrativeRoleScopeVault, assignments[0].ScopeType)
	assert.Equal(t, testVaultID, assignments[0].ScopeResourceId)
	assert.True(t, assignments[0].IsAssignable)
	require.Len(t, assignments[0].Members, 1)
	assert.Equal(t, testAssigneeID, assignments[0].Members[0].Id)
	assert.Equal(t, AdministrativeRoleScopeGlobal, assignments[1].ScopeType)
}

func TestRoleAssignmentsList_Unwrapped(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":100,"totalCount":1,"totalPage":1,"data":[
			{"administrativeRoleId":"` + testRoleID + `","roleName":"Custom Role","scopeType":0,"scopeDisplayName":"Global","members":[]}
		]}`))
	})

	client := newTestClient(t, mux)

	assignments, err := client.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{})
	require.NoError(t, err)
	require.Len(t, assignments, 1)
	assert.Equal(t, testRoleID, assignments[0].AdministrativeRoleId)
}

func TestRoleAssignmentsList_Filters(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, testRoleID, query.Get("roleId"))
		assert.Equal(t, "2", query.Get("scopeType"))
		assert.Equal(t, testVaultID, query.Get("scopeId"))
		assert.Equal(t, testAssigneeID, query.Get("memberId"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":{"currentPage":1,"pageSize":100,"totalCount":0,"totalPage":0,"data":[]}}`))
	})

	client := newTestClient(t, mux)

	scopeType := AdministrativeRoleScopeVault
	assignments, err := client.AdministrativeRoleAssignments.List(AdministrativeRoleAssignmentFilter{
		RoleId:    testRoleID,
		ScopeType: &scopeType,
		ScopeId:   testVaultID,
		MemberId:  testAssigneeID,
	})
	require.NoError(t, err)
	assert.Empty(t, assignments)
}

func TestRoleAssignmentsGetMembers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/"+testRoleID+"/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		assert.Equal(t, "2", r.URL.Query().Get("scopeType"))
		assert.Equal(t, testVaultID, r.URL.Query().Get("scopeResourceId"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"id":"assignment-1","administrativeRoleId":"` + testRoleID + `","assigneeId":"` + testAssigneeID + `","assigneeName":"Some Group","assigneeType":2}
		]}`))
	})

	client := newTestClient(t, mux)

	members, err := client.AdministrativeRoleAssignments.GetMembers(testRoleID, AdministrativeRoleScopeVault, testVaultID)
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, "assignment-1", members[0].Id)
	assert.Equal(t, testAssigneeID, members[0].AssigneeId)
	assert.Equal(t, AdministrativeRoleAssigneeUserGroup, members[0].AssigneeType)
}

func TestRoleAssignmentsGetMembers_GlobalScopeOmitsResourceId(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/"+testRoleID+"/members", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "0", r.URL.Query().Get("scopeType"))
		assert.False(t, r.URL.Query().Has("scopeResourceId"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[]}`))
	})

	client := newTestClient(t, mux)

	members, err := client.AdministrativeRoleAssignments.GetMembers(testRoleID, AdministrativeRoleScopeGlobal, "")
	require.NoError(t, err)
	assert.Empty(t, members)
}

func TestRoleAssignmentsListByAssignee(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/by-assignee/"+testAssigneeID, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("includeIndirect"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1,"data":[
			{"assignmentId":"assignment-1","administrativeRoleId":"` + BuiltinRoleVaultUserId + `","roleName":"Vault User","permissions":[332],"scopeType":2,"scopeResourceId":"` + testVaultID + `","scopeDisplayName":"Test Vault","isPermanent":true,"groupId":"group-1","groupName":"Some Group","startDate":"2026-08-14T12:00:00"}
		]}`))
	})

	client := newTestClient(t, mux)

	assignments, err := client.AdministrativeRoleAssignments.ListByAssignee(testAssigneeID, true)
	require.NoError(t, err)
	require.Len(t, assignments, 1)
	assert.Equal(t, BuiltinRoleVaultUserId, assignments[0].AdministrativeRoleId)
	assert.Equal(t, []AdministrativePermission{AdministrativePermissionRepositoriesContentView}, assignments[0].Permissions)
	assert.Equal(t, "Some Group", assignments[0].GroupName)
	require.NotNil(t, assignments[0].StartDate)
	assert.Equal(t, 2026, assignments[0].StartDate.Year())
}

func TestRoleAssignmentsAddMember(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody map[string]any
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, BuiltinRoleVaultUserId, reqBody["administrativeRoleId"])
		assert.Equal(t, testAssigneeID, reqBody["assigneeId"])
		assert.Equal(t, float64(2), reqBody["scopeType"])
		assert.Equal(t, testVaultID, reqBody["scopeResourceId"])

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1}`))
	})

	client := newTestClient(t, mux)

	scopeResourceId := testVaultID
	err := client.AdministrativeRoleAssignments.AddMember(AdministrativeRoleMemberRequest{
		AdministrativeRoleId: BuiltinRoleVaultUserId,
		AssigneeId:           testAssigneeID,
		ScopeType:            AdministrativeRoleScopeVault,
		ScopeResourceId:      &scopeResourceId,
	})
	assert.NoError(t, err)
}

func TestRoleAssignmentsUpdateMembers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/members/bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody []map[string]any
		require.NoError(t, json.Unmarshal(body, &reqBody))
		require.Len(t, reqBody, 2)
		assert.Equal(t, float64(0), reqBody[0]["action"])
		assert.NotContains(t, reqBody[0], "scopeResourceId")
		assert.Equal(t, float64(1), reqBody[1]["action"])
		assert.Equal(t, testVaultID, reqBody[1]["scopeResourceId"])

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1}`))
	})

	client := newTestClient(t, mux)

	scopeResourceId := testVaultID
	err := client.AdministrativeRoleAssignments.UpdateMembers([]AdministrativeRoleMemberUpdate{
		{
			AdministrativeRoleMemberRequest: AdministrativeRoleMemberRequest{
				AdministrativeRoleId: testRoleID,
				AssigneeId:           testAssigneeID,
				ScopeType:            AdministrativeRoleScopeGlobal,
			},
			Action: AdministrativeRoleMemberActionAdd,
		},
		{
			AdministrativeRoleMemberRequest: AdministrativeRoleMemberRequest{
				AdministrativeRoleId: BuiltinRoleVaultUserId,
				AssigneeId:           testAssigneeID,
				ScopeType:            AdministrativeRoleScopeVault,
				ScopeResourceId:      &scopeResourceId,
			},
			Action: AdministrativeRoleMemberActionDelete,
		},
	})
	assert.NoError(t, err)
}

func TestRoleAssignmentsRemoveMember(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/members/assignment-1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})

	client := newTestClient(t, mux)

	err := client.AdministrativeRoleAssignments.RemoveMember("assignment-1")
	assert.NoError(t, err)
}

func TestRoleAssignmentsDeleteScope(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/administrative-role-assignments/scope", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		assert.Equal(t, testRoleID, r.URL.Query().Get("roleId"))
		assert.Equal(t, "2", r.URL.Query().Get("scopeType"))
		assert.Equal(t, testVaultID, r.URL.Query().Get("scopeResourceId"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":1}`))
	})

	client := newTestClient(t, mux)

	err := client.AdministrativeRoleAssignments.DeleteScope(testRoleID, AdministrativeRoleScopeVault, testVaultID)
	assert.NoError(t, err)
}
