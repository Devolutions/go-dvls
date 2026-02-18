package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EntryCertificateService service

// EntryCertificate represents a certificate entry.
type EntryCertificate struct {
	Id                    string
	VaultId               string
	Name                  string
	Description           string
	EntryFolderPath       string
	Tags                  []string
	Expiration            time.Time
	Password              string
	UseDefaultCredentials bool

	// Can either be a URL or a file name.
	CertificateIdentifier string

	data entryCertificateData
}

type rawEntryCertificate struct {
	Id              string      `json:"id,omitempty"`
	VaultId         string      `json:"repositoryId"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	EntryFolderPath string      `json:"group"`
	ModifiedDate    *ServerTime `json:"modifiedDate,omitempty"`
	Tags            string      `json:"keywords,omitempty"`
	Expiration      time.Time   `json:"expiration,omitempty"`

	ConnectionType    ServerConnectionType    `json:"connectionType"`    // 45 - document
	ConnectionSubType ServerConnectionSubType `json:"connectionSubType"` // "Certificate"

	Data entryCertificateData `json:"data"`
}

type entryCertificateData struct {
	Mode                  int    `json:"dataMode"`     // 3 - URL, 2 - File
	FileSize              int    `json:"documentSize"` // 0 on mode 3
	FileName              string `json:"fileName"`
	Type                  any    `json:"type"` // "Certificate"
	UseDefaultCredentials bool   `json:"useWebDefaultCredentials"`
	Password              struct {
		HasSensitiveData bool   `json:"hasSensitiveData"`
		SensitiveData    string `json:"sensitiveData"`
	} `json:"password"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e EntryCertificate) MarshalJSON() ([]byte, error) {
	raw := rawEntryCertificate{
		Id:              e.Id,
		VaultId:         e.VaultId,
		Name:            e.Name,
		Description:     e.Description,
		EntryFolderPath: e.EntryFolderPath,
		Tags:            sliceToKeywords(e.Tags),
		Expiration:      e.Expiration,
		Data: entryCertificateData{
			Mode:                  e.data.Mode,
			FileName:              e.CertificateIdentifier,
			Type:                  "Certificate",
			UseDefaultCredentials: e.UseDefaultCredentials,
			FileSize:              e.data.FileSize,
			Password: struct {
				HasSensitiveData bool   `json:"hasSensitiveData"`
				SensitiveData    string `json:"sensitiveData"`
			}{
				HasSensitiveData: true,
				SensitiveData:    e.Password,
			},
		},
	}

	raw.ConnectionType = ServerConnectionDocument
	raw.ConnectionSubType = ServerConnectionSubTypeCertificate

	entryJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return entryJson, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *EntryCertificate) UnmarshalJSON(d []byte) error {
	rawString := struct {
		Data string
	}{}
	err := json.Unmarshal(d, &rawString)
	if err != nil && !strings.Contains(err.Error(), "cannot unmarshal object into Go struct") {
		return err
	}

	raw := struct {
		Data rawEntryCertificate
	}{}

	if rawString.Data != "" {
		err = json.Unmarshal([]byte(rawString.Data), &raw.Data)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(d, &raw)
		if err != nil {
			return err
		}
	}

	e.Id = raw.Data.Id
	e.VaultId = raw.Data.VaultId
	e.Name = raw.Data.Name
	e.Description = raw.Data.Description
	e.EntryFolderPath = raw.Data.EntryFolderPath
	e.Tags = keywordsToSlice(raw.Data.Tags)
	e.Expiration = raw.Data.Expiration

	e.data.Mode = raw.Data.Data.Mode
	e.CertificateIdentifier = raw.Data.Data.FileName
	e.UseDefaultCredentials = raw.Data.Data.UseDefaultCredentials
	e.Password = raw.Data.Data.Password.SensitiveData
	e.data.FileSize = raw.Data.Data.FileSize

	return nil
}

// Get returns a single Certificate specified by entryId.
func (c *EntryCertificateService) Get(entryId string) (EntryCertificate, error) {
	return c.GetWithContext(context.Background(), entryId)
}

// GetWithContext returns a single Certificate specified by entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetWithContext(ctx context.Context, entryId string) (EntryCertificate, error) {
	var entry EntryCertificate
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching entry: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return entry, nil
}

// GetFileContent returns the content of the file specified by entryId.
func (c *EntryCertificateService) GetFileContent(entryId string) ([]byte, error) {
	return c.GetFileContentWithContext(context.Background(), entryId)
}

// GetFileContentWithContext returns the content of the file specified by entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetFileContentWithContext(ctx context.Context, entryId string) ([]byte, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryConnectionsEndpoint, entryId, "document")
	if err != nil {
		return nil, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching entry content: %w", err)
	}

	return resp.Response, nil
}

// GetPassword returns the password of the entry specified by entry.
func (c *EntryCertificateService) GetPassword(entry EntryCertificate) (EntryCertificate, error) {
	return c.GetPasswordWithContext(context.Background(), entry)
}

// GetPasswordWithContext returns the password of the entry specified by entry.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetPasswordWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	var entryPassword EntryCertificate
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entry.Id, "/sensitive-data")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, nil)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching sensitive data: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entryPassword)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	entry.Password = entryPassword.Password

	return entry, nil
}

// NewURL creates a new EntryCertificate based on entry. Will use the url as the file content.
func (c *EntryCertificateService) NewURL(entry EntryCertificate) (EntryCertificate, error) {
	return c.NewURLWithContext(context.Background(), entry)
}

// NewURLWithContext creates a new EntryCertificate based on entry. Will use the url as the file content.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) NewURLWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	return c.newWithContext(ctx, entry, nil)
}

// NewFile creates a new EntryCertificate based on entry. Will upload the file content to the DVLS server.
func (c *EntryCertificateService) NewFile(entry EntryCertificate, content []byte) (EntryCertificate, error) {
	return c.NewFileWithContext(context.Background(), entry, content)
}

// NewFileWithContext creates a new EntryCertificate based on entry. Will upload the file content to the DVLS server.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) NewFileWithContext(ctx context.Context, entry EntryCertificate, content []byte) (EntryCertificate, error) {
	return c.newWithContext(ctx, entry, content)
}

func (c *EntryCertificateService) newWithContext(ctx context.Context, entry EntryCertificate, content []byte) (EntryCertificate, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	entry.data.Mode = 3

	if content != nil {
		entry.data.Mode = 2
		entry.data.FileSize = len(content)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while creating entry: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if content != nil {
		attachment := EntryAttachment{
			EntryId:   entry.Id,
			FileName:  entry.CertificateIdentifier,
			Size:      len(content),
			IsPrivate: true,
		}

		entryAttachment, err := c.client.newAttachmentRequest(ctx, attachment)
		if err != nil {
			return EntryCertificate{}, fmt.Errorf("error while creating entry attachment: %w", err)
		}

		err = c.client.uploadAttachment(ctx, content, entryAttachment)
		if err != nil {
			return EntryCertificate{}, fmt.Errorf("error while uploading attachment: %w", err)
		}
	}

	return entry, nil
}

// Update updates an EntryCertificate based on entry. Will replace all other fields whether included or not.
func (c *EntryCertificateService) Update(entry EntryCertificate) (EntryCertificate, error) {
	return c.UpdateWithContext(context.Background(), entry)
}

// UpdateWithContext updates an EntryCertificate based on entry. Will replace all other fields whether included or not.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) UpdateWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	oldEntry, err := c.GetWithContext(ctx, entry.Id)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching entry: %w", err)
	}

	entry.data.Mode = oldEntry.data.Mode
	entry.data.FileSize = oldEntry.data.FileSize

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPut, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while creating entry: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return entry, nil
}

// Delete deletes an EntryCertificate based on entryId.
func (c *EntryCertificateService) Delete(entryId string) error {
	return c.DeleteWithContext(context.Background(), entryId)
}

// DeleteWithContext deletes an EntryCertificate based on entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) DeleteWithContext(ctx context.Context, entryId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return fmt.Errorf("failed to delete entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// GetDataMode returns the data mode of the EntryCertificate. Can be either EntryCertificateDataModeURL or EntryCertificateDataModeFile.
func (c EntryCertificate) GetDataMode() EntryCertificateDataMode {
	return EntryCertificateDataMode(c.data.Mode)
}
