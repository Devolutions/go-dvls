package dvls

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ServerInfo represents the information the DVLS instance returns about itself.
type ServerInfo struct {
	Version string `json:"version"`
}

// ServerTime represents a time.Time that parses the correct server time layout.
type ServerTime struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (z *ServerTime) UnmarshalJSON(d []byte) error {
	s := strings.Trim(string(d), "\"")
	if s == "null" {
		return nil
	}

	for _, layout := range serverTimeLayouts {
		if dateParsed, err := time.Parse(layout, s); err == nil {
			z.Time = dateParsed
			return nil
		}
	}

	return fmt.Errorf("cannot parse server time %q", s)
}

const serverInfoEndpoint string = "server/info"

var serverTimeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.9999999Z07:00",
	"2006-01-02T15:04:05.9999999Z",
	"2006-01-02T15:04:05",
}

// GetServerInfo returns the information of the DVLS instance (its version).
func (c *Client) GetServerInfo() (ServerInfo, error) {
	return c.GetServerInfoWithContext(context.Background())
}

// GetServerInfoWithContext returns the information of the DVLS instance (its version).
// The provided context can be used to cancel the request.
func (c *Client) GetServerInfoWithContext(ctx context.Context) (ServerInfo, error) {
	reqUrl, err := c.v2Url(nil, serverInfoEndpoint)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("failed to build server info url: %w", err)
	}

	var info ServerInfo
	err = c.getJSON(ctx, reqUrl, &info)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("error while fetching server info: %w", err)
	}

	return info, nil
}
