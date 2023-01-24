package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	loginEndpoint    string = "/api/login/partial"
	isLoggedEndpoint string = "/api/is-logged"
)

func NewClient(username string, password string, baseUri string) (Client, error) {
	credential := credentials{username: username, password: password}
	client := Client{
		client:     &http.Client{},
		baseUri:    baseUri,
		credential: credential,
	}

	err := client.refreshToken()
	if err != nil {
		return Client{}, fmt.Errorf("login failed \"%w\"", err)
	}

	return client, nil
}

func (c *Client) refreshToken() error {
	loginBody := loginReqBody{
		Username: c.credential.username,
		LoginParameters: loginParameters{
			Password: c.credential.password,
			Client:   "Cli",
		},
	}
	loginJson, err := json.Marshal(loginBody)
	if err != nil {
		return fmt.Errorf("failed to marshall login body. error: %w", err)
	}

	reqUrl, err := url.JoinPath(c.baseUri, loginEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build login url. error: %w", err)
	}

	resp, err := c.client.Post(reqUrl, "application/json", bytes.NewBuffer(loginJson))
	if err != nil {
		return fmt.Errorf("error while submitting login request. error: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error while submitting login request. Unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body. error: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}
	if loginResponse.Data.Result != ServerLoginSuccess {
		return fmt.Errorf("failed to refresh token (%s) : %s", loginResponse.Data.Result, loginResponse.Data.Message)
	}

	c.credential.token = loginResponse.Data.TokenId

	return nil
}

func (c *Client) isLogged() (bool, error) {
	reqUrl, err := url.JoinPath(c.baseUri, isLoggedEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to isLogged url. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return false, fmt.Errorf("failed to make request. error: %w", err)
	}

	req.Header.Add("tokenId", c.credential.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error while submitting isLogged request. error: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body. error: %w", err)
	}

	if string(body) == "false" {
		return false, nil
	}

	return true, nil
}
