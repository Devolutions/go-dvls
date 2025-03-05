package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Entry represents a DVLS entry/connection.
type Entry struct {
	ID          string       `json:"id,omitempty"`
	VaultId     string       `json:"vaultId"`
	EntryName   string       `json:"name"`
	Description string       `json:"description"`
	Path        string       `json:"path"`
	ModifiedOn  *ServerTime  `json:"modifiedOn,omitempty"`
	ModifiedBy  string       `json:"modifiedBy,omitempty"`
	CreatedOn   *ServerTime  `json:"createdOn,omitempty"`
	CreatedBy   string       `json:"createdBy,omitempty"`
	Type        EntryType    `json:"type"`
	SubType     EntrySubType `json:"subType"`
	Tags        []string     `json:"tags,omitempty"`

	Credentials EntryCredentials `json:"data,omitempty"`
}

// EntryCredentials represents an Entry Credentials fields.
type EntryCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

const (
	entryEndpoint string = "/api/v1/vault/{vaultId}/entry/{id}"
)

func entryReplacer(vaultId string, entryId string) string {
	replacer := strings.NewReplacer("{vaultId}", vaultId, "{id}", entryId)
	return replacer.Replace(entryEndpoint)
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
	if entry.Type != EntryTypeCredential || entry.SubType != EntrySubTypeDefault {
		return Entry{}, fmt.Errorf("unsupported entry type (%s %s). Only %s %s is supported", entry.Type, entry.SubType, EntryTypeCredential, EntrySubTypeDefault)
	}

	entry.ID = ""
	entry.ModifiedOn = nil

	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, "save")
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
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Entry{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}

// UpdateEntry updates an Entry based on entry.
func (c *Client) UpdateEntry(entry Entry) (Entry, error) {
	if entry.Type != EntryTypeCredential || entry.SubType != EntrySubTypeDefault {
		return Entry{}, fmt.Errorf("unsupported entry type (%s %s). Only %s %s is supported", entry.Type, entry.SubType, EntryTypeCredential, EntrySubTypeDefault)
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

	resp, err := c.Request(reqUrl, http.MethodPatch, bytes.NewBuffer(entryJson))
	if err != nil {
		return Entry{}, fmt.Errorf("error while updating entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Entry{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}

// DeleteEntry deletes an entry based on entryId.
func (c *Client) DeleteEntry(entryId string) error {
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId)
	if err != nil {
		return fmt.Errorf("failed to delete entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}
