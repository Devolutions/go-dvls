package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type EntryUserCredentialService service

const (
	entryPublicEndpoint string = "/api/v1/vault/{vaultId}/entry/{id}"
	EntryTypeCredential string = "Credential"
	EntrySubTypeDefault string = "Default"
)

// EntryUserCredential represents a DVLS entry/connection.
type EntryUserCredential struct {
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

// EntryCredentials represents an EntryUserCredential Credentials fields.
type EntryCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func entryReplacer(vaultId string, entryId string) string {
	replacer := strings.NewReplacer("{vaultId}", vaultId, "{id}", entryId)
	return replacer.Replace(entryPublicEndpoint)
}

// validateEntry checks if an EntryUserCredential has the required fields and valid type/subtype.
func (c *EntryUserCredentialService) validateEntry(entry *EntryUserCredential) error {
	if entry.VaultId == "" {
		return fmt.Errorf("entry must have a VaultId")
	}

	if entry.Type != EntryTypeCredential {
		return fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.Type, EntryTypeCredential)
	}

	if entry.SubType == "" {
		entry.SubType = EntrySubTypeDefault
	} else if entry.SubType != EntrySubTypeDefault {
		return fmt.Errorf("unsupported entry subtype (%s). Only %s is supported", entry.SubType, EntrySubTypeDefault)
	}

	return nil
}

// Get returns a single EntryUserCredential specified by entryId.
func (c *EntryUserCredentialService) Get(vaultId string, entryId string) (EntryUserCredential, error) {
	if entryId == "" || vaultId == "" {
		return EntryUserCredential{}, fmt.Errorf("both entry ID and vault ID are required for deletion")
	}
	var entry EntryUserCredential
	entryUri := entryReplacer(vaultId, entryId)

	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while fetching entry. error: %w", err)
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	entry.VaultId = vaultId
	if entry.SubType == "" {
		entry.SubType = EntrySubTypeDefault
	}

	return entry, nil
}

// New creates a new EntryUserCredential based on entry.
func (c *EntryUserCredentialService) New(entry EntryUserCredential) (EntryUserCredential, error) {
	if err := c.validateEntry(&entry); err != nil {
		return EntryUserCredential{}, err
	}

	entry.ID = ""

	baseEntryEndpoint := strings.Replace(entryPublicEndpoint, "/{id}", "", 1)
	entryUri := strings.Replace(baseEntryEndpoint, "{vaultId}", entry.VaultId, 1)
	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while creating entry. error: %w", err)
	}
	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}
	return entry, nil
}

// Update updates an EntryUserCredential based on entry.
func (c *EntryUserCredentialService) Update(entry EntryUserCredential) (EntryUserCredential, error) {
	if err := c.validateEntry(&entry); err != nil {
		return EntryUserCredential{}, err
	}

	if entry.ID == "" {
		return EntryUserCredential{}, fmt.Errorf("entry ID is required for updates")
	}

	originalEntry, err := c.Get(entry.VaultId, entry.ID)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to fetch original entry. error: %w", err)
	}

	if originalEntry.SubType != entry.SubType {
		return EntryUserCredential{}, fmt.Errorf("entry subType cannot be changed")
	}

	entryUri := entryReplacer(entry.VaultId, entry.ID)

	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	_, err = c.client.Request(reqUrl, http.MethodPut, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while updating entry. error: %w", err)
	}

	entry, err = c.Get(entry.VaultId, entry.ID)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("update succeeded but failed to fetch updated entry: %w", err)
	}

	return entry, nil
}

// Delete deletes an entry based on entry.
func (c *EntryUserCredentialService) Delete(entry EntryUserCredential) error {
	if entry.ID == "" || entry.VaultId == "" {
		return fmt.Errorf("both entry ID and vault ID are required")
	}

	entryUri := entryReplacer(entry.VaultId, entry.ID)
	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return fmt.Errorf("failed to build delete entry url. error: %w", err)
	}

	_, err = c.client.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry. error: %w", err)
	}

	return nil
}
