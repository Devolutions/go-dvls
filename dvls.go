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

func getUntypedResultCode(respBody []byte) (int, error) {
	result := struct {
		Result int
	}{}
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to get result. error: %w", err)
	}
	return result.Result, nil
}

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
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId, "/sensitive-data")
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to make request. error: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
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

	result, err := getUntypedResultCode(body)
	if err != nil {
		return DvlsSecret{}, err
	} else if result != 1 {
		return DvlsSecret{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", result)
	}

	err = json.Unmarshal(body, &secret)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}
	secret.ID = entryId

	return secret, nil
}

func (c *Client) GetEntry(entryId string) (DvlsEntry, error) {
	islogged, err := c.isLogged()
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to fetch login status. error: %w", err)
	}
	if !islogged {
		err := c.refreshToken()
		if err != nil {
			return DvlsEntry{}, fmt.Errorf("failed to refresh login token. error: %w", err)
		}
	}

	var entry DvlsEntry
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to make request. error: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("tokenId", c.credential.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("error while submitting get secret request. error: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return DvlsEntry{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to read response body. error: %w", err)
	}

	result, err := getUntypedResultCode(body)
	if err != nil {
		return DvlsEntry{}, err
	} else if result != 1 {
		return DvlsEntry{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", result)
	}

	err = json.Unmarshal(body, &entry)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}
