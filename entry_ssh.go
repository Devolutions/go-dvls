package dvls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type EntrySSHService service

// EntrySSH represents a ssh entry in DVLS
type EntrySSH struct {
	ID                string                  `json:"id,omitempty"`
	VaultId           string                  `json:"repositoryId"`
	EntryName         string                  `json:"name"`
	Description       string                  `json:"description"`
	EntryFolderPath   string                  `json:"group"`
	ModifiedDate      *ServerTime             `json:"modifiedDate,omitempty"`
	ConnectionType    ServerConnectionType    `json:"connectionType"`
	ConnectionSubType ServerConnectionSubType `json:"connectionSubType"`
	Tags              []string                `json:"keywords,omitempty"`

	SSHDetails EntrySSHAuthDetails `json:"data"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e EntrySSH) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID           string `json:"id,omitempty"`
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

	raw.ID = e.ID
	raw.Keywords = sliceToKeywords(e.Tags)
	raw.Description = e.Description
	raw.RepositoryId = e.VaultId
	raw.Group = e.EntryFolderPath
	raw.ConnectionSubType = e.ConnectionSubType
	raw.ConnectionType = e.ConnectionType
	raw.Name = e.EntryName
	sensitiveJson, err := json.Marshal(e.SSHDetails)
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
func (e *EntrySSH) UnmarshalJSON(d []byte) error {
	raw := struct {
		ID                string                  `json:"id"`
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

	e.ID = raw.ID
	e.EntryName = raw.Name
	e.ConnectionType = raw.ConnectionType
	e.ConnectionSubType = raw.ConnectionSubType
	e.ModifiedDate = raw.ModifiedDate
	e.Description = raw.Description
	e.EntryFolderPath = raw.Group
	e.VaultId = raw.RepositoryId
	e.Tags = keywordsToSlice(raw.Keywords)

	if len(raw.Data) > 0 {
		if err := json.Unmarshal(raw.Data, &e.SSHDetails); err != nil {
			return fmt.Errorf("failed to unmarshal ssh details: %w", err)
		}
	}

	return nil
}

// EntrySSHAuthDetails represents ssh-specific fields
type EntrySSHAuthDetails struct {
	Username   string  `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	Host       string  `json:"host"`
	PrivateKey *string `json:"privateKey,omitempty"`
	Passphrase *string `json:"passphrase,omitempty"`
	Port       int     `json:"port"`
}

type SensitiveItem struct {
	HasSensitiveData bool    `json:"hasSensitiveData"`
	SensitiveData    *string `json:"sensitiveData,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (s EntrySSHAuthDetails) MarshalJSON() ([]byte, error) {
	raw := struct {
		BeforeDisconnectMacrosMore []interface{} `json:"beforeDisconnectMacrosMore"`
		AfterConnectMacros         []interface{} `json:"afterConnectMacros"`
		BeforeDisconnectMacros     []interface{} `json:"beforeDisconnectMacros"`
		CloseOnDisconnect          bool          `json:"closeOnDisconnect"`
		DisconnectAction           int           `json:"disconnectAction"`
		Host                       string        `json:"host"`
		HostPort                   int           `json:"hostPort"`
		PasswordItem               struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"passwordItem"`
		PortForwards   []interface{} `json:"portForwards"`
		PrivateKeyData struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"privateKeyData"`
		PrivateKeyType    int `json:"privateKeyType"`
		ProxyPasswordItem struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"proxyPasswordItem"`
		ProxyType                               int           `json:"proxyType"`
		SSHGatewayCredentialSource              int           `json:"sshGatewayCredentialSource"`
		SSHGatewayPort                          int           `json:"sshGatewayPort"`
		SSHGatewayPrivateKeyPromptForPassPhrase bool          `json:"sshGatewayPrivateKeyPromptForPassPhrase"`
		SSHGateways                             []interface{} `json:"sshGateways"`
		X11Protocol                             int           `json:"x11Protocol"`
		AfterConnectMacrosMore                  []interface{} `json:"afterConnectMacrosMore"`
		CredentialConnectionID                  string        `json:"credentialConnectionId"`
		CredentialMode                          int           `json:"credentialMode"`
		VPN                                     struct {
			EnableAutoDetectIsOnlineVPN           int    `json:"enableAutoDetectIsOnlineVPN"`
			TerminalHost                          string `json:"terminalHost"`
			TerminalHostPort                      int    `json:"terminalHostPort"`
			TerminalPrivateKeyType                int    `json:"terminalPrivateKeyType"`
			TerminalPrivateKeyPromptForPassPhrase bool   `json:"terminalPrivateKeyPromptForPassPhrase"`
			TerminalDisconnectAction              int    `json:"terminalDisconnectAction"`
			TerminalShowLogs                      bool   `json:"terminalShowLogs"`
		} `json:"vpn"`
		PrivateKeyPassPhraseItem struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"privateKeyPassPhraseItem"`
	}{}

	if s.Password != nil {
		raw.PasswordItem.HasSensitiveData = true
		raw.PasswordItem.SensitiveData = *s.Password
	} else {
		raw.PasswordItem.HasSensitiveData = false
	}

	if s.PrivateKey != nil {
		raw.PrivateKeyData.HasSensitiveData = true
		raw.PrivateKeyData.SensitiveData = *s.PrivateKey
	} else {
		raw.PrivateKeyData.HasSensitiveData = false
	}

	if s.Passphrase != nil {
		raw.PrivateKeyPassPhraseItem.HasSensitiveData = true
		raw.PrivateKeyPassPhraseItem.SensitiveData = *s.Passphrase
	} else {
		raw.PrivateKeyPassPhraseItem.HasSensitiveData = false
	}

	//raw.Username = s.Username
	raw.Host = s.Host
	raw.HostPort = s.Port

	secretJson, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return secretJson, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for EntrySSHAuthDetails.
func (s *EntrySSHAuthDetails) UnmarshalJSON(data []byte) error {
	raw := struct {
		Username     string `json:"username"`
		Host         string `json:"host"`
		HostPort     int    `json:"hostPort"`
		PasswordItem struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"passwordItem"`
		PrivateKeyData struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"privateKeyData"`
		PrivateKeyPassPhraseItem struct {
			HasSensitiveData bool   `json:"hasSensitiveData"`
			SensitiveData    string `json:"sensitiveData,omitempty"`
		} `json:"privateKeyPassPhraseItem"`
	}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Username = raw.Username
	s.Host = raw.Host
	s.Port = raw.HostPort

	if raw.PasswordItem.HasSensitiveData {
		s.Password = &raw.PasswordItem.SensitiveData
	} else {
		s.Password = nil
	}

	if raw.PrivateKeyData.HasSensitiveData {
		s.PrivateKey = &raw.PrivateKeyData.SensitiveData
	} else {
		s.PrivateKey = nil
	}

	if raw.PrivateKeyPassPhraseItem.HasSensitiveData {
		s.Passphrase = &raw.PrivateKeyPassPhraseItem.SensitiveData
	} else {
		s.Passphrase = nil
	}

	return nil
}

// GetSSHDetails returns entry with the entry.SSHDetails.Password field.
func (c *EntrySSHService) GetSSHDetails(entry EntrySSH) (EntrySSH, error) {
	var respData struct {
		Data string `json:"data"`
	}

	reqUrl, err := url.JoinPath(c.client.baseUri, entryEndpoint, entry.ID, "/sensitive-data")
	if err != nil {
		return EntrySSH{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.client.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return EntrySSH{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return EntrySSH{}, err
	}

	if err := json.Unmarshal(resp.Response, &respData); err != nil {
		return EntrySSH{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	var sensitiveDataResponse struct {
		Data struct {
			PasswordItem             SensitiveItem `json:"passwordItem"`
			PrivateKeyData           SensitiveItem `json:"privateKeyData"`
			PrivateKeyPassPhraseItem SensitiveItem `json:"privateKeyPassPhraseItem"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(respData.Data), &sensitiveDataResponse); err != nil {
		return EntrySSH{}, fmt.Errorf("failed to unmarshal inner data. error: %w", err)
	}

	if sensitiveDataResponse.Data.PasswordItem.HasSensitiveData {
		entry.SSHDetails.Password = sensitiveDataResponse.Data.PasswordItem.SensitiveData
	} else {
		entry.SSHDetails.Password = nil
	}

	if sensitiveDataResponse.Data.PrivateKeyData.HasSensitiveData {
		entry.SSHDetails.PrivateKey = sensitiveDataResponse.Data.PrivateKeyData.SensitiveData
	} else {
		entry.SSHDetails.PrivateKey = nil
	}

	if sensitiveDataResponse.Data.PrivateKeyPassPhraseItem.HasSensitiveData {
		entry.SSHDetails.Passphrase = sensitiveDataResponse.Data.PrivateKeyPassPhraseItem.SensitiveData
	} else {
		entry.SSHDetails.Passphrase = nil
	}

	return entry, nil
}

// Get returns a single Entry specified by entryId. Call GetSSHDetails with
// the returned Entry to fetch the password.
func (s *EntrySSHService) Get(entryId string) (EntrySSH, error) {
	var respData struct {
		Data EntrySSH `json:"data"`
	}

	reqUrl, err := url.JoinPath(s.client.baseUri, entryEndpoint, entryId)
	if err != nil {
		return EntrySSH{}, fmt.Errorf("failed to build entry url: %w", err)
	}

	resp, err := s.client.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return EntrySSH{}, fmt.Errorf("error fetching entry: %w", err)
	}

	fmt.Println("Raw response:", string(resp.Response))

	if err = resp.CheckRespSaveResult(); err != nil {
		return EntrySSH{}, err
	}
	if resp.Response == nil {
		return EntrySSH{}, fmt.Errorf("response body is nil for request to %s", reqUrl)
	}

	if err := json.Unmarshal(resp.Response, &respData); err != nil {
		return EntrySSH{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return respData.Data, nil
}
