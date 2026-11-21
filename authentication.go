package dvls

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// Client represents the DVLS client used to communicate with the API.
type Client struct {
	client     *http.Client
	baseUri    string
	credential *credentials

	common service

	Entries                       *Entries
	Vaults                        *Vaults
	AdministrativeRoles           *AdministrativeRoles
	AdministrativeRoleAssignments *AdministrativeRoleAssignments
	Users                         *Users
	UserGroups                    *UserGroups
}

type service struct {
	client *Client
}

// credentials is shared by pointer: NewClient returns the Client by value, and every copy must see the token
// a re-login stores.
type credentials struct {
	appKey    string
	appSecret string
	apiKey    string

	mu    sync.Mutex
	token string
}

func (c *credentials) useApiKey() bool {
	return c.apiKey != ""
}

func (c *credentials) currentToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.token
}

type loginRequest struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

type loginResponse struct {
	TokenId string `json:"tokenId"`
}

const loginEndpoint string = "login"

// NewClient returns a new Client configured with the specified application identity credentials (application key
// and secret) and base URI. baseUri should be the full URI to your DVLS instance (ex.: https://dvls.your-dvls-instance.com)
// The initial login has no deadline. To bound or cancel it, use NewClientWithContext instead.
func NewClient(appKey string, appSecret string, baseUri string) (Client, error) {
	return NewClientWithContext(context.Background(), appKey, appSecret, baseUri)
}

// NewClientWithContext returns a new Client configured with the specified application identity credentials
// (application key and secret) and base URI. baseUri should be the full URI to your DVLS instance
// (ex.: https://dvls.your-dvls-instance.com)
// The provided context applies to the initial login request only and can be used to cancel it
// or to give it a deadline. It is not retained by the returned Client.
func NewClientWithContext(ctx context.Context, appKey string, appSecret string, baseUri string) (Client, error) {
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: &credentials{appKey: appKey, appSecret: appSecret},
	}

	err := client.loginWithContext(ctx)
	if err != nil {
		return Client{}, fmt.Errorf("login failed \"%w\"", err)
	}

	client.initServices()

	return client, nil
}

// NewClientWithApiKey returns a new Client configured with the specified API key and base URI.
// The API key is sent directly on each request using the Authorization Bearer scheme, so no
// login exchange or token refresh is performed. baseUri should be the full URI to your DVLS
// instance (ex.: https://dvls.your-dvls-instance.com)
func NewClientWithApiKey(apiKey string, baseUri string) (Client, error) {
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: &credentials{apiKey: apiKey},
	}

	client.initServices()

	return client, nil
}

func (c *Client) initServices() {
	c.common.client = c

	c.Entries = &Entries{
		Certificate: (*EntryCertificateService)(&c.common),
		Credential:  (*EntryCredentialService)(&c.common),
		Folder:      (*EntryFolderService)(&c.common),
		Host:        (*EntryHostService)(&c.common),
		Website:     (*EntryWebsiteService)(&c.common),
		Permissions: (*EntryPermissionsService)(&c.common),
	}
	c.Vaults = (*Vaults)(&c.common)
	c.AdministrativeRoles = (*AdministrativeRoles)(&c.common)
	c.AdministrativeRoleAssignments = (*AdministrativeRoleAssignments)(&c.common)
	c.Users = (*Users)(&c.common)
	c.UserGroups = (*UserGroups)(&c.common)
}

func (c *Client) loginWithContext(ctx context.Context) error {
	c.credential.mu.Lock()
	defer c.credential.mu.Unlock()

	return c.loginLocked(ctx)
}

// reloginWithContext logs in again unless another request already replaced staleToken.
func (c *Client) reloginWithContext(ctx context.Context, staleToken string) error {
	c.credential.mu.Lock()
	defer c.credential.mu.Unlock()

	if c.credential.token != staleToken {
		return nil
	}

	return c.loginLocked(ctx)
}

func (c *Client) loginLocked(ctx context.Context) error {
	reqUrl, err := url.JoinPath(c.baseUri, apiV2Prefix, loginEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build login url: %w", err)
	}

	body, err := jsonBody(loginRequest{AppKey: c.credential.appKey, AppSecret: c.credential.appSecret})
	if err != nil {
		return err
	}

	resp, err := c.rawRequestWithContext(ctx, reqUrl, http.MethodPost, defaultContentType, "", body)
	if err != nil {
		return fmt.Errorf("error while submitting login request: %w", err)
	}

	var login loginResponse
	err = unmarshalBody(resp.Response, &login)
	if err != nil {
		return err
	}
	if login.TokenId == "" {
		return fmt.Errorf("login response has no tokenId")
	}

	c.credential.token = login.TokenId

	return nil
}
