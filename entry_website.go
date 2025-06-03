package dvls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EntryWebsiteService service

// EntryWebsite represents a website entry in DVLS
type EntryWebsite struct {
	Id                string                  `json:"id,omitempty"`
	VaultId           string                  `json:"repositoryId"`
	EntryName         string                  `json:"name"`
	Description       string                  `json:"description"`
	EntryFolderPath   string                  `json:"group"`
	ModifiedDate      *ServerTime             `json:"modifiedDate,omitempty"`
	ConnectionType    ServerConnectionType    `json:"connectionType"`
	ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
	Tags              []string                `json:"keywords,omitempty"`

	WebsiteDetails EntryWebsiteAuthDetails `json:"data"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e EntryWebsite) MarshalJSON() ([]byte, error) {
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

	raw.Id = e.Id
	raw.Keywords = sliceToKeywords(e.Tags)
	raw.Description = e.Description
	raw.RepositoryId = e.VaultId
	raw.Group = e.EntryFolderPath
	raw.ConnectionSubType = e.ConnectionSubType
	raw.ConnectionType = e.ConnectionType
	raw.Name = e.EntryName
	sensitiveJson, err := json.Marshal(e.WebsiteDetails)
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
func (e *EntryWebsite) UnmarshalJSON(d []byte) error {
	raw := struct {
		Id                string                  `json:"id"`
		Description       string                  `json:"description"`
		Name              string                  `json:"name"`
		Group             string                  `json:"group"`
		ModifiedDate      *ServerTime             `json:"modifiedDate"`
		Keywords          string                  `json:"keywords"`
		RepositoryId      string                  `json:"repositoryId"`
		ConnectionType    ServerConnectionType    `json:"connectionType"`
		ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
		Data              json.RawMessage         `json:"data"`
	}{}

	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	e.Id = raw.Id
	e.EntryName = raw.Name
	e.ConnectionType = raw.ConnectionType
	e.ConnectionSubType = raw.ConnectionSubType
	e.ModifiedDate = raw.ModifiedDate
	e.Description = raw.Description
	e.EntryFolderPath = raw.Group
	e.VaultId = raw.RepositoryId
	e.Tags = keywordsToSlice(raw.Keywords)

	if len(raw.Data) > 0 {
		if err := json.Unmarshal(raw.Data, &e.WebsiteDetails); err != nil {
			return fmt.Errorf("failed to unmarshal website details: %w", err)
		}
	}

	return nil
}

// EntryWebsiteAuthDetails represents website-specific fields
type EntryWebsiteAuthDetails struct {
	Username              string
	Password              *string
	URL                   string
	WebBrowserApplication int
}

// MarshalJSON implements the json.Marshaler interface.
func (s EntryWebsiteAuthDetails) MarshalJSON() ([]byte, error) {
	raw := struct {
		AutoFillLogin         bool   `json:"AutoFillLogin"`
		AutoSubmit            bool   `json:"AutoSubmit"`
		AutomaticRefreshTime  int    `json:"AutomaticRefreshTime"`
		ChromeProxyType       int    `json:"ChromeProxyType"`
		CustomJavaScript      string `json:"CustomJavaScript"`
		Host                  string `json:"Host"`
		URL                   string `json:"URL"`
		Username              string `json:"Username"`
		WebBrowserApplication int    `json:"WebBrowserApplication"`
		PasswordItem          struct {
			HasSensitiveData bool   `json:"HasSensitiveData"`
			SensitiveData    string `json:"SensitiveData"`
		} `json:"PasswordItem"`
		VPN struct {
			EnableAutoDetectIsOnlineVPN int `json:"EnableAutoDetectIsOnlineVPN"`
		} `json:"VPN"`
	}{}

	if s.Password != nil {
		raw.PasswordItem.HasSensitiveData = true
		raw.PasswordItem.SensitiveData = *s.Password
	} else {
		raw.PasswordItem.HasSensitiveData = false
	}

	raw.Username = s.Username
	raw.URL = s.URL
	raw.WebBrowserApplication = s.WebBrowserApplication

	secretJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return secretJson, nil
}

// GetWebsiteDetails returns entry with the entry.WebsiteDetails.Password field.
func (c *EntryWebsiteService) GetWebsiteDetails(entry EntryWebsite) (EntryWebsite, error) {
	var respData struct {
		Data string `json:"data"`
	}

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entry.Id, "/sensitive-data")
	if err != nil {
		return EntryWebsite{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return EntryWebsite{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntryWebsite{}, err
	}

	if err := json.Unmarshal(resp.Response, &respData); err != nil {
		return EntryWebsite{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	var sensitiveDataResponse struct {
		Data struct {
			PasswordItem struct {
				HasSensitiveData bool    `json:"hasSensitiveData"`
				SensitiveData    *string `json:"sensitiveData,omitempty"`
			} `json:"passwordItem"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(respData.Data), &sensitiveDataResponse); err != nil {
		return EntryWebsite{}, fmt.Errorf("failed to unmarshal inner data. error: %w", err)
	}

	if sensitiveDataResponse.Data.PasswordItem.HasSensitiveData {
		entry.WebsiteDetails.Password = sensitiveDataResponse.Data.PasswordItem.SensitiveData
	} else {
		entry.WebsiteDetails.Password = nil
	}

	return entry, nil
}

// Get returns a single Entry specified by entryId. Call GetWebsiteDetails with
// the returned Entry to fetch the password.
func (s *EntryWebsiteService) Get(entryId string) (EntryWebsite, error) {
	var respData struct {
		Data EntryWebsite `json:"data"`
	}

	reqUrl, err := url.JoinPath(s.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return EntryWebsite{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := s.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntryWebsite{}, fmt.Errorf("error fetching entry: %w", err)
	}
	if err = resp.CheckRespSaveResult(); err != nil {
		return EntryWebsite{}, err
	}
	if resp.Response == nil {
		return EntryWebsite{}, fmt.Errorf("response body is nil for request to %s", reqUrl)
	}

	if err := json.Unmarshal(resp.Response, &respData); err != nil {
		return EntryWebsite{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return respData.Data, nil
}
