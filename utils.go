package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func keywordsToSlice(kw string) []string {
	var spacedTag bool
	tags := strings.FieldsFunc(string(kw), func(r rune) bool {
		if r == '"' {
			spacedTag = !spacedTag
		}
		return !spacedTag && r == ' '
	})
	for i, v := range tags {
		unquotedTag, err := strconv.Unquote(v)
		if err != nil {
			continue
		}

		tags[i] = unquotedTag
	}

	return tags
}

func sliceToKeywords(kw []string) string {
	keywords := []string(kw)
	for i, v := range keywords {
		if strings.Contains(v, " ") {
			kw[i] = "\"" + v + "\""
		}
	}

	kString := strings.Join(keywords, " ")

	return kString
}

type dataListResponse[T any] struct {
	Result  SaveResult `json:"result"`
	Message string     `json:"message"`
	Data    []T        `json:"data"`
}

// fetchDataList performs a GET on an unpaged endpoint that wraps its items in
// a {result, data: [...]} envelope and returns the items.
func fetchDataList[T any](ctx context.Context, c *Client, endpoint string, resource string) ([]T, error) {
	reqUrl, err := url.JoinPath(c.baseUri, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s url: %w", resource, err)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching %s: %w", resource, err)
	}

	var listResp dataListResponse[T]
	if err := json.Unmarshal(resp.Response, &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	if listResp.Result != SaveResultSuccess {
		return nil, fmt.Errorf("unexpected result code %d (%s) %s", listResp.Result, listResp.Result, listResp.Message)
	}

	return listResp.Data, nil
}

// singleByName returns the only item whose name equals name, or notFound /
// multiple when zero or several items match.
func singleByName[T any](items []T, name string, nameOf func(T) string, notFound error, multiple error) (T, error) {
	var match T
	count := 0
	for _, item := range items {
		if nameOf(item) == name {
			match = item
			count++
		}
	}

	switch count {
	case 0:
		return match, notFound
	case 1:
		return match, nil
	default:
		var zero T
		return zero, multiple
	}
}

// jsonValuePresent reports whether a raw JSON value holds anything other than
// nothing or null.
func jsonValuePresent(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && !bytes.Equal(trimmed, []byte("null"))
}
