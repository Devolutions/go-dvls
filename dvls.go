package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const apiV2Prefix string = "/api/v2"

// Response represents an HTTP response from the DVLS API. Response holds the raw response body.
type Response struct {
	Response []byte `json:"-"`
}

// RequestError is returned when a request fails or the server answers with a non-2xx status code.
// Body holds the raw response body; Error includes the message of a DVLS error body.
type RequestError struct {
	Url        string
	StatusCode int
	Body       []byte
	Err        error
}

const defaultContentType string = "application/json"

// RequestOptions contains optional settings for Request and RequestWithContext.
type RequestOptions struct {
	// ContentType replaces the default application/json Content-Type header.
	ContentType string
}

type apiErrorBody struct {
	Message string `json:"message"`
	Errors  []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (e RequestError) Error() string {
	msg := e.Err.Error()
	if detail := errorBodyDetail(e.Body); detail != "" {
		msg += ": " + detail
	}

	if e.StatusCode != 0 {
		return fmt.Sprintf("error while submitting request on url %s (status %d). error: %s", e.Url, e.StatusCode, msg)
	}

	return fmt.Sprintf("error while submitting request on url %s. error: %s", e.Url, msg)
}

func (e RequestError) Unwrap() error {
	return e.Err
}

// errorBodyDetail returns the messages of a DVLS error body ({message} or {errors:[{field,message}]}), or an
// empty string when the body is not one.
func errorBodyDetail(body []byte) string {
	var apiErr apiErrorBody
	if len(body) == 0 || json.Unmarshal(body, &apiErr) != nil {
		return ""
	}

	parts := make([]string, 0, len(apiErr.Errors)+1)
	if apiErr.Message != "" {
		parts = append(parts, apiErr.Message)
	}
	for _, fieldErr := range apiErr.Errors {
		if fieldErr.Field != "" {
			parts = append(parts, fieldErr.Field+": "+fieldErr.Message)
		} else {
			parts = append(parts, fieldErr.Message)
		}
	}

	return strings.Join(parts, "; ")
}

// IsNotFound reports whether the error is a DVLS RequestError with an HTTP 404 status code.
func IsNotFound(err error) bool {
	return hasStatusCode(err, http.StatusNotFound)
}

func hasStatusCode(err error, statusCode int) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode == statusCode
	}

	return false
}

// Request returns a Response that contains the HTTP response body.
// For context support, use RequestWithContext instead.
func (c *Client) Request(url string, reqMethod string, reqBody io.Reader, options ...RequestOptions) (Response, error) {
	return c.RequestWithContext(context.Background(), url, reqMethod, reqBody, options...)
}

// RequestWithContext returns a Response that contains the HTTP response body.
// When the session token was rejected (HTTP 401), the client logs in again and retries the request once.
// The provided context can be used to cancel the request.
func (c *Client) RequestWithContext(ctx context.Context, url string, reqMethod string, reqBody io.Reader, options ...RequestOptions) (Response, error) {
	contentType := defaultContentType
	if len(options) > 0 && options[0].ContentType != "" {
		contentType = options[0].ContentType
	}

	var body []byte
	if reqBody != nil {
		var err error
		body, err = io.ReadAll(reqBody)
		if err != nil {
			return Response{}, &RequestError{Err: fmt.Errorf("failed to read request body: %w", err), Url: url}
		}
	}

	token := c.credential.currentToken()
	resp, err := c.rawRequestWithContext(ctx, url, reqMethod, contentType, token, body)
	if err == nil || c.credential.useApiKey() || !hasStatusCode(err, http.StatusUnauthorized) {
		return resp, err
	}

	// DVLS tokens expire a fixed time after login and the API has no refresh operation.
	err = c.reloginWithContext(ctx, token)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to refresh login token: %w", err), Url: url}
	}

	return c.rawRequestWithContext(ctx, url, reqMethod, contentType, c.credential.currentToken(), body)
}

func (c *Client) rawRequestWithContext(ctx context.Context, url string, reqMethod string, contentType string, token string, body []byte) (Response, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, reqMethod, url, reqBody)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to make request: %w", err), Url: url}
	}

	req.Header.Set("Content-Type", contentType)
	if c.credential.useApiKey() {
		req.Header.Set("Authorization", "Bearer "+c.credential.apiKey)
	} else if token != "" {
		req.Header.Set("tokenId", token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("error while submitting request: %w", err), Url: url}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to read response body: %w", err), Url: url, StatusCode: resp.StatusCode}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Response{}, &RequestError{Err: fmt.Errorf("unexpected status code %d", resp.StatusCode), Url: url, StatusCode: resp.StatusCode, Body: respBody}
	}

	return Response{Response: respBody}, nil
}

// v2Url returns the URL of a v2 API route built from path segments.
func (c *Client) v2Url(query url.Values, segments ...string) (string, error) {
	reqUrl, err := url.JoinPath(c.baseUri, append([]string{apiV2Prefix}, segments...)...)
	if err != nil {
		return "", err
	}

	if len(query) > 0 {
		reqUrl += "?" + query.Encode()
	}

	return reqUrl, nil
}

// getJSON performs a GET on reqUrl and unmarshals the response body into out.
func (c *Client) getJSON(ctx context.Context, reqUrl string, out any) error {
	return c.sendJSON(ctx, http.MethodGet, reqUrl, nil, out)
}

// sendJSON sends in as the JSON body (none when in is nil) and unmarshals the response body into out when out is
// not nil.
func (c *Client) sendJSON(ctx context.Context, method string, reqUrl string, in any, out any) error {
	var reqBody io.Reader
	if in != nil {
		body, err := jsonBody(in)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(body)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, method, reqBody)
	if err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	return unmarshalBody(resp.Response, out)
}
