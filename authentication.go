package dvls

import (
	"bytes"
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

	Entry   *Entry
	Entries *Entries
	Vaults  *Vaults
}

type service struct {
	client *Client
}

type credentials struct {
	appKey    string
	appSecret string
	token     string
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

	client.common.client = &client

	client.Entries = &Entries{
		Certificate: (*EntryCertificateService)(&client.common),
		Website:     (*EntryWebsiteService)(&client.common),
		Host:        (*EntryHostService)(&client.common),
	}
	client.Vaults = (*Vaults)(&client.common)

	return client, nil
}

func (c *Client) login() error {
	loginBody := fmt.Sprintf("AppKey=%s&AppSecret=%s", c.credential.appKey, c.credential.appSecret)

	reqUrl, err := url.JoinPath(c.baseUri, loginEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build login url. error: %w", err)
	}

	resp, err := c.rawRequest(reqUrl, http.MethodPost, loginContentType, bytes.NewBufferString(loginBody))
	if err != nil {
		return fmt.Errorf("error while submitting login request. error: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(resp.Response, &loginResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	c.credential.token = loginResponse.TokenId

	return nil
}

func (c *Client) isLogged() (bool, error) {
	reqUrl, err := url.JoinPath(c.baseUri, isLoggedEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to build isLogged url. error: %w", err)
	}

	resp, err := c.rawRequest(reqUrl, http.MethodGet, defaultContentType, nil)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool into Go value") {
		return false, fmt.Errorf("error while submitting isLogged request. error: %w", err)
	}

	if string(resp.Response) == "false" {
		return false, nil
	}

	return true, nil
}
