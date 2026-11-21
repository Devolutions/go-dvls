package dvls

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const listPageSize int = 500

type pagedResponse[T any] struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalCount  int `json:"totalCount"`
	TotalPage   int `json:"totalPage"`
	Data        []T `json:"data"`
}

// listAll fetches every page of a v2 list route and returns the items of all pages.
func listAll[T any](ctx context.Context, c *Client, query url.Values, segments ...string) ([]T, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("pageSize", strconv.Itoa(listPageSize))

	var items []T
	for pageNumber := 1; ; pageNumber++ {
		query.Set("pageNumber", strconv.Itoa(pageNumber))

		reqUrl, err := c.v2Url(query, segments...)
		if err != nil {
			return nil, fmt.Errorf("failed to build url: %w", err)
		}

		var page pagedResponse[T]
		err = c.getJSON(ctx, reqUrl, &page)
		if err != nil {
			return nil, fmt.Errorf("page %d: %w", pageNumber, err)
		}

		items = append(items, page.Data...)

		// A page past the last one is a 400, so stop on the paging metadata.
		if len(page.Data) == 0 || pageNumber >= page.TotalPage || len(items) >= page.TotalCount {
			return items, nil
		}
	}
}
