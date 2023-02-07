package dvls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type DvlsEntry struct {
	ID                string
	EntryName         string
	Username          string
	ModifiedDate      time.Time
	ConnectionType    ServerConnectionType
	ConnectionSubType ServerConnectionSubType
}

func (e *DvlsEntry) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			ID                string
			Name              string
			Username          string
			ModifiedDate      string
			ConnectionType    ServerConnectionType
			ConnectionSubType ServerConnectionSubType
		}
		Result int
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02T15:04:05", raw.Data.ModifiedDate)
	if err != nil {
		return err
	}

	e.ID = raw.Data.ID
	e.EntryName = raw.Data.Name
	e.ConnectionType = raw.Data.ConnectionType
	e.ConnectionSubType = raw.Data.ConnectionSubType
	e.ModifiedDate = date
	e.Username = raw.Data.Username

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

const (
	entryEndpoint string = "/api/connections/partial"
)

func (c *Client) GetSecret(entryId string) (DvlsSecret, error) {
	var secret DvlsSecret
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId, "/sensitive-data")
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, nil)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("error while fetching sensitive data. error: %w", err)
	} else if resp.Result != 1 {
		return DvlsSecret{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", resp.Result)
	}

	err = json.Unmarshal(resp.Response, &secret)
	if err != nil {
		return DvlsSecret{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}
	secret.ID = entryId

	return secret, nil
}

func (c *Client) GetEntry(entryId string) (DvlsEntry, error) {
	var entry DvlsEntry
	reqUrl, err := url.JoinPath(c.baseUri, entryEndpoint, entryId)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to build entry url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("error while fetching entry. error: %w", err)
	} else if resp.Result != 1 {
		return DvlsEntry{}, fmt.Errorf("unexpected result code %d. Make sure the entry ID is correct and the user has access to the entry", resp.Result)
	}

	err = json.Unmarshal(resp.Response, &entry)
	if err != nil {
		return DvlsEntry{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return entry, nil
}
