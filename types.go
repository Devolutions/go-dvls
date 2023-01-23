package dvls

import (
	"encoding/json"
	"net/http"
)

type Client struct {
	client     *http.Client
	baseUri    string
	credential credentials
}

type credentials struct {
	username string
	password string
	token    string
}

type loginReqBody struct {
	Username        string          `json:"userName"`
	LoginParameters loginParameters `json:"LoginParameters"`
}

type loginParameters struct {
	Password         string `json:"Password"`
	Client           string `json:"Client"`
	Version          string `json:"Version,omitempty"`
	LocalMachineName string `json:"LocalMachineName,omitempty"`
	LocalUserName    string `json:"LocalUserName,omitempty"`
}

type DvlsSecret struct {
	ID       string
	Username string
	Password string
}

func (s *DvlsSecret) UnmarshalJSON(d []byte) error {
	var raw dvlsSecretRaw

	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}
	s.ID = raw.Data.ID
	s.Username = raw.Data.Data.Username
	if len(raw.Data.MetaInformationData.MetaInformation.PasswordHistory) > 0 {
		s.Password = raw.Data.MetaInformationData.MetaInformation.PasswordHistory[0].Password
	}

	return nil
}

type dvlsSecretRaw struct {
	Data struct {
		ID                  string `json:"id,omitempty"`
		MetaInformationData struct {
			MetaInformation struct {
				PasswordHistory []struct {
					Password string `json:"password,omitempty"`
				} `json:"passwordHistory,omitempty"`
			} `json:"metaInformation,omitempty"`
		} `json:"metaInformationData,omitempty"`
		Data struct {
			Username string `json:"username,omitempty"`
		} `json:"data,omitempty"`
	} `json:"data,omitempty"`
	Result int
}
