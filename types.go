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

type loginResponse struct {
	Data struct {
		Message string
		Result  ServerLoginResult
		TokenId string
	}
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
		s.Password = newRaw.Data.Credentials.Password
	}

	return nil
}

//go:generate stringer -type=ServerLoginResult -trimprefix ServerLogin
type ServerLoginResult int

const (
	ServerLoginError ServerLoginResult = iota
	ServerLoginSuccess
	ServerLoginInvalidUserNamePassword
	ServerLoginInvalidDataSource
	ServerLoginDisabledDataSource
	ServerLoginInvalidSubscription
	ServerLoginTooManyUserForTheLicense
	ServerLoginExpiredSubscription
	ServerLoginInGracePeriod
	ServerLoginDisabledUser
	ServerLoginUserNotFound
	ServerLoginLockedUser
	ServerLoginNotApprovedUser
	ServerLoginBlackListed
	ServerLoginInvalidIP
	ServerLoginUnableToCreateUser
	ServerLoginTwoFactorTypeNotConfigured
	ServerLoginTwoFactorTypeActivatedNotAllowedClientSide
	ServerLoginDomainNotTrusted
	ServerLoginUserDoesNotBelongToDefaultDomain
	ServerLoginInvalidGeoIP
	ServerLoginTwoFactorIsRequired
	ServerLoginTwoFactorPreconfigured
	ServerLoginTwoFactorSecondStepIsRequired
	ServerLoginTwoFactorUserIsDenied
	ServerLoginTwoFactorSmsSended
	ServerLoginTwoFactorTimeout
	ServerLoginTwoFactorUserLockedOut
	ServerLoginTwoFactorUserFraud
	ServerLoginTwoFactorUserEmailNotConfigured
	ServerLoginTwoFactorUserSmsNotConfigured
	ServerLoginNotInTrustedGroup
	ServerLoginServerNotResponding
	ServerLoginNotAccessToApplication
	ServerLoginDirectoryNotResponding
	ServerLoginWindowsAuthenticationFailure
	ServerLoginForcePasswordChange
	ServerLoginTwoFactorInvalid
	ServerLoginOutsideValidUsageTimePeriod
)
