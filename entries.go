package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DvlsEntry struct {
	ID                string                  `json:"id,omitempty"`
	VaultId           string                  `json:"repositoryId"`
	EntryName         string                  `json:"name"`
	Description       string                  `json:"description"`
	EntryFolderPath   string                  `json:"group"`
	ModifiedDate      *time.Time              `json:"modifiedDate,omitempty"`
	ConnectionType    ServerConnectionType    `json:"connectionType"`
	ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
	Tags              []string                `json:"keywords,omitempty"`

	Credentials DvlsEntryCredentials `json:"data,omitempty"`
}

func (e DvlsEntry) MarshalJSON() ([]byte, error) {
	raw := struct {
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

	for i, v := range e.Tags {
		if strings.Contains(v, " ") {
			e.Tags[i] = "\"" + v + "\""
		}
	}

	raw.Keywords = strings.Join(e.Tags, " ")
	raw.Description = e.Description
	raw.RepositoryId = e.VaultId
	raw.Group = e.EntryFolderPath
	raw.ConnectionSubType = e.ConnectionSubType
	raw.ConnectionType = e.ConnectionType
	raw.Name = e.EntryName
	sensitiveJson, err := json.Marshal(e.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall sensitive data. error: %w", err)
	}

	raw.Data = string(sensitiveJson)

	entryJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return entryJson, nil
}

func (e *DvlsEntry) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			ID                string
			Description       string
			Name              string
			Group             string
			Username          string
			ModifiedDate      string
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

	var date *time.Time
	if raw.Data.ModifiedDate != "" {
		dateParsed, err := time.Parse("2006-01-02T15:04:05", raw.Data.ModifiedDate)
		if err != nil {
			return err
		}
		date = &dateParsed
	}

	e.ID = raw.Data.ID
	e.EntryName = raw.Data.Name
	e.ConnectionType = raw.Data.ConnectionType
	e.ConnectionSubType = raw.Data.ConnectionSubType
	e.ModifiedDate = date
	e.Credentials.Username = raw.Data.Username
	e.Description = raw.Data.Description
	e.EntryFolderPath = raw.Data.Group
	e.VaultId = raw.Data.RepositoryId

	var spacedTag bool
	tags := strings.FieldsFunc(raw.Data.Keywords, func(r rune) bool {
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

	e.Tags = tags

	return nil
}

type DvlsEntryCredentials struct {
	Username string
	Password *string
}

func (s DvlsEntryCredentials) MarshalJSON() ([]byte, error) {
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

func (s *DvlsEntryCredentials) UnmarshalJSON(d []byte) error {
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

const (
	entryEndpoint string = "/api/connections/partial"
)

func (c *Client) GetEntryCredentialsPassword(entry DvlsEntry) (DvlsEntry, error) {
	var secret DvlsEntryCredentials
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entry.ID, "/sensitive-data")
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if resp.Result != 1 {
		return DvlsEntry{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", resp.Result)
	}

	err = json.Unmarshal(resp.Response, &secret)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	entry.Credentials = secret

	return entry, nil
}

func (c *Client) GetEntry(entryId string) (DvlsEntry, error) {
	var entry DvlsEntry
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("error while fetching entry. error: %w", err)
	} else if resp.Result != 1 {
		return DvlsEntry{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", resp.Result)
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}

func (c *Client) NewEntry(entry DvlsEntry) (DvlsEntry, error) {
	if entry.ConnectionType != ServerConnectionCredential || entry.ConnectionSubType != ServerConnectionSubTypeDefault {
		return DvlsEntry{}, fmt.Errorf("unsupported entry type (%s %s). Only %s %s is supported", entry.ConnectionType, entry.ConnectionSubType, ServerConnectionCredential, ServerConnectionSubTypeDefault)
	}

	entry.ID = ""
	entry.ModifiedDate = nil

	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, "save")
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	entryJson, err := json.Marshal(entry)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to marshall body. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, bytes.NewBuffer(entryJson))
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("error while creating entry. error: %w", err)
	} else if resp.Result != 1 {
		return DvlsEntry{}, fmt.Errorf("unexpected result code %d %s", resp.Result, resp.Message)
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}
