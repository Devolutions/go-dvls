package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	entryConnectionsEndpoint string = "/api/connections"
	entryEndpoint            string = "/api/connections/partial"
	entryPublicEndpoint      string = "/api/v1/vault/{vaultId}/entry/{id}"
	EntryTypeCredential      string = "Credential"
)

type Entries struct {
	Certificate *EntryCertificateService
	Host        *EntryHostService
	Website     *EntryWebsiteService
}

type Entry struct {
	ID          string      `json:"id,omitempty"`
	VaultId     string      `json:"vaultId"`
	EntryName   string      `json:"name"`
	Description string      `json:"description"`
	Path        string      `json:"path"`
	ModifiedOn  *ServerTime `json:"modifiedOn,omitempty"`
	ModifiedBy  string      `json:"modifiedBy,omitempty"`
	CreatedOn   *ServerTime `json:"createdOn,omitempty"`
	CreatedBy   string      `json:"createdBy,omitempty"`
	Type        string      `json:"type"`
	SubType     string      `json:"subType"`
	Tags        []string    `json:"tags,omitempty"`

	Credentials EntryCredentials `json:"data,omitempty"`
}

// EntryCredentials represents an Entry Credentials fields.
type EntryCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func entryReplacer(vaultId string, entryId string) string {
	replacer := strings.NewReplacer("{vaultId}", vaultId, "{id}", entryId)
	return replacer.Replace(entryPublicEndpoint)
}

// GetEntry returns a single Entry specified by entryId.
func (c *Client) GetEntry(vaultId string, entryId string) (Entry, error) {
	var entry Entry
	entryUri := entryReplacer(vaultId, entryId)

	reqUrl, err := url.JoinPath(c.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return Entry{}, fmt.Errorf("error while fetching entry. error: %w", err)
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	entry.VaultId = vaultId

	return entry, nil
}

// NewEntry creates a new Entry based on entry.
func (c *Client) NewEntry(entry Entry) (Entry, error) {
	if entry.Type != EntryTypeCredential {
		return Entry{}, fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.Type, EntryTypeCredential)
	}
	if entry.SubType == "" {
		entry.SubType = "Default"
	}

	entry.ID = ""
	entry.ModifiedOn = nil

	baseEntryEndpoint := strings.Replace(entryPublicEndpoint, "/{id}", "", 1)
	entryUri := strings.Replace(baseEntryEndpoint, "{vaultId}", entry.VaultId, 1)
	reqUrl, err := url.JoinPath(c.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}
	entryJson, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to marshall body. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return Entry{}, fmt.Errorf("error while creating entry. error: %w", err)
	}
	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}
	return entry, nil
}

// UpdateEntry updates an Entry based on entry.
func (c *Client) UpdateEntry(entry Entry) (Entry, error) {
	if entry.Type != EntryTypeCredential {
		return Entry{}, fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.Type, EntryTypeCredential)
	}

	// Always set SubType to Default
	entry.SubType = "Default"

	if entry.ID == "" {
		return Entry{}, fmt.Errorf("entry ID is required for updates")
	}

	entryUri := entryReplacer(entry.VaultId, entry.ID)

	reqUrl, err := url.JoinPath(c.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to marshall body. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPut, bytes.NewBuffer(entryJson))
	if err != nil {
		return Entry{}, fmt.Errorf("error while updating entry. error: %w", err)
	}

	// Handle empty response from server
	if len(resp.Response) == 0 {
		return entry, nil
	}

	// If we have a response body, try to unmarshal it
	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}

// DeleteEntry deletes an entry based on entry.
func (c *Client) DeleteEntry(entry Entry) error {
	if entry.ID == "" || entry.VaultId == "" {
		return fmt.Errorf("both entry ID and vault ID are required for deletion")
	}

	entryUri := entryReplacer(entry.VaultId, entry.ID)
	reqUrl, err := url.JoinPath(c.baseUri, entryUri)
	if err != nil {
		return fmt.Errorf("failed to build delete entry url. error: %w", err)
	}

	_, err = c.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry. error: %w", err)
	}

	return nil
}

func keywordsToSlice(kw string) []string {
	var spacedTag bool
	tags := strings.FieldsFunc(string(kw), func(r rune) bool {
		if r == '"' {
			spacedTag = !spacedTag
		}
		return !spacedTag && r == ' '
	})
	for i, v := range tags {
		unquotedTag, err := strconv.Unquote(v)
		if err != nil {
			continue
		}

		tags[i] = unquotedTag
	}

	return tags
}

func sliceToKeywords(kw []string) string {
	keywords := []string(kw)
	for i, v := range keywords {
		if strings.Contains(v, " ") {
			kw[i] = "\"" + v + "\""
		}
	}

	kString := strings.Join(keywords, " ")

	return kString
}
