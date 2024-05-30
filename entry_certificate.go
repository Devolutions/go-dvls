package dvls

import (
	"bytes"
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
	ID                    string
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
	ID              string      `json:"id,omitempty"`
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
		ID:              e.ID,
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

	e.ID = raw.Data.ID
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
	var entry EntryCertificate
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return entry, nil
}

// GetFileContent returns the content of the file specified by entryId.
func (c *EntryCertificateService) GetFileContent(entryId string) ([]byte, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryConnectionsEndpoint, entryId, "document")
	if err != nil {
		return nil, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodGet, nil, RequestOptions{RawBody: true})
	if err != nil {
		return nil, fmt.Errorf("error while fetching entry content. error: %w", err)
	}

	return resp.Response, nil
}

// GetPassword returns the password of the entry specified by entry.
func (c *EntryCertificateService) GetPassword(entry EntryCertificate) (EntryCertificate, error) {
	var entryPassword EntryCertificate
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entry.ID, "/sensitive-data")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entryPassword)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	entry.Password = entryPassword.Password

	return entry, nil
}

// NewURL creates a new EntryCertificate based on entry. Will use the url as the file content.
func (c *EntryCertificateService) NewURL(entry EntryCertificate) (EntryCertificate, error) {
	return c.new(entry, nil)
}

// NewFile creates a new EntryCertificate based on entry. Will upload the file content to the DVLS server.
func (c *EntryCertificateService) NewFile(entry EntryCertificate, content []byte) (EntryCertificate, error) {
	return c.new(entry, content)
}

func (c *EntryCertificateService) new(entry EntryCertificate, content []byte) (EntryCertificate, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entry.data.Mode = 3

	if content != nil {
		entry.data.Mode = 2
		entry.data.FileSize = len(content)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while creating entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	if content != nil {
		attachment := EntryAttachment{
			EntryID:   entry.ID,
			FileName:  entry.CertificateIdentifier,
			Size:      len(content),
			IsPrivate: true,
		}

		entryAttachment, err := c.client.newAttachmentRequest(attachment)
		if err != nil {
			return EntryCertificate{}, fmt.Errorf("error while creating entry attachment. error: %w", err)
		}

		err = c.client.uploadAttachment(content, entryAttachment)
		if err != nil {
			return EntryCertificate{}, fmt.Errorf("error while uploading attachment. error: %w", err)
		}
	}

	return entry, nil
}

// Update updates an EntryCertificate based on entry. Will replace all other fields whether included or not.
func (c *EntryCertificateService) Update(entry EntryCertificate) (EntryCertificate, error) {
	_, err := c.Get(entry.ID)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while fetching entry. error: %w", err)
	}

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPut, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("error while creating entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryCertificate{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return entry, nil
}

// Delete deletes an EntryCertificate based on entryId.
func (c *EntryCertificateService) Delete(entryId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return fmt.Errorf("failed to delete entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}
