package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	EntryFolderType string = "Folder"

	// 16 subtypes (all have the same behavior, difference is UI only)
	EntryFolderSubTypeCompany          string = "Company"
	EntryFolderSubTypeCredentials      string = "Credentials"
	EntryFolderSubTypeCustomer         string = "Customer"
	EntryFolderSubTypeDatabase         string = "Database"
	EntryFolderSubTypeDevice           string = "Device"
	EntryFolderSubTypeDomain           string = "Domain"
	EntryFolderSubTypeFolder           string = "Folder" // default
	EntryFolderSubTypeIdentity         string = "Identity"
	EntryFolderSubTypeMacroScriptTools string = "MacroScriptTools"
	EntryFolderSubTypePrinter          string = "Printer"
	EntryFolderSubTypeServer           string = "Server"
	EntryFolderSubTypeSite             string = "Site"
	EntryFolderSubTypeSmartFolder      string = "SmartFolder"
	EntryFolderSubTypeSoftware         string = "Software"
	EntryFolderSubTypeTeam             string = "Team"
	EntryFolderSubTypeWorkstation      string = "Workstation"
)

// supportedFolderSubTypes is generated from entryFactories to ensure a single source of truth.
var supportedFolderSubTypes = getSupportedSubTypes(EntryFolderType)

type EntryFolderService service

type EntryFolderData struct {
	Domain   string `json:"domain,omitempty"`
	Username string `json:"username,omitempty"`
}

func (e *Entry) GetFolderData() (*EntryFolderData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryFolderData)
	return data, ok
}

// validateEntry checks if an Entry has the required fields and valid type/subtype.
func (c *EntryFolderService) validateEntry(entry *Entry) error {
	if entry.VaultId == "" {
		return fmt.Errorf("entry must have a VaultId")
	}

	if entry.GetType() != EntryFolderType {
		return fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.GetType(), EntryFolderType)
	}

	subType := entry.GetSubType()
	if _, isSupported := supportedFolderSubTypes[subType]; !isSupported {
		var supportedList []string
		for st := range supportedFolderSubTypes {
			supportedList = append(supportedList, st)
		}
		return fmt.Errorf("unsupported entry subtype (%s). Supported subtypes: %v", subType, supportedList)
	}

	return nil
}

// Get returns a single EntryFolder based on the entry's VaultId and Id.
func (c *EntryFolderService) Get(entry Entry) (Entry, error) {
	return c.GetWithContext(context.Background(), entry)
}

// GetWithContext returns a single EntryFolder based on the entry's VaultId and Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetWithContext(ctx context.Context, entry Entry) (Entry, error) {
	return c.GetByIdWithContext(ctx, entry.VaultId, entry.Id)
}

// GetById returns a single EntryFolder based on vault Id and entry Id.
func (c *EntryFolderService) GetById(vaultId string, entryId string) (Entry, error) {
	return c.GetByIdWithContext(context.Background(), vaultId, entryId)
}

// GetByIdWithContext returns a single EntryFolder based on vault Id and entry Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetByIdWithContext(ctx context.Context, vaultId string, entryId string) (Entry, error) {
	if vaultId == "" || entryId == "" {
		return Entry{}, fmt.Errorf("both entry Id and vault Id are required")
	}

	var entry Entry
	entryUri := entryPublicEndpointReplacer(vaultId, entryId)

	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return Entry{}, fmt.Errorf("error while fetching entry: %w", err)
	}

	err = entry.UnmarshalJSON(resp.Response)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	entry.VaultId = vaultId

	return entry, nil
}

// New creates a new EntryFolder and returns the new entry's Id.
func (c *EntryFolderService) New(entry Entry) (string, error) {
	return c.NewWithContext(context.Background(), entry)
}

// NewWithContext creates a new EntryFolder and returns the new entry's Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) NewWithContext(ctx context.Context, entry Entry) (string, error) {
	if err := c.validateEntry(&entry); err != nil {
		return "", err
	}

	newEntryRequest := struct {
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		Path        string    `json:"path,omitempty"`
		Type        string    `json:"type"`
		SubType     string    `json:"subType"`
		Tags        []string  `json:"tags,omitempty"`
		Data        EntryData `json:"data"`
	}{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.Path,
		Type:        entry.GetType(),
		SubType:     entry.GetSubType(),
		Tags:        entry.Tags,
		Data:        entry.Data,
	}

	baseEntryEndpoint := entryPublicBaseEndpointReplacer(entry.VaultId)
	reqUrl, err := url.JoinPath(c.client.baseUri, baseEntryEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to build entry url: %w", err)
	}

	body, err := json.Marshal(newEntryRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error while creating entry: %w", err)
	}

	newEntryResponse := struct {
		Id string `json:"id"`
	}{}

	err = json.Unmarshal(resp.Response, &newEntryResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return newEntryResponse.Id, nil
}

// Update updates an EntryFolder and returns the updated entry.
func (c *EntryFolderService) Update(entry Entry) (Entry, error) {
	return c.UpdateWithContext(context.Background(), entry)
}

// UpdateWithContext updates an EntryFolder and returns the updated entry.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) UpdateWithContext(ctx context.Context, entry Entry) (Entry, error) {
	if err := c.validateEntry(&entry); err != nil {
		return Entry{}, err
	}

	if entry.Id == "" {
		return Entry{}, fmt.Errorf("entry Id is required for updates")
	}

	updateEntryRequest := struct {
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		Path        string    `json:"path,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
		Data        EntryData `json:"data"`
	}{
		Name:        entry.Name,
		Description: entry.Description,
		Path:        entry.Path,
		Tags:        entry.Tags,
		Data:        entry.Data,
	}

	entryUri := entryPublicEndpointReplacer(entry.VaultId, entry.Id)
	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	body, err := json.Marshal(updateEntryRequest)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	_, err = c.client.RequestWithContext(ctx, reqUrl, http.MethodPut, bytes.NewBuffer(body))
	if err != nil {
		return Entry{}, fmt.Errorf("error while updating entry: %w", err)
	}

	entry, err = c.GetByIdWithContext(ctx, entry.VaultId, entry.Id)
	if err != nil {
		return Entry{}, fmt.Errorf("update succeeded but failed to fetch updated entry: %w", err)
	}

	return entry, nil
}

// Delete deletes an entry based on the entry's VaultId and Id.
func (c *EntryFolderService) Delete(e Entry) error {
	return c.DeleteWithContext(context.Background(), e)
}

// DeleteWithContext deletes an entry based on the entry's VaultId and Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) DeleteWithContext(ctx context.Context, e Entry) error {
	return c.DeleteByIdWithContext(ctx, e.VaultId, e.Id)
}

// DeleteById deletes an entry based on vault Id and entry Id.
func (c *EntryFolderService) DeleteById(vaultId string, entryId string) error {
	return c.DeleteByIdWithContext(context.Background(), vaultId, entryId)
}

// DeleteByIdWithContext deletes an entry based on vault Id and entry Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) DeleteByIdWithContext(ctx context.Context, vaultId string, entryId string) error {
	if vaultId == "" || entryId == "" {
		return fmt.Errorf("both entry Id and vault Id are required")
	}

	entryUri := entryPublicEndpointReplacer(vaultId, entryId)
	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return fmt.Errorf("failed to build delete entry url: %w", err)
	}

	_, err = c.client.RequestWithContext(ctx, reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry: %w", err)
	}

	return nil
}

// GetByName retrieves a single folder entry by name and optional filters.
// Returns ErrEntryNotFound if no match exists.
func (c *EntryFolderService) GetByName(vaultId, name string, opts GetByNameOptions) (Entry, error) {
	return c.GetByNameWithContext(context.Background(), vaultId, name, opts)
}

// GetByNameWithContext retrieves a single folder entry by name and optional filters.
// Returns ErrEntryNotFound if no match exists.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetByNameWithContext(ctx context.Context, vaultId, name string, opts GetByNameOptions) (Entry, error) {
	entries, err := c.GetEntriesWithContext(ctx, vaultId, GetEntriesOptions{Name: &name, Path: opts.Path})
	if err != nil {
		return Entry{}, err
	}

	if len(entries) == 0 {
		return Entry{}, ErrEntryNotFound
	}

	return c.GetByIdWithContext(ctx, vaultId, entries[0].Id)
}

// GetEntries returns a list of folder entries from a vault with optional filters.
// Note: The API does not support filtering by entry type, so all entries are fetched and filtered client-side.
func (c *EntryFolderService) GetEntries(vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	return c.GetEntriesWithContext(context.Background(), vaultId, opts)
}

// GetEntriesWithContext returns a list of folder entries from a vault with optional filters.
// The provided context can be used to cancel the request.
// Note: The API does not support filtering by entry type, so all entries are fetched and filtered client-side.
func (c *EntryFolderService) GetEntriesWithContext(ctx context.Context, vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	entries, err := c.client.getEntries(ctx, vaultId, GetEntriesOptions{
		Name: opts.Name,
		Path: opts.Path,
	})
	if err != nil {
		return nil, err
	}

	// Filter only Folder type entries
	var folders []Entry
	for _, entry := range entries {
		if entry.GetType() == EntryFolderType {
			folders = append(folders, entry)
		}
	}

	return folders, nil
}
