package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetchAllPages_StopsAtTotalPage covers termination on the totalPage bound when
// totalCount alone would not stop the loop.
func TestFetchAllPages_StopsAtTotalPage(t *testing.T) {
	var pagesRequested []int

	err := fetchAllPages(func(pageNumber int) (pagedResponse, int, error) {
		pagesRequested = append(pagesRequested, pageNumber)
		return pagedResponse{TotalPage: 3, TotalCount: 10}, 2, nil
	})

	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, pagesRequested)
}

// TestFetchAllPages_StopsAtTotalCount covers a server that reports totalPage > 1 but
// whose totalCount is already satisfied by the first page: the driver must stop
// without requesting further pages.
func TestFetchAllPages_StopsAtTotalCount(t *testing.T) {
	calls := 0

	err := fetchAllPages(func(pageNumber int) (pagedResponse, int, error) {
		calls++
		return pagedResponse{CurrentPage: pageNumber, TotalPage: 3, TotalCount: 2}, 2, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

// TestFetchAllPages_EmptyFirstPage covers an empty result set: one request, no error.
func TestFetchAllPages_EmptyFirstPage(t *testing.T) {
	calls := 0

	err := fetchAllPages(func(pageNumber int) (pagedResponse, int, error) {
		calls++
		return pagedResponse{}, 0, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}
