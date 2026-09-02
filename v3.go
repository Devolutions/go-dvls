package dvls

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// v3ResultError represents a logical failure reported in a v3 response
// envelope.
type v3ResultError struct {
	Result       SaveResult
	ErrorMessage string
}

func (e *v3ResultError) Error() string {
	return fmt.Sprintf("unexpected result code %d (%s) %s", e.Result, e.Result, e.ErrorMessage)
}

// v3Envelope represents the DVLS v3 API response envelope.
// Depending on the authentication method, mutation responses are returned
// unwrapped (empty body on success) while GET responses are always wrapped in
// this envelope. The result field is omitted by the server when it equals 0
// (SaveResultError), so a missing result must never be interpreted as success.
type v3Envelope struct {
	Result       SaveResult      `json:"result"`
	ErrorMessage string          `json:"errorMessage"`
	Data         json.RawMessage `json:"data"`
}

func (e v3Envelope) checkResult() error {
	if e.Result != SaveResultSuccess {
		return &v3ResultError{Result: e.Result, ErrorMessage: e.ErrorMessage}
	}

	return nil
}

// unwrapV3Data validates a wrapped v3 response and returns its data payload.
// GET responses are wrapped for every authentication method; only mutation
// responses (checkV3SaveResponse) and the paged list (decodeV3PagedResponse)
// vary in shape.
func unwrapV3Data(body []byte) (json.RawMessage, error) {
	var envelope v3Envelope
	err := json.Unmarshal(body, &envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if err := envelope.checkResult(); err != nil {
		return nil, err
	}

	return envelope.Data, nil
}

// decodeV3 validates a wrapped v3 response and unmarshals its data payload.
func decodeV3[T any](body []byte) (T, error) {
	var value T

	data, err := unwrapV3Data(body)
	if err != nil {
		return value, err
	}

	err = json.Unmarshal(data, &value)
	if err != nil {
		return value, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return value, nil
}

// isV3ResultNotFound reports whether the error is a v3ResultError with a
// SaveResultNotFound result code.
func isV3ResultNotFound(err error) bool {
	var resultErr *v3ResultError
	return errors.As(err, &resultErr) && resultErr.Result == SaveResultNotFound
}

// checkV3SaveResponse validates a v3 mutation response. An empty or null body
// on HTTP 200 is a success (sessions using HTTP status codes report failures
// with 4xx/5xx responses); otherwise the envelope result must be
// SaveResultSuccess.
func checkV3SaveResponse(resp Response) error {
	if !jsonValuePresent(resp.Response) {
		return nil
	}

	_, err := unwrapV3Data(resp.Response)
	return err
}

type v3PagedBody struct {
	v3Envelope
	pagedResponse
}

// decodeV3PagedResponse decodes a v3 paged list response and returns the raw
// page items along with the paging metadata. It supports both response
// shapes: the paged payload at the root (data is the item array) and the
// paged payload wrapped in the envelope (data is an object).
func decodeV3PagedResponse(body []byte) (json.RawMessage, pagedResponse, error) {
	var root v3PagedBody
	err := json.Unmarshal(body, &root)
	if err != nil {
		return nil, pagedResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	data := bytes.TrimSpace(root.Data)
	if len(data) > 0 && data[0] == '[' {
		return root.Data, root.pagedResponse, nil
	}

	if err := root.checkResult(); err != nil {
		return nil, pagedResponse{}, err
	}

	if !jsonValuePresent(data) {
		return nil, pagedResponse{}, nil
	}

	var paged struct {
		Data json.RawMessage `json:"data"`
		pagedResponse
	}
	err = json.Unmarshal(root.Data, &paged)
	if err != nil {
		return nil, pagedResponse{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return paged.Data, paged.pagedResponse, nil
}
