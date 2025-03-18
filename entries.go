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

// EntryCredentialsData interface for all credential types
type EntryCredentialsData interface {
	GetSubType() string
}

// Entry represents an entry in the vault
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

	// The actual credentials data - only one will be populated based on SubType
	credentialsData EntryCredentialsData `json:"-"`
}

// MarshalJSON customizes how Entry is converted to JSON
func (e Entry) MarshalJSON() ([]byte, error) {
	type EntryAlias Entry // Avoid infinite recursion

	// Create a struct that will be marshaled to JSON
	aliasValue := struct {
		EntryAlias
		Data interface{} `json:"data,omitempty"`
	}{
		EntryAlias: EntryAlias(e),
		Data:       e.credentialsData,
	}

	return json.Marshal(aliasValue)
}

// UnmarshalJSON customizes how JSON is converted to Entry
func (e *Entry) UnmarshalJSON(data []byte) error {
	type EntryAlias Entry

	// Create a struct that will hold the unmarshaled data
	aliasValue := struct {
		EntryAlias
		Data json.RawMessage `json:"data,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	// Copy all the fields except credentialsData
	*e = Entry(aliasValue.EntryAlias)

	// If no data field, return early
	if len(aliasValue.Data) == 0 {
		return nil
	}

	// Unmarshal into the appropriate type based on SubType
	var err error
	switch e.SubType {
	case "":
	case "Default":
		var creds DefaultCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "PrivateKey":
		var creds PrivateKeyCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "AccessCode":
		var creds AccessCodeCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "ApiKey":
		var creds ApiKeyCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "AzureServicePrincipal":
		var creds AzureServicePrincipalCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "ConnectionString":
		var creds ConnectionStringCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	case "Passkey":
		var creds PasskeyCredentials
		err = json.Unmarshal(aliasValue.Data, &creds)
		e.credentialsData = &creds
	default:
		return fmt.Errorf("unknown credential subtype: %s", e.SubType)
	}

	return err
}

// Helper methods to get and set specific credential types
func (e *Entry) GetCredentials() EntryCredentialsData {
	return e.credentialsData
}

func (e *Entry) SetCredentials(creds EntryCredentialsData) {
	e.credentialsData = creds
	e.SubType = creds.GetSubType()
}

// Define all the credential types with their specific fields

// DefaultCredentials represents the Default credential subtype
type DefaultCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Domain   string `json:"domain,omitempty"`
}

func (c *DefaultCredentials) GetSubType() string {
	return "Default"
}

// PrivateKeyCredentials represents the PrivateKey credential subtype
type PrivateKeyCredentials struct {
	PrivateKeyData             string `json:"privateKeyData,omitempty"`
	PublicKeyData              string `json:"publicKeyData,omitempty"`
	PrivateKeyOverridePassword string `json:"privateKeyOverridePassword,omitempty"`
	PrivateKeyPassPhrase       string `json:"privateKeyPassPhrase,omitempty"`
}

func (c *PrivateKeyCredentials) GetSubType() string {
	return "PrivateKey"
}

// AccessCodeCredentials represents the AccessCode credential subtype
type AccessCodeCredentials struct {
	Password string `json:"password,omitempty"`
}

func (c *AccessCodeCredentials) GetSubType() string {
	return "AccessCode"
}

// ApiKeyCredentials represents the ApiKey credential subtype
type ApiKeyCredentials struct {
	ApiId    string `json:"apiId,omitempty"`
	ApiKey   string `json:"apiKey,omitempty"`
	TenantId string `json:"tenantId,omitempty"`
}

func (c *ApiKeyCredentials) GetSubType() string {
	return "ApiKey"
}

// AzureServicePrincipalCredentials represents the AzureServicePrincipal credential subtype
type AzureServicePrincipalCredentials struct {
	ClientId     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	TenantId     string `json:"tenantId,omitempty"`
}

func (c *AzureServicePrincipalCredentials) GetSubType() string {
	return "AzureServicePrincipal"
}

// ConnectionStringCredentials represents the ConnectionString credential subtype
type ConnectionStringCredentials struct {
	ConnectionString string `json:"connectionString,omitempty"`
}

func (c *ConnectionStringCredentials) GetSubType() string {
	return "ConnectionString"
}

// PasskeyCredentials represents the Passkey credential subtype
type PasskeyCredentials struct {
	PasskeyPrivateKey string `json:"passkeyPrivateKey,omitempty"`
	PasskeyRpID       string `json:"passkeyRpID,omitempty"`
}

func (c *PasskeyCredentials) GetSubType() string {
	return "Passkey"
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
