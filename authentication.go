package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client represents the DVLS client used to communicate with the API.
type Client struct {
	client     *http.Client
	baseUri    string
	credential credentials

	common service

	Entries *Entries
	Vaults  *Vaults
}

type service struct {
	client *Client
}

type credentials struct {
	appKey    string
	appSecret string
	apiKey    string
	token     string
}

func (c credentials) useApiKey() bool {
	return c.apiKey != ""
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
func NewClient(appKey string, appSecret string, baseUri string) (Client, error) {
	credential := credentials{appKey: appKey, appSecret: appSecret}
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: credential,
	}

	err := client.login()
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
	credential := credentials{apiKey: apiKey}
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
	}
	c.Vaults = (*Vaults)(&c.common)
}

func (c *Client) login() error {
	return c.loginWithContext(context.Background())
}

func (c *Client) loginWithContext(ctx context.Context) error {
	form := url.Values{}
	form.Set("AppKey", c.credential.appKey)
	form.Set("AppSecret", c.credential.appSecret)
	loginBody := form.Encode()

	reqUrl, err := url.JoinPath(c.baseUri, loginEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build login url: %w", err)
	}

	resp, err := c.rawRequestWithContext(ctx, reqUrl, http.MethodPost, loginContentType, bytes.NewBufferString(loginBody))
	if err != nil {
		return fmt.Errorf("error while submitting login request: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(resp.Response, &loginResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	c.credential.token = loginResponse.TokenId

	return nil
}

func (c *Client) isLogged() (bool, error) {
	return c.isLoggedWithContext(context.Background())
}

func (c *Client) isLoggedWithContext(ctx context.Context) (bool, error) {
	reqUrl, err := url.JoinPath(c.baseUri, isLoggedEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to build isLogged url: %w", err)
	}

	resp, err := c.rawRequestWithContext(ctx, reqUrl, http.MethodGet, defaultContentType, nil)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool into Go value") {
		return false, fmt.Errorf("error while submitting isLogged request: %w", err)
	}

	if string(resp.Response) == "false" {
		return false, nil
	}

	return true, nil
}
