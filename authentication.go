package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	client     *http.Client
	baseUri    string
	credential credentials

	ClientUser User
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

type User struct {
	ID       string
	Username string
	UserType UserAuthenticationType
}

func (u *User) UnmarshalJSON(d []byte) error {
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

	return nil
}

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

	user, err := client.login()
	if err != nil {
		return Client{}, fmt.Errorf("login failed \"%w\"", err)
	}

	client.ClientUser = user

	return client, nil
}

func (c *Client) login() (User, error) {
	loginBody := loginReqBody{
		Username: c.credential.username,
		LoginParameters: loginParameters{
			Password: c.credential.password,
			Client:   "Cli",
		},
	}
	loginJson, err := json.Marshal(loginBody)
	if err != nil {
		return User{}, fmt.Errorf("failed to marshall login body. error: %w", err)
	}

	reqUrl, err := url.JoinPath(c.baseUri, loginEndpoint)
	if err != nil {
		return User{}, fmt.Errorf("failed to build login url. error: %w", err)
	}

	resp, err := c.rawRequest(reqUrl, http.MethodPost, bytes.NewBuffer(loginJson))
	if err != nil {
		return User{}, fmt.Errorf("error while submitting refreshtoken request. error: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(resp.Response, &loginResponse)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}
	if loginResponse.Data.Result != ServerLoginSuccess {
		return User{}, fmt.Errorf("failed to refresh token (%s) : %s", loginResponse.Data.Result, loginResponse.Data.Message)
	}

	var user User
	err = json.Unmarshal(resp.Response, &user)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshall user body. error: %w", err)
	}

	c.credential.token = loginResponse.Data.TokenId

	return user, nil
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

	resp, err := c.rawRequest(reqUrl, http.MethodPost, bytes.NewBuffer(loginJson))
	if err != nil {
		return fmt.Errorf("error while submitting refreshtoken request. error: %w", err)
	}

	var loginResponse loginResponse
	err = json.Unmarshal(resp.Response, &loginResponse)
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
		return false, fmt.Errorf("failed to build isLogged url. error: %w", err)
	}

	resp, err := c.rawRequest(reqUrl, http.MethodGet, nil)
	if err != nil && !strings.Contains(err.Error(), "json: cannot unmarshal bool into Go value") {
		return false, fmt.Errorf("error while submitting isLogged request. error: %w", err)
	}

	if string(resp.Response) == "false" {
		return false, nil
	}

	return true, nil
}
