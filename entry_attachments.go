package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EntryAttachment struct {
	ID            string `json:"id,omitempty"`
	IDString      string `json:"idString"`
	EntryID       string `json:"connectionID"`
	EntryIDString string `json:"connectionIDString"`
	Description   string `json:"description"`
	FileName      string `json:"filename"`
	IsPrivate     bool   `json:"isPrivate"`
	Size          int    `json:"size"`
	Title         string `json:"title"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *EntryAttachment) UnmarshalJSON(d []byte) error {
	type rawEntryAttachment EntryAttachment
	raw := struct {
		Data rawEntryAttachment `json:"data"`
	}{}

	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	*e = EntryAttachment(raw.Data)

	return nil
}

const attachmentEndpoint = "/api/attachment"

func (c *Client) newAttachmentRequest(attachment EntryAttachment) (string, error) {
	reqUrl, err := url.JoinPath(c.baseUri, attachmentEndpoint, "save?=&private=false&useSensitiveMode=true")
	if err != nil {
		return "", fmt.Errorf("failed to build attachment url. error: %w", err)
	}

	reqUrl, err = url.QueryUnescape(reqUrl)
	if err != nil {
		return "", fmt.Errorf("failed to unescape query url. error: %w", err)
	}

	entryJson, err := json.Marshal(attachment)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return "", fmt.Errorf("error while submitting entry attachment request. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return "", err
	}

	err = json.Unmarshal(resp.Response, &attachment)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return attachment.ID, nil
}

func (c *Client) uploadAttachment(fileBytes []byte, attachmentId string) error {
	reqUrl, err := url.JoinPath(c.baseUri, attachmentEndpoint, attachmentId, "document")
	if err != nil {
		return fmt.Errorf("failed to build attachment url. error: %w", err)
	}

	contentType := http.DetectContentType(fileBytes)

	resp, err := c.Request(reqUrl, http.MethodPost, bytes.NewBuffer(fileBytes), RequestOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("error while uploading entry attachment. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}
