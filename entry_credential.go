package dvls

import (
	"context"
	"errors"
	"fmt"
)

var ErrEntryNotFound = errors.New("entry not found")
var ErrMultipleEntriesFound = errors.New("multiple entries found")

const (
	EntryCredentialType string = "Credential"

	EntryCredentialSubTypeAccessCode            string = "AccessCode"
	EntryCredentialSubTypeApiKey                string = "ApiKey"
	EntryCredentialSubTypeAzureServicePrincipal string = "AzureServicePrincipal"
	EntryCredentialSubTypeConnectionString      string = "ConnectionString"
	EntryCredentialSubTypeDefault               string = "Default"
	EntryCredentialSubTypePrivateKey            string = "PrivateKey"
)

// supportedCredentialSubTypes is generated from entryFactories to ensure a single source of truth.
var supportedCredentialSubTypes = getSupportedSubTypes(EntryCredentialType)

type EntryCredentialService service

// The data types send their secret fields (passwords, keys, connection strings) only when set, so updating an entry
// read from a list, which has no secrets, keeps them. Their other fields are always sent, so an empty one is cleared.

type EntryCredentialAccessCodeData struct {
	Password string `json:"password,omitempty"`
}

type EntryCredentialApiKeyData struct {
	ApiId    string `json:"apiId"`
	ApiKey   string `json:"apiKey,omitempty"`
	TenantId string `json:"tenantId"`
}

type EntryCredentialAzureServicePrincipalData struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"`
	TenantId     string `json:"tenantId"`
}

type EntryCredentialConnectionStringData struct {
	ConnectionString string `json:"connectionString,omitempty"`
}

type EntryCredentialDefaultData struct {
	Domain   string `json:"domain"`
	Password string `json:"password,omitempty"`
	Username string `json:"username"`
}

type EntryCredentialPrivateKeyData struct {
	Username   string `json:"privateKeyOverrideUsername"`
	Password   string `json:"privateKeyOverridePassword,omitempty"`
	PrivateKey string `json:"privateKeyData,omitempty"`
	PublicKey  string `json:"publicKeyData"`
	Passphrase string `json:"privateKeyPassPhrase,omitempty"`
}

func (e *Entry) GetCredentialAccessCodeData() (*EntryCredentialAccessCodeData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialAccessCodeData)
	return data, ok
}

func (e *Entry) GetCredentialApiKeyData() (*EntryCredentialApiKeyData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialApiKeyData)
	return data, ok
}

func (e *Entry) GetCredentialAzureServicePrincipalData() (*EntryCredentialAzureServicePrincipalData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialAzureServicePrincipalData)
	return data, ok
}

func (e *Entry) GetCredentialConnectionStringData() (*EntryCredentialConnectionStringData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialConnectionStringData)
	return data, ok
}

func (e *Entry) GetCredentialDefaultData() (*EntryCredentialDefaultData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialDefaultData)
	return data, ok
}

func (e *Entry) GetCredentialPrivateKeyData() (*EntryCredentialPrivateKeyData, bool) {
	if e == nil {
		return nil, false
	}

	data, ok := e.Data.(*EntryCredentialPrivateKeyData)
	return data, ok
}

// ToCredentialMap flattens a credential entry into a map of fields keyed by a stable name.
// It always includes "entry-id" and "entry-name" and then subtype-specific keys.
func (e *Entry) ToCredentialMap() (map[string]string, error) {
	if e.GetType() != EntryCredentialType {
		return nil, fmt.Errorf("unsupported entry type (%s). Only %s is supported", e.GetType(), EntryCredentialType)
	}

	secretMap := map[string]string{
		"entry-id":   e.Id,
		"entry-name": e.Name,
	}

	switch e.SubType {
	case EntryCredentialSubTypeDefault:
		if data, ok := e.GetCredentialDefaultData(); ok {
			if data.Username != "" {
				secretMap["username"] = data.Username
			}
			if data.Password != "" {
				secretMap["password"] = data.Password
			}
			if data.Domain != "" {
				secretMap["domain"] = data.Domain
			}
		}

	case EntryCredentialSubTypeAccessCode:
		if data, ok := e.GetCredentialAccessCodeData(); ok {
			if data.Password != "" {
				secretMap["password"] = data.Password
			}
		}

	case EntryCredentialSubTypeApiKey:
		if data, ok := e.GetCredentialApiKeyData(); ok {
			if data.ApiId != "" {
				secretMap["api-id"] = data.ApiId
			}
			if data.ApiKey != "" {
				secretMap["api-key"] = data.ApiKey
			}
			if data.TenantId != "" {
				secretMap["tenant-id"] = data.TenantId
			}
		}

	case EntryCredentialSubTypeAzureServicePrincipal:
		if data, ok := e.GetCredentialAzureServicePrincipalData(); ok {
			if data.ClientId != "" {
				secretMap["client-id"] = data.ClientId
			}
			if data.ClientSecret != "" {
				secretMap["client-secret"] = data.ClientSecret
			}
			if data.TenantId != "" {
				secretMap["tenant-id"] = data.TenantId
			}
		}

	case EntryCredentialSubTypeConnectionString:
		if data, ok := e.GetCredentialConnectionStringData(); ok {
			if data.ConnectionString != "" {
				secretMap["connection-string"] = data.ConnectionString
			}
		}

	case EntryCredentialSubTypePrivateKey:
		if data, ok := e.GetCredentialPrivateKeyData(); ok {
			if data.Username != "" {
				secretMap["username"] = data.Username
			}
			if data.Password != "" {
				secretMap["password"] = data.Password
			}
			if data.PrivateKey != "" {
				secretMap["private-key"] = data.PrivateKey
			}
			if data.PublicKey != "" {
				secretMap["public-key"] = data.PublicKey
			}
			if data.Passphrase != "" {
				secretMap["passphrase"] = data.Passphrase
			}
		}

	default:
		return nil, fmt.Errorf("unsupported credential subtype (%s)", e.SubType)
	}

	return secretMap, nil
}

// SetCredentialSecret mutates the Entry data to update the secret value for supported subtypes.
// It preserves existing fields and only updates the password/secret field.
func (e *Entry) SetCredentialSecret(secret string) error {
	if e.GetType() != EntryCredentialType {
		return fmt.Errorf("unsupported entry type (%s). Only %s is supported", e.GetType(), EntryCredentialType)
	}

	switch e.SubType {
	case EntryCredentialSubTypeDefault:
		if data, ok := e.GetCredentialDefaultData(); ok {
			data.Password = secret
		} else {
			e.Data = &EntryCredentialDefaultData{Password: secret}
		}
	case EntryCredentialSubTypeAccessCode:
		if data, ok := e.GetCredentialAccessCodeData(); ok {
			data.Password = secret
		} else {
			e.Data = &EntryCredentialAccessCodeData{Password: secret}
		}
	default:
		return fmt.Errorf("cannot set secret for credential subtype (%s)", e.SubType)
	}

	return nil
}

// validateEntry checks if an Entry has the required fields and valid type/subtype.
func (c *EntryCredentialService) validateEntry(entry *Entry) error {
	return validateEntryType(entry, EntryCredentialType, supportedCredentialSubTypes)
}

// Get returns a single EntryCredential based on the entry's VaultId and Id.
func (c *EntryCredentialService) Get(entry Entry) (Entry, error) {
	return c.GetWithContext(context.Background(), entry)
}

// GetWithContext returns a single EntryCredential based on the entry's VaultId and Id.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) GetWithContext(ctx context.Context, entry Entry) (Entry, error) {
	return c.GetByIdWithContext(ctx, entry.VaultId, entry.Id)
}

// GetById returns a single EntryCredential, its secrets included, based on vault Id and entry Id.
func (c *EntryCredentialService) GetById(vaultId string, entryId string) (Entry, error) {
	return c.GetByIdWithContext(context.Background(), vaultId, entryId)
}

// GetByIdWithContext returns a single EntryCredential, its secrets included, based on vault Id and entry Id.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) GetByIdWithContext(ctx context.Context, vaultId string, entryId string) (Entry, error) {
	return c.client.getEntry(ctx, vaultId, entryId)
}

// New creates a new EntryCredential and returns the new entry's Id.
func (c *EntryCredentialService) New(entry Entry) (string, error) {
	return c.NewWithContext(context.Background(), entry)
}

// NewWithContext creates a new EntryCredential and returns the new entry's Id.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) NewWithContext(ctx context.Context, entry Entry) (string, error) {
	if err := c.validateEntry(&entry); err != nil {
		return "", err
	}

	return c.client.newEntry(ctx, entry)
}

// Update updates an EntryCredential and returns the updated entry.
// A secret field left empty keeps its stored value; every other field is replaced, so an empty one is cleared.
func (c *EntryCredentialService) Update(entry Entry) (Entry, error) {
	return c.UpdateWithContext(context.Background(), entry)
}

// UpdateWithContext updates an EntryCredential and returns the updated entry.
// A secret field left empty keeps its stored value; every other field is replaced, so an empty one is cleared.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) UpdateWithContext(ctx context.Context, entry Entry) (Entry, error) {
	if err := c.validateEntry(&entry); err != nil {
		return Entry{}, err
	}

	return c.client.updateEntry(ctx, entry)
}

// Delete deletes an entry based on the entry's VaultId and Id.
func (c *EntryCredentialService) Delete(e Entry) error {
	return c.DeleteWithContext(context.Background(), e)
}

// DeleteWithContext deletes an entry based on the entry's VaultId and Id.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) DeleteWithContext(ctx context.Context, e Entry) error {
	return c.DeleteByIdWithContext(ctx, e.VaultId, e.Id)
}

// DeleteById deletes an entry based on vault Id and entry Id.
func (c *EntryCredentialService) DeleteById(vaultId string, entryId string) error {
	return c.DeleteByIdWithContext(context.Background(), vaultId, entryId)
}

// DeleteByIdWithContext deletes an entry based on vault Id and entry Id.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) DeleteByIdWithContext(ctx context.Context, vaultId string, entryId string) error {
	return c.client.deleteEntry(ctx, vaultId, entryId)
}

// GetEntries returns the credential entries of a vault, without their secrets, with optional filters.
func (c *EntryCredentialService) GetEntries(vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	return c.GetEntriesWithContext(context.Background(), vaultId, opts)
}

// GetEntriesWithContext returns the credential entries of a vault, without their secrets, with optional filters.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) GetEntriesWithContext(ctx context.Context, vaultId string, opts GetEntriesOptions) ([]Entry, error) {
	return c.client.getEntries(ctx, vaultId, EntryCredentialType, "", opts)
}

// GetByName retrieves a single credential entry, its secrets included, by name (ignoring case), subType, and
// optional filters.
// Returns ErrEntryNotFound if no match exists.
// Returns ErrMultipleEntriesFound if more than one match exists.
func (c *EntryCredentialService) GetByName(vaultId, name, subType string, opts GetByNameOptions) (Entry, error) {
	return c.GetByNameWithContext(context.Background(), vaultId, name, subType, opts)
}

// GetByNameWithContext retrieves a single credential entry, its secrets included, by name (ignoring case), subType,
// and optional filters.
// Returns ErrEntryNotFound if no match exists.
// Returns ErrMultipleEntriesFound if more than one match exists.
// The provided context can be used to cancel the request.
func (c *EntryCredentialService) GetByNameWithContext(ctx context.Context, vaultId, name, subType string, opts GetByNameOptions) (Entry, error) {
	matches, err := c.client.getEntries(ctx, vaultId, EntryCredentialType, subType, GetEntriesOptions{Name: &name, Path: opts.Path})
	if err != nil {
		return Entry{}, err
	}

	switch len(matches) {
	case 0:
		return Entry{}, ErrEntryNotFound
	case 1:
		return c.GetByIdWithContext(ctx, vaultId, matches[0].Id)
	default:
		return Entry{}, ErrMultipleEntriesFound
	}
}
