package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EntryAttachment struct {
	Id            string `json:"id,omitempty"`
	IdString      string `json:"idString"`
	EntryId       string `json:"connectionID"`
	EntryIdString string `json:"connectionIDString"`
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

func (c *Client) newAttachmentRequest(ctx context.Context, attachment EntryAttachment) (string, error) {
	reqUrl, err := url.JoinPath(c.baseUri, attachmentEndpoint, "save?=&private=false&useSensitiveMode=true")
	if err != nil {
		return "", fmt.Errorf("failed to build attachment url: %w", err)
	}

	reqUrl, err = url.QueryUnescape(reqUrl)
	if err != nil {
		return "", fmt.Errorf("failed to unescape query url: %w", err)
	}

	entryJson, err := json.Marshal(attachment)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return "", fmt.Errorf("error while submitting entry attachment request: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return "", err
	}

	err = json.Unmarshal(resp.Response, &attachment)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return attachment.Id, nil
}

func (c *Client) uploadAttachment(ctx context.Context, fileBytes []byte, attachmentId string) error {
	reqUrl, err := url.JoinPath(c.baseUri, attachmentEndpoint, attachmentId, "document")
	if err != nil {
		return fmt.Errorf("failed to build attachment url: %w", err)
	}

	contentType := http.DetectContentType(fileBytes)

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(fileBytes), RequestOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("error while uploading entry attachment: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}
