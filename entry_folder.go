package dvls

import (
	"context"
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
	Domain   string `json:"domain"`
	Username string `json:"username"`
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
	return validateEntryType(entry, EntryFolderType, supportedFolderSubTypes)
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

// GetById returns a single EntryFolder based on vault Id and entry Id. The Path of a folder is its parent folder's
// path, empty at the root of the vault.
func (c *EntryFolderService) GetById(vaultId string, entryId string) (Entry, error) {
	return c.GetByIdWithContext(context.Background(), vaultId, entryId)
}

// GetByIdWithContext returns a single EntryFolder based on vault Id and entry Id. The Path of a folder is its parent
// folder's path, empty at the root of the vault.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetByIdWithContext(ctx context.Context, vaultId string, entryId string) (Entry, error) {
	return c.client.getEntry(ctx, vaultId, entryId)
}

// New creates a new EntryFolder and returns the new entry's Id. Path is the parent folder's path, empty for the
// root of the vault.
func (c *EntryFolderService) New(entry Entry) (string, error) {
	return c.NewWithContext(context.Background(), entry)
}

// NewWithContext creates a new EntryFolder and returns the new entry's Id. Path is the parent folder's path, empty
// for the root of the vault.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) NewWithContext(ctx context.Context, entry Entry) (string, error) {
	if err := c.validateEntry(&entry); err != nil {
		return "", err
	}

	return c.client.newEntry(ctx, entry)
}

// Update updates an EntryFolder and returns the updated entry. Path is the parent folder's path: another path moves
// the folder with its content, and Name renames it.
func (c *EntryFolderService) Update(entry Entry) (Entry, error) {
	return c.UpdateWithContext(context.Background(), entry)
}

// UpdateWithContext updates an EntryFolder and returns the updated entry. Path is the parent folder's path: another
// path moves the folder with its content, and Name renames it.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) UpdateWithContext(ctx context.Context, entry Entry) (Entry, error) {
	if err := c.validateEntry(&entry); err != nil {
		return Entry{}, err
	}

	return c.client.updateEntry(ctx, entry)
}

// Delete deletes an empty folder based on the entry's VaultId and Id.
func (c *EntryFolderService) Delete(e Entry) error {
	return c.DeleteWithContext(context.Background(), e)
}

// DeleteWithContext deletes an empty folder based on the entry's VaultId and Id.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) DeleteWithContext(ctx context.Context, e Entry) error {
	return c.DeleteByIdWithContext(ctx, e.VaultId, e.Id)
}

// DeleteById deletes an empty folder based on vault Id and entry Id.
func (c *EntryFolderService) DeleteById(vaultId string, entryId string) error {
	return c.DeleteByIdWithContext(context.Background(), vaultId, entryId)
}

// DeleteByIdWithContext deletes an empty folder based on vault Id and entry Id. A folder that still holds entries
// is refused with an HTTP 409.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) DeleteByIdWithContext(ctx context.Context, vaultId string, entryId string) error {
	return c.client.deleteEntry(ctx, vaultId, entryId)
}

// GetByName retrieves a single folder entry by name (ignoring case) and optional filters.
// Returns ErrEntryNotFound if no match exists.
// Returns ErrMultipleEntriesFound if more than one match exists.
func (c *EntryFolderService) GetByName(vaultId, name string, opts GetByNameOptions) (Entry, error) {
	return c.GetByNameWithContext(context.Background(), vaultId, name, opts)
}

// GetByNameWithContext retrieves a single folder entry by name (ignoring case) and optional filters.
// Returns ErrEntryNotFound if no match exists.
// Returns ErrMultipleEntriesFound if more than one match exists.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetByNameWithContext(ctx context.Context, vaultId, name string, opts GetByNameOptions) (Entry, error) {
	entries, err := c.GetEntriesWithContext(ctx, vaultId, GetEntriesOptions{Name: &name, Path: opts.Path})
	if err != nil {
		return Entry{}, err
	}

	switch len(entries) {
	case 0:
		return Entry{}, ErrEntryNotFound
	case 1:
		return c.GetByIdWithContext(ctx, vaultId, entries[0].Id)
	default:
		return Entry{}, ErrMultipleEntriesFound
	}
}

// GetEntries returns the folder entries of a vault with optional filters.
func (c *EntryFolderService) GetEntries(vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	return c.GetEntriesWithContext(context.Background(), vaultId, opts)
}

// GetEntriesWithContext returns the folder entries of a vault with optional filters.
// The provided context can be used to cancel the request.
func (c *EntryFolderService) GetEntriesWithContext(ctx context.Context, vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	return c.client.getEntries(ctx, vaultId, EntryFolderType, "", opts)
}
