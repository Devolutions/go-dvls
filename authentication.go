package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Client represents the DVLS client used to communicate with the API.
type Client struct {
	client     *http.Client
	baseUri    string
	credential *credentials // a pointer, so copies of a Client and their services share one session

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

type credentials struct {
	appKey    string
	appSecret string
	apiKey    string
	token     string
	mu        sync.Mutex
}

func (c *credentials) useApiKey() bool {
	return c.apiKey != ""
}

func (c *credentials) currentToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.token
}

type loginResponse struct {
	TokenId string
}

const (
	loginEndpoint    string = "/api/v1/login"
	isLoggedEndpoint string = "/api/is-logged"
)

const loginContentType = "application/x-www-form-urlencoded"

// NewClient returns a new Client configured with the specified credentials and
// base URI. baseUri should be the full URI to your DVLS instance (ex.: https://dvls.your-dvls-instance.com)
// The initial login has no deadline. To bound or cancel it, use NewClientWithContext instead.
func NewClient(appKey string, appSecret string, baseUri string) (Client, error) {
	return NewClientWithContext(context.Background(), appKey, appSecret, baseUri)
}

// NewClientWithContext returns a new Client configured with the specified credentials and
// base URI. baseUri should be the full URI to your DVLS instance (ex.: https://dvls.your-dvls-instance.com)
// The provided context applies to the initial login request only and can be used to cancel it
// or to give it a deadline. It is not retained by the returned Client.
func NewClientWithContext(ctx context.Context, appKey string, appSecret string, baseUri string) (Client, error) {
	credential := &credentials{appKey: appKey, appSecret: appSecret}
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: credential,
	}

	token, err := client.loginWithContext(ctx)
	if err != nil {
		return Client{}, fmt.Errorf("login failed \"%w\"", err)
	}
	credential.token = token

	client.initServices()

	return client, nil
}

// NewClientWithApiKey returns a new Client configured with the specified API key and base URI.
// The API key is sent directly on each request using the Authorization Bearer scheme, so no
// login exchange or token refresh is performed. baseUri should be the full URI to your DVLS
// instance (ex.: https://dvls.your-dvls-instance.com)
func NewClientWithApiKey(apiKey string, baseUri string) (Client, error) {
	credential := &credentials{apiKey: apiKey}
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: credential,
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

func (c *Client) loginWithContext(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("AppKey", c.credential.appKey)
	form.Set("AppSecret", c.credential.appSecret)
	loginBody := form.Encode()

	reqUrl, err := url.JoinPath(c.baseUri, loginEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to build login url: %w", err)
	}

	resp, err := c.rawRequestWithContext(ctx, reqUrl, http.MethodPost, loginContentType, "", bytes.NewBufferString(loginBody))
	if err != nil {
		return "", fmt.Errorf("error while submitting login request: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(resp.Response, &loginResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return loginResponse.TokenId, nil
}

func (c *Client) renewSession(ctx context.Context, staleToken string) (bool, error) {
	c.credential.mu.Lock()
	defer c.credential.mu.Unlock()

	if c.credential.token != staleToken {
		return true, nil
	}

	islogged, err := c.isLoggedWithContext(ctx, staleToken)
	if err != nil {
		return false, fmt.Errorf("failed to fetch login status: %w", err)
	}
	if islogged {
		return false, nil
	}

	token, err := c.loginWithContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to refresh login token: %w", err)
	}
	c.credential.token = token

	return true, nil
}

func (c *Client) isLogged() (bool, error) {
	return c.isLoggedWithContext(context.Background(), c.credential.currentToken())
}

func (c *Client) isLoggedWithContext(ctx context.Context, token string) (bool, error) {
	reqUrl, err := url.JoinPath(c.baseUri, isLoggedEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to build isLogged url: %w", err)
	}

	resp, err := c.rawRequestWithContext(ctx, reqUrl, http.MethodGet, defaultContentType, token, nil)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool into Go value") {
		return false, fmt.Errorf("error while submitting isLogged request: %w", err)
	}

	if string(resp.Response) == "false" {
		return false, nil
	}

	return true, nil
}
