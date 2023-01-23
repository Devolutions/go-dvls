package dvls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	entryEndpoint string = "/api/connections/partial"
)

func (c *Client) GetSecret(entryId string) (DvlsSecret, error) {
	islogged, err := c.isLogged()
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to fetch login status. error: %w", err)
	}
	if !islogged {
		err := c.refreshToken()
		if err != nil {
			return DvlsSecret{}, fmt.Errorf("failed to refresh login token. error: %w", err)
		}
	}

	var secret DvlsSecret
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to make request. error: %w", err)
	}

	req.Header.Add("tokenId", c.credential.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("error while submitting get secret request. error: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return DvlsSecret{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to read response body. error: %w", err)
	}

	err = json.Unmarshal(body, &secret)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return secret, nil
}
