package dvls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Server struct {
	AccessUri     string
	TimeZone      string
	ServerName    string `json:"servername"`
	Version       string
	SystemMessage string
}

func (s *Server) UnmarshalJSON(d []byte) error {
	raw := struct {
		Data struct {
			AccessUri          string
			SelectedTimeZoneId string
			ServerName         string
			Version            string
			SystemMessage      string
		}
	}{}
	err := json.Unmarshal(d, &raw)
	if err != nil {
		return err
	}

	s.TimeZone = raw.Data.SelectedTimeZoneId
	s.AccessUri = raw.Data.AccessUri
	s.ServerName = raw.Data.ServerName
	s.Version = raw.Data.Version
	s.SystemMessage = raw.Data.SystemMessage

	return nil
}

const (
	serverInfoEndpoint string = "/api/server-information"
)

func (c *Client) GetServerInfo() (Server, error) {
	var server Server
	reqUrl, err := url.JoinPath(c.baseUri, serverInfoEndpoint)
	if err != nil {
		return Server{}, fmt.Errorf("failed to build server info url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return Server{}, fmt.Errorf("error while fetching server info. error: %w", err)
	} else if resp.Result != 1 {
		return Server{}, fmt.Errorf("unexpected result code %d", resp.Result)
	}

	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return Server{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return server, nil
}
