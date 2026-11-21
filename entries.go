package dvls

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const entryEndpoint string = "entry"

// ErrUnsupportedEntryType is returned when an entry type/subtype is not supported by this client.
type ErrUnsupportedEntryType struct {
	Type    string
	SubType string
}

func (e ErrUnsupportedEntryType) Error() string {
	return fmt.Sprintf("unsupported entry type/subtype: %s/%s", e.Type, e.SubType)
}

// IsUnsupportedEntryType returns true if the error is an ErrUnsupportedEntryType.
func IsUnsupportedEntryType(err error) bool {
	var unsupportedErr ErrUnsupportedEntryType
	return errors.As(err, &unsupportedErr)
}

type Entries struct {
	Certificate *EntryCertificateService
	Host        *EntryHostService
	Credential  *EntryCredentialService
	Website     *EntryWebsiteService
	Folder      *EntryFolderService
	Permissions *EntryPermissionsService
}

type Entry struct {
	Id          string   `json:"id,omitempty"`
	VaultId     string   `json:"vaultId,omitempty"`
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Type        string   `json:"type"`
	SubType     string   `json:"subType"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`

	Data EntryData `json:"data,omitempty"`

	CreatedBy  string      `json:"createdBy,omitempty"`
	CreatedOn  *ServerTime `json:"createdOn,omitempty"`
	ModifiedBy string      `json:"modifiedBy,omitempty"`
	ModifiedOn *ServerTime `json:"modifiedOn,omitempty"`
}

type EntryData any

func (e *Entry) GetType() string {
	return e.Type
}

func (e *Entry) GetSubType() string {
	return e.SubType
}

var entryFactories = map[string]func() EntryData{
	"Credential/AccessCode":            func() EntryData { return &EntryCredentialAccessCodeData{} },
	"Credential/ApiKey":                func() EntryData { return &EntryCredentialApiKeyData{} },
	"Credential/AzureServicePrincipal": func() EntryData { return &EntryCredentialAzureServicePrincipalData{} },
	"Credential/ConnectionString":      func() EntryData { return &EntryCredentialConnectionStringData{} },
	"Credential/Default":               func() EntryData { return &EntryCredentialDefaultData{} },
	"Credential/PrivateKey":            func() EntryData { return &EntryCredentialPrivateKeyData{} },
	"Folder/Company":                   func() EntryData { return &EntryFolderData{} },
	"Folder/Credentials":               func() EntryData { return &EntryFolderData{} },
	"Folder/Customer":                  func() EntryData { return &EntryFolderData{} },
	"Folder/Database":                  func() EntryData { return &EntryFolderData{} },
	"Folder/Device":                    func() EntryData { return &EntryFolderData{} },
	"Folder/Domain":                    func() EntryData { return &EntryFolderData{} },
	"Folder/Folder":                    func() EntryData { return &EntryFolderData{} },
	"Folder/Identity":                  func() EntryData { return &EntryFolderData{} },
	"Folder/MacroScriptTools":          func() EntryData { return &EntryFolderData{} },
	"Folder/Printer":                   func() EntryData { return &EntryFolderData{} },
	"Folder/Server":                    func() EntryData { return &EntryFolderData{} },
	"Folder/Site":                      func() EntryData { return &EntryFolderData{} },
	"Folder/SmartFolder":               func() EntryData { return &EntryFolderData{} },
	"Folder/Software":                  func() EntryData { return &EntryFolderData{} },
	"Folder/Team":                      func() EntryData { return &EntryFolderData{} },
	"Folder/Workstation":               func() EntryData { return &EntryFolderData{} },
}

// getSupportedSubTypes extracts all supported subtypes for a given entry type from entryFactories.
// This ensures a single source of truth for supported entry types/subtypes.
func getSupportedSubTypes(entryType string) map[string]struct{} {
	result := make(map[string]struct{})
	prefix := entryType + "/"
	for key := range entryFactories {
		if subType, found := strings.CutPrefix(key, prefix); found {
			result[subType] = struct{}{}
		}
	}
	return result
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	type alias Entry
	raw := &struct {
		Data json.RawMessage `json:"data"`
		*alias
	}{
		alias: (*alias)(e),
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", raw.Type, raw.SubType)
	factory, ok := entryFactories[key]
	if !ok {
		return ErrUnsupportedEntryType{Type: raw.Type, SubType: raw.SubType}
	}

	dataStruct := factory()
	if err := json.Unmarshal(raw.Data, dataStruct); err != nil {
		return fmt.Errorf("failed to unmarshal entry data: %w", err)
	}

	e.Data = dataStruct

	return nil
}

func (e Entry) MarshalJSON() ([]byte, error) {
	type alias Entry

	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Data json.RawMessage `json:"data"`
		*alias
	}{
		Data:  dataBytes,
		alias: (*alias)(&e),
	})
}

// GetByNameOptions contains optional filters for GetByName.
// A nil field means the filter is not applied.
type GetByNameOptions struct {
	Path *string
}

// GetEntriesOptions contains optional filters for listing entries.
// A nil value means the filter is not applied. Name matches ignoring case. Path returns the entries of that folder
// and of its sub-folders (a folder is listed under its parent, not under itself); an empty Path returns the entries
// at the root of the vault.
type GetEntriesOptions struct {
	Name *string
	Path *string
}

func entryUrl(c *Client, query url.Values, vaultId string, entryId ...string) (string, error) {
	reqUrl, err := c.v2Url(query, append([]string{vaultEndpoint, vaultId, entryEndpoint}, entryId...)...)
	if err != nil {
		return "", fmt.Errorf("failed to build entry url: %w", err)
	}

	return reqUrl, nil
}

// getEntries returns the entries of a vault of the given type, and of the given subtype unless subType is empty,
// with optional filters.
// Entries with unsupported types are skipped.
// This function handles pagination automatically and returns all entries across all pages.
func (c *Client) getEntries(ctx context.Context, vaultId string, entryType string, subType string, opts GetEntriesOptions) ([]Entry, error) {
	if vaultId == "" {
		return nil, fmt.Errorf("vaultId is required")
	}

	query := url.Values{"type": {entryType}}
	if subType != "" {
		query.Set("subType", subType)
	}
	if opts.Name != nil {
		query.Set("name", *opts.Name)
	}
	if opts.Path != nil {
		if *opts.Path == "" {
			query.Set("exactPath", "true")
		} else {
			query.Set("path", *opts.Path)
		}
	}

	raws, err := listAll[json.RawMessage](ctx, c, query, vaultEndpoint, vaultId, entryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching entries: %w", err)
	}

	var entries []Entry
	for _, raw := range raws {
		var entry Entry
		if err := json.Unmarshal(raw, &entry); err != nil {
			if IsUnsupportedEntryType(err) {
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
		}

		entry.VaultId = vaultId
		entries = append(entries, entry)
	}

	return entries, nil
}

func (c *Client) getEntry(ctx context.Context, vaultId string, entryId string) (Entry, error) {
	if vaultId == "" || entryId == "" {
		return Entry{}, fmt.Errorf("both entry Id and vault Id are required")
	}

	reqUrl, err := entryUrl(c, nil, vaultId, entryId)
	if err != nil {
		return Entry{}, err
	}

	var entry Entry
	err = c.getJSON(ctx, reqUrl, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("error while fetching entry: %w", err)
	}

	entry.VaultId = vaultId

	return entry, nil
}

type entryCreateRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Path        string    `json:"path,omitempty"`
	Type        string    `json:"type"`
	SubType     string    `json:"subType,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Data        EntryData `json:"data"`
}

// entryUpdateRequest is the body of a PUT, where every field is required: tags and data are never null.
type entryUpdateRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Tags        []string  `json:"tags"`
	Data        EntryData `json:"data"`
}

func (c *Client) newEntry(ctx context.Context, entry Entry) (string, error) {
	reqUrl, err := entryUrl(c, nil, entry.VaultId)
	if err != nil {
		return "", err
	}

	reqBody := entryCreateRequest{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.Path,
		Type:        entry.Type,
		SubType:     entry.SubType,
		Tags:        entry.Tags,
		Data:        nonNilData(entry.Data),
	}

	var created struct {
		Id string `json:"id"`
	}
	err = c.sendJSON(ctx, http.MethodPost, reqUrl, reqBody, &created)
	if err != nil {
		return "", fmt.Errorf("error while creating entry: %w", err)
	}

	return created.Id, nil
}

func (c *Client) updateEntry(ctx context.Context, entry Entry) (Entry, error) {
	if entry.Id == "" {
		return Entry{}, fmt.Errorf("entry Id is required for updates")
	}

	reqUrl, err := entryUrl(c, nil, entry.VaultId, entry.Id)
	if err != nil {
		return Entry{}, err
	}

	tags := entry.Tags
	if tags == nil {
		tags = []string{}
	}

	reqBody := entryUpdateRequest{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.Path,
		Tags:        tags,
		Data:        nonNilData(entry.Data),
	}

	err = c.sendJSON(ctx, http.MethodPut, reqUrl, reqBody, nil)
	if err != nil {
		return Entry{}, fmt.Errorf("error while updating entry: %w", err)
	}

	// The PUT answer has no secrets: read the entry back to return them.
	updated, err := c.getEntry(ctx, entry.VaultId, entry.Id)
	if err != nil {
		return Entry{}, fmt.Errorf("update succeeded but failed to fetch updated entry: %w", err)
	}

	return updated, nil
}

func (c *Client) deleteEntry(ctx context.Context, vaultId string, entryId string) error {
	if vaultId == "" || entryId == "" {
		return fmt.Errorf("both entry Id and vault Id are required")
	}

	reqUrl, err := entryUrl(c, nil, vaultId, entryId)
	if err != nil {
		return err
	}

	err = c.sendJSON(ctx, http.MethodDelete, reqUrl, nil, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry: %w", err)
	}

	return nil
}

// rawEntry is an entry as the API sends it, its data left undecoded, for the types Entry does not decode.
type rawEntry struct {
	Id          string          `json:"id,omitempty"`
	VaultId     string          `json:"vaultId,omitempty"`
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	Type        string          `json:"type"`
	SubType     string          `json:"subType,omitempty"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Data        json.RawMessage `json:"data"`
	ModifiedOn  *ServerTime     `json:"modifiedOn,omitempty"`
}

func (c *Client) getRawEntry(ctx context.Context, vaultId string, entryId string, query url.Values) (rawEntry, error) {
	if vaultId == "" || entryId == "" {
		return rawEntry{}, fmt.Errorf("both entry Id and vault Id are required")
	}

	reqUrl, err := entryUrl(c, query, vaultId, entryId)
	if err != nil {
		return rawEntry{}, err
	}

	var entry rawEntry
	err = c.getJSON(ctx, reqUrl, &entry)
	if err != nil {
		return rawEntry{}, fmt.Errorf("error while fetching entry: %w", err)
	}

	return entry, nil
}

func nonNilData(data EntryData) EntryData {
	if data == nil {
		return struct{}{}
	}

	return data
}

// validateEntryType checks that an Entry has a VaultId and a supported type/subtype.
func validateEntryType(entry *Entry, entryType string, supportedSubTypes map[string]struct{}) error {
	if entry.VaultId == "" {
		return fmt.Errorf("entry must have a VaultId")
	}

	if entry.GetType() != entryType {
		return fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.GetType(), entryType)
	}

	subType := entry.GetSubType()
	if _, isSupported := supportedSubTypes[subType]; !isSupported {
		var supportedList []string
		for st := range supportedSubTypes {
			supportedList = append(supportedList, st)
		}
		return fmt.Errorf("unsupported entry subtype (%s). Supported subtypes: %v", subType, supportedList)
	}

	return nil
}
