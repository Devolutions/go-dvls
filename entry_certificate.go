package dvls

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type EntryCertificateService service

const (
	entryDocumentType           string = "Document"
	entryCertificateSubType     string = "Certificate"
	entryCertificateContentPath string = "content"
)

type EntryCertificateDataMode string

const (
	EntryCertificateDataModeURL  EntryCertificateDataMode = "Url"
	EntryCertificateDataModeFile EntryCertificateDataMode = "File"
)

// EntryCertificate represents a certificate entry (type Document, subtype Certificate).
type EntryCertificate struct {
	Id              string
	VaultId         string
	Name            string
	Description     string
	EntryFolderPath string
	Tags            []string
	Password        string

	// Can either be a URL or a file name.
	CertificateIdentifier string

	mode EntryCertificateDataMode
}

type entryCertificateData struct {
	DataMode    EntryCertificateDataMode `json:"dataMode"`
	Url         string                   `json:"url,omitempty"`
	FileName    string                   `json:"fileName,omitempty"`
	FileContent string                   `json:"fileContent,omitempty"`
	Password    string                   `json:"password,omitempty"`
}

func certificateFromRaw(raw rawEntry) (EntryCertificate, error) {
	if raw.Type != entryDocumentType || raw.SubType != entryCertificateSubType {
		return EntryCertificate{}, ErrUnsupportedEntryType{Type: raw.Type, SubType: raw.SubType}
	}

	var data entryCertificateData
	if len(raw.Data) > 0 {
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return EntryCertificate{}, fmt.Errorf("failed to unmarshal certificate data: %w", err)
		}
	}

	identifier := data.Url
	if data.DataMode == EntryCertificateDataModeFile {
		identifier = data.FileName
	}

	return EntryCertificate{
		Id:                    raw.Id,
		VaultId:               raw.VaultId,
		Name:                  raw.Name,
		Description:           raw.Description,
		EntryFolderPath:       raw.Path,
		Tags:                  raw.Tags,
		Password:              data.Password,
		CertificateIdentifier: identifier,
		mode:                  data.DataMode,
	}, nil
}

func (e EntryCertificate) data(mode EntryCertificateDataMode, content []byte) entryCertificateData {
	data := entryCertificateData{DataMode: mode, Password: e.Password}
	if mode == EntryCertificateDataModeFile {
		data.FileName = e.CertificateIdentifier
	} else {
		data.Url = e.CertificateIdentifier
	}

	if content != nil {
		data.FileContent = base64.StdEncoding.EncodeToString(content)
	}

	return data
}

// Get returns a single Certificate, its password included, specified by vaultId and entryId.
func (c *EntryCertificateService) Get(vaultId string, entryId string) (EntryCertificate, error) {
	return c.GetWithContext(context.Background(), vaultId, entryId)
}

// GetWithContext returns a single Certificate, its password included, specified by vaultId and entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetWithContext(ctx context.Context, vaultId string, entryId string) (EntryCertificate, error) {
	raw, err := c.client.getRawEntry(ctx, vaultId, entryId, nil)
	if err != nil {
		return EntryCertificate{}, err
	}

	return certificateFromRaw(raw)
}

// GetFileContent returns the content of the file of a certificate kept in file mode, specified by vaultId and entryId.
func (c *EntryCertificateService) GetFileContent(vaultId string, entryId string) ([]byte, error) {
	return c.GetFileContentWithContext(context.Background(), vaultId, entryId)
}

// GetFileContentWithContext returns the content of the file of a certificate kept in file mode, specified by vaultId
// and entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetFileContentWithContext(ctx context.Context, vaultId string, entryId string) ([]byte, error) {
	if vaultId == "" || entryId == "" {
		return nil, fmt.Errorf("both entry Id and vault Id are required")
	}

	reqUrl, err := entryUrl(c.client, nil, vaultId, entryId, entryCertificateContentPath)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching entry content: %w", err)
	}

	return resp.Response, nil
}

// GetPassword returns entry with its Password read from the server.
func (c *EntryCertificateService) GetPassword(entry EntryCertificate) (EntryCertificate, error) {
	return c.GetPasswordWithContext(context.Background(), entry)
}

// GetPasswordWithContext returns entry with its Password read from the server.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) GetPasswordWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	stored, err := c.GetWithContext(ctx, entry.VaultId, entry.Id)
	if err != nil {
		return EntryCertificate{}, err
	}

	entry.Password = stored.Password

	return entry, nil
}

// NewURL creates a new EntryCertificate based on entry, with CertificateIdentifier as the certificate's URL.
func (c *EntryCertificateService) NewURL(entry EntryCertificate) (EntryCertificate, error) {
	return c.NewURLWithContext(context.Background(), entry)
}

// NewURLWithContext creates a new EntryCertificate based on entry, with CertificateIdentifier as the certificate's URL.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) NewURLWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	return c.newWithContext(ctx, entry, EntryCertificateDataModeURL, nil)
}

// NewFile creates a new EntryCertificate based on entry and stores content on the DVLS server, with
// CertificateIdentifier as the file name.
func (c *EntryCertificateService) NewFile(entry EntryCertificate, content []byte) (EntryCertificate, error) {
	return c.NewFileWithContext(context.Background(), entry, content)
}

// NewFileWithContext creates a new EntryCertificate based on entry and stores content on the DVLS server, with
// CertificateIdentifier as the file name.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) NewFileWithContext(ctx context.Context, entry EntryCertificate, content []byte) (EntryCertificate, error) {
	if content == nil {
		content = []byte{}
	}

	return c.newWithContext(ctx, entry, EntryCertificateDataModeFile, content)
}

func (c *EntryCertificateService) newWithContext(ctx context.Context, entry EntryCertificate, mode EntryCertificateDataMode, content []byte) (EntryCertificate, error) {
	if entry.VaultId == "" {
		return EntryCertificate{}, fmt.Errorf("entry must have a VaultId")
	}

	reqUrl, err := entryUrl(c.client, nil, entry.VaultId)
	if err != nil {
		return EntryCertificate{}, err
	}

	reqBody := entryCreateRequest{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.EntryFolderPath,
		Type:        entryDocumentType,
		SubType:     entryCertificateSubType,
		Tags:        entry.Tags,
		Data:        entry.data(mode, content),
	}

	return c.send(ctx, http.MethodPost, reqUrl, reqBody, entry.Password, "error while creating entry")
}

// Update updates an EntryCertificate based on entry and keeps its data mode and stored file. Every other field is
// replaced; an empty Password keeps the stored one.
func (c *EntryCertificateService) Update(entry EntryCertificate) (EntryCertificate, error) {
	return c.UpdateWithContext(context.Background(), entry)
}

// UpdateWithContext updates an EntryCertificate based on entry and keeps its data mode and stored file. Every other
// field is replaced; an empty Password keeps the stored one.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) UpdateWithContext(ctx context.Context, entry EntryCertificate) (EntryCertificate, error) {
	oldEntry, err := c.GetWithContext(ctx, entry.VaultId, entry.Id)
	if err != nil {
		return EntryCertificate{}, err
	}

	reqUrl, err := entryUrl(c.client, nil, entry.VaultId, entry.Id)
	if err != nil {
		return EntryCertificate{}, err
	}

	tags := entry.Tags
	if tags == nil {
		tags = []string{}
	}

	reqBody := entryUpdateRequest{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.EntryFolderPath,
		Tags:        tags,
		Data:        entry.data(oldEntry.mode, nil),
	}

	return c.send(ctx, http.MethodPut, reqUrl, reqBody, entry.Password, "error while updating entry")
}

// send writes a certificate and returns the server's answer, which has no password: the password sent is put back.
func (c *EntryCertificateService) send(ctx context.Context, method string, reqUrl string, reqBody any, password string, errorPrefix string) (EntryCertificate, error) {
	var raw rawEntry
	err := c.client.sendJSON(ctx, method, reqUrl, reqBody, &raw)
	if err != nil {
		return EntryCertificate{}, fmt.Errorf("%s: %w", errorPrefix, err)
	}

	entry, err := certificateFromRaw(raw)
	if err != nil {
		return EntryCertificate{}, err
	}
	entry.Password = password

	return entry, nil
}

// Delete deletes an EntryCertificate based on vaultId and entryId.
func (c *EntryCertificateService) Delete(vaultId string, entryId string) error {
	return c.DeleteWithContext(context.Background(), vaultId, entryId)
}

// DeleteWithContext deletes an EntryCertificate based on vaultId and entryId.
// The provided context can be used to cancel the request.
func (c *EntryCertificateService) DeleteWithContext(ctx context.Context, vaultId string, entryId string) error {
	return c.client.deleteEntry(ctx, vaultId, entryId)
}

// GetDataMode returns the data mode of the EntryCertificate. Can be either EntryCertificateDataModeURL or EntryCertificateDataModeFile.
func (c EntryCertificate) GetDataMode() EntryCertificateDataMode {
	return c.mode
}
