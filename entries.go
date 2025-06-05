package dvls

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	entryEndpoint            string = "/api/connections/partial"
	entryConnectionsEndpoint string = "/api/connections"
	entryBasePublicEndpoint  string = "/api/v1/vault/{vaultId}/entry"
	entryPublicEndpoint      string = "/api/v1/vault/{vaultId}/entry/{id}"
)

type Entries struct {
	Certificate *EntryCertificateService
	Host        *EntryHostService
	Credential  *EntryCredentialService
	Website     *EntryWebsiteService
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
		return fmt.Errorf("unsupported entry type/subtype: %s", key)
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
