package dvls

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAll_Pagination(t *testing.T) {
	var pages []string

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user-group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "500", r.URL.Query().Get("pageSize"))
		pages = append(pages, r.URL.Query().Get("pageNumber"))

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("pageNumber") {
		case "1":
			w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":3,"totalPage":2,"data":[{"id":"1","name":"a"},{"id":"2","name":"b"}]}`))
		case "2":
			w.Write([]byte(`{"currentPage":2,"pageSize":500,"totalCount":3,"totalPage":2,"data":[{"id":"3","name":"c"}]}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errors":[{"field":"pageNumber","message":"pageNumber is past the last page."}]}`))
		}
	})

	client := newTestClient(t, mux)

	groups, err := listAll[UserGroup](context.Background(), client, nil, userGroupEndpoint)
	require.NoError(t, err)
	assert.Len(t, groups, 3)
	assert.Equal(t, []string{"1", "2"}, pages)
}

func TestListAll_Empty(t *testing.T) {
	calls := 0

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user-group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"currentPage":1,"pageSize":500,"totalCount":0,"totalPage":0,"data":[]}`))
	})

	client := newTestClient(t, mux)

	groups, err := listAll[UserGroup](context.Background(), client, nil, userGroupEndpoint)
	require.NoError(t, err)
	assert.Empty(t, groups)
	assert.Equal(t, 1, calls)
}

func TestListAll_PageError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/user-group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"Access denied."}`))
	})

	client := newTestClient(t, mux)

	_, err := listAll[UserGroup](context.Background(), client, nil, userGroupEndpoint)
	require.Error(t, err)
	assert.ErrorContains(t, err, "page 1")
	assert.ErrorContains(t, err, "Access denied.")
}
