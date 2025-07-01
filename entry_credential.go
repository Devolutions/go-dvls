package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	EntryCredentialType string = "Credential"

	EntryCredentialSubTypeAccessCode            string = "AccessCode"
	EntryCredentialSubTypeApiKey                string = "ApiKey"
	EntryCredentialSubTypeAzureServicePrincipal string = "AzureServicePrincipal"
	EntryCredentialSubTypeConnectionString      string = "ConnectionString"
	EntryCredentialSubTypeDefault               string = "Default"
	EntryCredentialSubTypePrivateKey            string = "PrivateKey"
)

type EntryCredentialService service

type EntryCredentialAccessCodeData struct {
	Password string `json:"password,omitempty"`
}

type EntryCredentialApiKeyData struct {
	ApiId    string `json:"apiId,omitempty"`
	ApiKey   string `json:"apiKey,omitempty"`
	TenantId string `json:"tenantId,omitempty"`
}

type EntryCredentialAzureServicePrincipalData struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	TenantId     string `json:"tenantId,omitempty"`
}

type EntryCredentialConnectionStringData struct {
	ConnectionString string `json:"connectionString,omitempty"`
}

type EntryCredentialDefaultData struct {
	Domain   string `json:"domain,omitempty"`
	Password string `json:"password,omitempty"`
	Username string `json:"username,omitempty"`
}

type EntryCredentialPrivateKeyData struct {
	Username   string `json:"privateKeyOverrideUsername,omitempty"`
	Password   string `json:"privateKeyOverridePassword,omitempty"`
	PrivateKey string `json:"privateKeyData,omitempty"`
	PublicKey  string `json:"publicKeyData,omitempty"`
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

// validateEntry checks if an Entry has the required fields and valid type/subtype.
func (c *EntryCredentialService) validateEntry(entry *Entry) error {
	if entry.VaultId == "" {
		return fmt.Errorf("entry must have a VaultId")
	}

	if entry.GetType() != EntryCredentialType {
		return fmt.Errorf("unsupported entry type (%s). Only %s is supported", entry.GetType(), EntryCredentialType)
	}

	supportedSubTypes := []string{
		EntryCredentialSubTypeAccessCode,
		EntryCredentialSubTypeApiKey,
		EntryCredentialSubTypeAzureServicePrincipal,
		EntryCredentialSubTypeConnectionString,
		EntryCredentialSubTypeDefault,
		EntryCredentialSubTypePrivateKey,
	}

	subType := entry.GetSubType()
	isSupported := false
	for _, t := range supportedSubTypes {
		if subType == t {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("unsupported entry subtype (%s). Supported subtypes: %v", subType, supportedSubTypes)
	}

	return nil
}

// Get returns a single EntryCredential
func (c *EntryCredentialService) Get(entry Entry) (Entry, error) {
	return c.GetById(entry.VaultId, entry.Id)
}

// Get returns a single EntryCredential based on vault Id and entry Id.
func (c *EntryCredentialService) GetById(vaultId string, entryId string) (Entry, error) {
	if vaultId == "" || entryId == "" {
		return Entry{}, fmt.Errorf("both entry Id and vault Id are required")
	}

	var entry Entry
	entryUri := entryPublicEndpointReplacer(vaultId, entryId)

	reqUrl, err := url.JoinPath(c.client.baseUri, entryUri)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return Entry{}, fmt.Errorf("error while fetching entry. error: %w", err)
	}

	err = entry.UnmarshalJSON(resp.Response)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	entry.VaultId = vaultId

	return entry, nil
}

// New creates a new EntryCredential
func (c *EntryCredentialService) New(entry Entry) (string, error) {
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
		return "", fmt.Errorf("failed to build entry url. error: %w", err)
	}

	body, err := json.Marshal(newEntryRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error while creating entry. error: %w", err)
	}

	newEntryResponse := struct {
		Id string `json:"id"`
	}{}

	err = json.Unmarshal(resp.Response, &newEntryResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}
	return newEntryResponse.Id, nil
}

// Update updates an EntryCredential
func (c *EntryCredentialService) Update(entry Entry) (Entry, error) {
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
		return Entry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	body, err := json.Marshal(updateEntryRequest)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	_, err = c.client.Request(reqUrl, http.MethodPut, bytes.NewBuffer(body))
	if err != nil {
		return Entry{}, fmt.Errorf("error while updating entry. error: %w", err)
	}

	entry, err = c.GetById(entry.VaultId, entry.Id)
	if err != nil {
		return Entry{}, fmt.Errorf("update succeeded but failed to fetch updated entry: %w", err)
	}

	return entry, nil
}

// Delete deletes an entry
func (c *EntryCredentialService) Delete(e Entry) error {
	return c.DeleteById(e.VaultId, e.Id)
}

// Delete deletes an entry based on vault Id and entry Id
func (c *EntryCredentialService) DeleteById(vaultId string, entryId string) error {
	if vaultId == "" || entryId == "" {
		return fmt.Errorf("both entry Id and vault Id are required")
	}

	entryUri := entryPublicEndpointReplacer(vaultId, entryId)
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
