package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EntryUserCredentialService service

// EntryUserCredential represents a DVLS entry/connection.
type EntryUserCredential struct {
	ID                string                  `json:"id,omitempty"`
	VaultId           string                  `json:"repositoryId"`
	EntryName         string                  `json:"name"`
	Description       string                  `json:"description"`
	EntryFolderPath   string                  `json:"group"`
	ModifiedDate      *ServerTime             `json:"modifiedDate,omitempty"`
	ConnectionType    ServerConnectionType    `json:"connectionType"`
	ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
	Tags              []string                `json:"keywords,omitempty"`

	Credentials EntryUserAuthDetails `json:"data,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e EntryUserCredential) MarshalJSON() ([]byte, error) {
	raw := struct {
		Id           string `json:"id,omitempty"`
		RepositoryId string `json:"repositoryId"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Events       struct {
			OpenCommentPrompt                        bool `json:"openCommentPrompt"`
			CredentialViewedPrompt                   bool `json:"credentialViewedPrompt"`
			TicketNumberIsRequiredOnCredentialViewed bool `json:"ticketNumberIsRequiredOnCredentialViewed"`
			TicketNumberIsRequiredOnClose            bool `json:"ticketNumberIsRequiredOnClose"`
			CredentialViewedCommentIsRequired        bool `json:"credentialViewedCommentIsRequired"`
			TicketNumberIsRequiredOnOpen             bool `json:"ticketNumberIsRequiredOnOpen"`
			CloseCommentIsRequired                   bool `json:"closeCommentIsRequired"`
			OpenCommentPromptOnBrowserExtensionLink  bool `json:"openCommentPromptOnBrowserExtensionLink"`
			CloseCommentPrompt                       bool `json:"closeCommentPrompt"`
			OpenCommentIsRequired                    bool `json:"openCommentIsRequired"`
			WarnIfAlreadyOpened                      bool `json:"warnIfAlreadyOpened"`
		} `json:"events"`
		Data              string                  `json:"data"`
		Expiration        string                  `json:"expiration"`
		CheckOutMode      int                     `json:"checkOutMode"`
		Group             string                  `json:"group"`
		ConnectionType    ServerConnectionType    `json:"connectionType"`
		ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
		Keywords          string                  `json:"keywords"`
	}{}

	raw.Id = e.ID
	raw.Keywords = sliceToKeywords(e.Tags)
	raw.Description = e.Description
	raw.RepositoryId = e.VaultId
	raw.Group = e.EntryFolderPath
	raw.ConnectionSubType = e.ConnectionSubType
	raw.ConnectionType = e.ConnectionType
	raw.Name = e.EntryName
	sensitiveJson, err := json.Marshal(e.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sensitive data. error: %w", err)
	}

	raw.Data = string(sensitiveJson)

	entryJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return entryJson, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *EntryUserCredential) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			ID                string
			Description       string
			Name              string
			Group             string
			Username          string
			ModifiedDate      *ServerTime
			Keywords          string
			RepositoryId      string
			ConnectionType    ServerConnectionType
			ConnectionSubType ServerConnectionSubType
		}
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	e.ID = raw.Data.ID
	e.EntryName = raw.Data.Name
	e.ConnectionType = raw.Data.ConnectionType
	e.ConnectionSubType = raw.Data.ConnectionSubType
	e.ModifiedDate = raw.Data.ModifiedDate
	e.Credentials.Username = raw.Data.Username
	e.Description = raw.Data.Description
	e.EntryFolderPath = raw.Data.Group
	e.VaultId = raw.Data.RepositoryId

	e.Tags = keywordsToSlice(raw.Data.Keywords)

	return nil
}

// EntryUserAuthDetails represents an Entry User Authentication Details fields.
type EntryUserAuthDetails struct {
	Username string
	Password *string
}

// MarshalJSON implements the json.Marshaler interface.
func (s EntryUserAuthDetails) MarshalJSON() ([]byte, error) {
	raw := struct {
		AllowClipboard         bool    `json:"allowClipboard"`
		CredentialConnectionId string  `json:"credentialConnectionId"`
		PamCredentialId        string  `json:"pamCredentialId"`
		PamCredentialName      string  `json:"pamCredentialName"`
		CredentialMode         int     `json:"credentialMode"`
		Credentials            *string `json:"credentials"`
		Domain                 string  `json:"domain"`
		MnemonicPassword       string  `json:"mnemonicPassword"`
		PasswordItem           struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData"`
		} `json:"passwordItem"`
		PromptForPassword bool   `json:"promptForPassword"`
		UserName          string `json:"userName"`
	}{}

	if s.Password != nil {
		raw.PasswordItem.HasSensitiveData = true
		raw.PasswordItem.SensitiveData = *s.Password
	}
	raw.UserName = s.Username

	secretJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return secretJson, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *EntryUserAuthDetails) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data string
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	if raw.Data != "" {
		newRaw := struct {
			Data struct {
				Credentials struct {
					Username string
					Password string
				}
			}
		}{}
		err = json.Unmarshal([]byte(raw.Data), &newRaw)
		if err != nil {
			return err
		}

		s.Username = newRaw.Data.Credentials.Username
		s.Password = &newRaw.Data.Credentials.Password
	}

	return nil
}

// GetUserAuthDetails returns entry with the entry.Credentials.Password field.
func (c *EntryUserCredentialService) GetUserAuthDetails(entry EntryUserCredential) (EntryUserCredential, error) {
	var secret EntryUserAuthDetails
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entry.ID, "/sensitive-data")
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryUserCredential{}, err
	}

	err = json.Unmarshal(resp.Response, &secret)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	entry.Credentials = secret

	return entry, nil
}

// Get returns a single Entry specified by entryId. Call GetEntryCredentialsPassword with
// the returned Entry to fetch the password.
func (c *EntryUserCredentialService) Get(entryId string) (EntryUserCredential, error) {
	var entry EntryUserCredential
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while fetching entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryUserCredential{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return entry, nil
}

// New creates a new EntryUserCredential based on entry.
func (c *EntryUserCredentialService) New(entry EntryUserCredential) (EntryUserCredential, error) {
	if entry.ConnectionType != ServerConnectionCredential || entry.ConnectionSubType != ServerConnectionSubTypeDefault {
		return EntryUserCredential{}, fmt.Errorf("unsupported entry type (%s %s). Only %s %s is supported", entry.ConnectionType, entry.ConnectionSubType, ServerConnectionCredential, ServerConnectionSubTypeDefault)
	}

	entry.ID = ""
	entry.ModifiedDate = nil

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
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
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryUserCredential{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return entry, nil
}

// Update updates an EntryUserCredential based on entry. Will replace all other fields whether included or not.
func (c *EntryUserCredentialService) Update(entry EntryUserCredential) (EntryUserCredential, error) {
	if entry.ConnectionType != ServerConnectionCredential || entry.ConnectionSubType != ServerConnectionSubTypeDefault {
		return EntryUserCredential{}, fmt.Errorf("unsupported entry type (%s %s). Only %s %s is supported", entry.ConnectionType, entry.ConnectionSubType, ServerConnectionCredential, ServerConnectionSubTypeDefault)
	}
	_, err := c.Get(entry.ID)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while fetching entry. error: %w", err)
	}

	entry.ModifiedDate = nil

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, "save")
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPut, bytes.NewBuffer(entryJson))
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("error while creating entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryUserCredential{}, err
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return EntryUserCredential{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return entry, nil
}

// Delete deletes an EntryUserCredential based on entryId.
func (c *EntryUserCredentialService) Delete(entryId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return fmt.Errorf("failed to delete entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting entry. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// NewEntryUserAuthDetails returns an EntryUserAuthDetails with an initialised EntryUserAuthDetails.Password.
func (c *EntryUserCredentialService) NewUserAuthDetails(username string, password string) EntryUserAuthDetails {
	creds := EntryUserAuthDetails{
		Username: username,
		Password: &password,
	}
	return creds
}
