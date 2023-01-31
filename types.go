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

type DvlsEntry struct {
	ID                string
	Name              string
	ConnectionType    ServerConnectionType
	ConnectionSubType ServerConnectionSubType
}

func (e *DvlsEntry) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			ID                string
			Name              string
			ConnectionType    ServerConnectionType
			ConnectionSubType ServerConnectionSubType
		}
		Result int
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	e.ID = raw.Data.ID
	e.Name = raw.Data.Name
	e.ConnectionType = raw.Data.ConnectionType
	e.ConnectionSubType = raw.Data.ConnectionSubType

	return nil
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

type DvlsUser struct {
	ID       string
	Username string
	UserType UserAuthenticationType
	result   ServerLoginResult
	message  string
	tokenId  string
}

func (u *DvlsUser) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			TokenId    string
			UserEntity struct {
				Id           string
				Display      string
				UserSecurity struct {
					AuthenticationType UserAuthenticationType
				}
			}
		}
		Result  ServerLoginResult
		Message string
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	u.ID = raw.Data.UserEntity.Id
	u.Username = raw.Data.UserEntity.Display
	u.UserType = raw.Data.UserEntity.UserSecurity.AuthenticationType
	u.result = raw.Result
	u.tokenId = raw.Data.TokenId
	u.message = raw.Message

	return nil
}

//go:generate stringer -type=UserAuthenticationType -trimprefix UserAuthentication
type UserAuthenticationType int

const (
	UserAuthenticationBuiltin UserAuthenticationType = iota
	UserAuthenticationLocalWindows
	UserAuthenticationSqlServer
	UserAuthenticationDomain
	UserAuthenticationOffice365
	UserAuthenticationNone
	UserAuthenticationCloud
	UserAuthenticationLegacy
	UserAuthenticationAzureAD
	UserAuthenticationApplication
	UserAuthenticationOkta
)

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

//go:generate stringer -type=ServerConnectionType -trimprefix ServerConnection
type ServerConnectionType int

const (
	ServerConnectionUndefined ServerConnectionType = iota
	ServerConnectionRDPConfigured
	ServerConnectionRDPFilename
	ServerConnectionCommandLine
	ServerConnectionVNC
	ServerConnectionWebBrowser
	ServerConnectionLogMeIn
	ServerConnectionTeamViewer
	ServerConnectionPutty
	ServerConnectionFtp
	ServerConnectionVirtualPC
	ServerConnectionRadmin
	ServerConnectionDameware
	ServerConnectionVMWare
	ServerConnectionPCAnywhere
	ServerConnectionICA
	ServerConnectionXWindow
	ServerConnectionHyperV
	ServerConnectionAddOn
	ServerConnectionRemoteAssistance
	ServerConnectionVPN
	ServerConnectionVirtualBox
	ServerConnectionVMRC
	ServerConnectionXenServer
	ServerConnectionWindowsVirtualPC
	ServerConnectionGroup
	ServerConnectionCredential
	ServerConnectionHpRgs
	ServerConnectionDesktone
	ServerConnectionApplicationTool
	ServerConnectionSessionTool
	ServerConnectionContact
	ServerConnectionDataEntry
	ServerConnectionDataReport
	ServerConnectionAgent
	ServerConnectionComputer
	ServerConnectionDropBox
	ServerConnectionS3
	ServerConnectionAzureStorage
	ServerConnectionCitrixWeb
	ServerConnectionPowerShell
	ServerConnectionHostSessionTool
	ServerConnectionShortcut
	ServerConnectionIntelAMT
	ServerConnectionAzure
	ServerConnectionDocument
	ServerConnectionVMWareConsole
	ServerConnectionInventoryReport
	ServerConnectionSkyDrive
	ServerConnectionScreenConnect
	ServerConnectionAzureTableStorage
	ServerConnectionAzureQueueStorage
	ServerConnectionTemplateGroup
	ServerConnectionHost
	ServerConnectionDatabase
	ServerConnectionCustomer
	ServerConnectionADConsole
	ServerConnectionAws
	ServerConnectionSNMPReport
	ServerConnectionSync
	ServerConnectionGateway
	ServerConnectionPlayList
	ServerConnectionTerminalConsole
	ServerConnectionPSExec
	ServerConnectionAppleRemoteDesktop
	ServerConnectionSpiceworks
	ServerConnectionDeskRoll
	ServerConnectionSecureCRT
	ServerConnectionIterm
	ServerConnectionSheet
	ServerConnectionSplunk
	ServerConnectionPortForward
	ServerConnectionTeamViewerConsole
	ServerConnectionScreenHero
	ServerConnectionTelnet
	ServerConnectionSerial
	ServerConnectionSSHTunnel
	ServerConnectionSSHShell
	ServerConnectionResetPassword
	ServerConnectionWayk
	ServerConnectionControlUp
	ServerConnectionDataSource
	ServerConnectionChromeRemoteDesktop
	ServerConnectionRDCommander
	ServerConnectionIDrac
	ServerConnectionIlo
	ServerConnectionWebDav
	ServerConnectionBeyondTrustPasswordSafeConsole
	ServerConnectionDevolutionsProxy
	ServerConnectionFtpNative
	ServerConnectionPowerShellRemoteConsole
	ServerConnectionProxyTunnel
	ServerConnectionRoot
	ServerConnectionBeyondTrustPasswordSafe
	ServerConnectionFileExplorer
	ServerConnectionScp
	ServerConnectionSftp
	ServerConnectionAzureBlobStorage
	ServerConnectionTFtp
	ServerConnectionGoToAssist
	ServerConnectionIPTable
	ServerConnectionHub
	ServerConnectionGoogleDrive
	ServerConnectionGoogleCloud
	ServerConnectionNoVNC
	ServerConnectionSplashtop
	ServerConnectionJumpDesktop
	ServerConnectionBoxNet
	ServerConnectionMSPAnywhere
	ServerConnectionRepository
	ServerConnectionCyberArkPSM
	ServerConnectionCloudBerryRemoteAssistant
	ServerConnectionITGlue
	ServerConnectionSmartFolder
	ServerConnectionCyberArkJump
	ServerConnectionWindowsAdminCenter
	ServerConnectionDevolutionsGateway
	ServerConnectionWaykDenConsole
	ServerConnectionRDGatewayConsole
	ServerConnectionCyberArkDashboard
	ServerConnectionDVLSPamDashboard
	ServerConnectionSMB
	ServerConnectionAppleRemoteManagement
	ServerConnectionRustDesk
	ServerConnectionPAM
	ServerConnectionITManager
	ServerConnectionCustomImage
)

type ServerConnectionSubType string

const (
	ServerConnectionSubTypeDefault ServerConnectionSubType = "Default"
)
