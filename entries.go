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

const (
	entryEndpoint            string = "/api/connections/partial"
	entryConnectionsEndpoint string = "/api/connections"
	entryBasePublicEndpoint  string = "/api/v1/vault/{vaultId}/entry"
	entryPublicEndpoint      string = "/api/v1/vault/{vaultId}/entry/{id}"
)

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

func entryPublicEndpointReplacer(vaultId string, entryId string) string {
	replacer := strings.NewReplacer("{vaultId}", vaultId, "{id}", entryId)
	return replacer.Replace(entryPublicEndpoint)
}

func entryPublicBaseEndpointReplacer(vaultId string) string {
	replacer := strings.NewReplacer("{vaultId}", vaultId)
	return replacer.Replace(entryBasePublicEndpoint)
}

// entryListRawResponse represents the raw paginated response from the entry list endpoint.
type entryListRawResponse struct {
	Data        []json.RawMessage `json:"data"`
	CurrentPage int               `json:"currentPage"`
	PageSize    int               `json:"pageSize"`
	TotalCount  int               `json:"totalCount"`
	TotalPage   int               `json:"totalPage"`
}

// getEntriesOptions contains optional filters for listing entries.
// A nil value means the filter is not applied.
type GetEntriesOptions struct {
	Name *string
	Path *string
}

// getEntries returns a list of entries from a vault with optional filters.
// Entries with unsupported types are skipped.
// This function handles pagination automatically and returns all entries across all pages.
func (c *Client) getEntries(ctx context.Context, vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	if vaultId == "" {
		return nil, fmt.Errorf("vaultId is required")
	}

	baseEndpoint := entryPublicBaseEndpointReplacer(vaultId)
	reqUrl, err := url.JoinPath(c.baseUri, baseEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build entry url: %w", err)
	}

	parsedUrl, err := url.Parse(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry url: %w", err)
	}

	var allEntries []Entry
	currentPage := 1

	for {
		q := parsedUrl.Query()
		if opts.Name != nil {
			q.Set("name", *opts.Name)
		}
		if opts.Path != nil && *opts.Path != "" {
			q.Set("path", *opts.Path)
		}
		q.Set("page", fmt.Sprintf("%d", currentPage))
		parsedUrl.RawQuery = q.Encode()

		resp, err := c.RequestWithContext(ctx, parsedUrl.String(), http.MethodGet, nil)
		if err != nil {
			return nil, fmt.Errorf("error while fetching entries (page %d): %w", currentPage, err)
		}

		var rawResp entryListRawResponse
		if err := json.Unmarshal(resp.Response, &rawResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry list response (page %d): %w", currentPage, err)
		}

		for _, raw := range rawResp.Data {
			var entry Entry
			if err := json.Unmarshal(raw, &entry); err != nil {
				if IsUnsupportedEntryType(err) {
					continue
				}
				return nil, fmt.Errorf("failed to unmarshal entry (page %d): %w", currentPage, err)
			}
			entry.VaultId = vaultId
			allEntries = append(allEntries, entry)
		}

		// Check if we've fetched all pages
		if currentPage >= rawResp.TotalPage {
			break
		}
		currentPage++
	}

	// The server path filter is not exact, so we always apply client-side filtering when path
	// is set. We match entries at the exact path or any sub-path (prefix + backslash separator).
	// When path is "", the server ignores the filter, so we also handle root-level filtering here.
	if opts.Path != nil {
		var filtered []Entry
		for _, entry := range allEntries {
			if entry.Path == *opts.Path || (*opts.Path != "" && strings.HasPrefix(entry.Path, *opts.Path+"\\")) {
				filtered = append(filtered, entry)
			}
		}
		return filtered, nil
	}

	return allEntries, nil
}
