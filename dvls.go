package dvls

import (
	"encoding/json"
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
	Url string
	Err error
}

const defaultContentType string = "application/json"

func (e RequestError) Error() string {
	return fmt.Sprintf("error while submitting request on url %s. error: %s", e.Url, e.Err.Error())
}

// Request returns a Response that contains the HTTP response body in bytes, the result code and result message.
func (c *Client) Request(url string, reqMethod string, reqBody io.Reader) (Response, error) {
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

	resp, err := c.rawRequest(url, reqMethod, defaultContentType, reqBody)
	if err != nil {
		return Response{}, err
	}
	return resp, nil
}

func (c *Client) rawRequest(url string, reqMethod string, contentType string, reqBody io.Reader) (Response, error) {
	req, err := http.NewRequest(reqMethod, url, reqBody)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to make request. error: %w", err), Url: url}
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("tokenId", c.credential.token)

	// TODO: remove this on next server release
	req.Header.Add("X-Forwarded-For", "null")

	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("error while submitting request. error: %w", err), Url: url}
	} else if resp.StatusCode != http.StatusOK {
		return Response{}, &RequestError{Err: fmt.Errorf("unexpected status code %s", resp.Status), Url: url}
	}
	var response Response
	response.Response, err = io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, &RequestError{Err: fmt.Errorf("failed to read response body. error: %w", err), Url: url}
	}
	defer resp.Body.Close()

	err = json.Unmarshal(response.Response, &response)
	if err != nil {
		return response, &RequestError{Err: fmt.Errorf("failed to unmarshall response body. error: %w", err), Url: url}
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
