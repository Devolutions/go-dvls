package dvls

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Response represents an HTTP response from the DVLS API. Contains the response body in bytes, the result code
// and the result message.
type Response struct {
	Response []byte `json:"-"`
	Result   uint8
	Message  string
}

type RequestError struct {
	Url        string
	StatusCode int
	Body       []byte
	Err        error
}

const defaultContentType string = "application/json"

type RequestOptions struct {
	ContentType string
	RawBody     bool
}

func (e RequestError) Error() string {
	if e.StatusCode != 0 {
		return fmt.Sprintf("error while submitting request on url %s (status %d). error: %s", e.Url, e.StatusCode, e.Err.Error())
	}

	return fmt.Sprintf("error while submitting request on url %s. error: %s", e.Url, e.Err.Error())
}

// IsNotFound reports whether the error is a DVLS RequestError with an HTTP 404 status code.
func IsNotFound(err error) bool {
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode == http.StatusNotFound
	}

	return false
}

// Request returns a Response that contains the HTTP response body in bytes, the result code and result message.
func (c *Client) Request(url string, reqMethod string, reqBody io.Reader, options ...RequestOptions) (Response, error) {
	islogged, err := c.isLogged()
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to fetch login status. error: %w", err), Url: url}
	}
	if !islogged {
		err := c.login()
		if err != nil {
			return Response{}, &RequestError{Err: fmt.Errorf("failed to refresh login token. error: %w", err), Url: url}
		}
	}

	var opts RequestOptions
	if len(options) > 0 {
		opts = options[0]
	}

	resp, err := c.rawRequest(url, reqMethod, defaultContentType, reqBody, opts)
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}

func (c *Client) rawRequest(url string, reqMethod string, contentType string, reqBody io.Reader, options ...RequestOptions) (Response, error) {
	var opts RequestOptions
	if len(options) > 0 {
		opts = options[0]
	}

	req, err := http.NewRequest(reqMethod, url, reqBody)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to make request. error: %w", err), Url: url}
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("tokenId", c.credential.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("error while submitting request. error: %w", err), Url: url}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)

		return Response{}, &RequestError{Err: fmt.Errorf("unexpected status code %d", resp.StatusCode), Url: url, StatusCode: resp.StatusCode, Body: body}
	}

	var response Response
	response.Response, err = io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to read response body. error: %w", err), Url: url}
	}

	if !opts.RawBody && len(response.Response) > 0 {
		err = json.Unmarshal(response.Response, &response)
		if err != nil {
			return response, &RequestError{Err: fmt.Errorf("failed to unmarshal response body. error: %w", err), Url: url}
		}
	}

	return response, nil
}

func (r Response) CheckRespSaveResult() error {
	resultCode := SaveResult(r.Result)
	if resultCode != SaveResultSuccess {
		return fmt.Errorf("unexpected result code %d (%s) %s", resultCode, resultCode, r.Message)
	}
	return nil
}
