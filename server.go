package dvls

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Server represents the available server instance information.
type Server struct {
	AccessUri     string
	TimeZone      string
	ServerName    string `json:"servername"`
	Version       string
	SystemMessage string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

// Timezone represents a Server timezone.
type Timezone struct {
	Id                         string
	DisplayName                string
	StandardName               string
	DaylightName               string
	BaseUtcOffset              string
	AdjustmentRules            []TimezoneAdjustmentRule
	SupportsDaylightSavingTime bool
}

// TimezoneAdjustmentRule represents a Timezone Adjustment Rule.
type TimezoneAdjustmentRule struct {
	DateStart               ServerTime
	DateEnd                 ServerTime
	DaylightDelta           string
	DaylightTransitionStart TimezoneAdjustmentRuleTransitionTime
	DaylightTransitionEnd   TimezoneAdjustmentRuleTransitionTime
	BaseUtcOffsetDelta      string
	NoDaylightTransitions   bool
}

// TimezoneAdjustmentRuleTransitionTime represents a Timezone Adjustment Rule Transition Time.
type TimezoneAdjustmentRuleTransitionTime struct {
	TimeOfDay       ServerTime
	Month           int
	Week            int
	Day             int
	DayOfWeek       int
	IsFixedDateRule bool
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

const (
	serverPublicInfoEndpoint  string = "api/public-instance-information"
	serverPrivateInfoEndpoint string = "api/private-instance-information"
	serverTimezonesEndpoint   string = "/api/configuration/timezones"
)

var serverTimeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.9999999Z07:00",
	"2006-01-02T15:04:05.9999999Z",
	"2006-01-02T15:04:05",
}

// GetPublicServerInfo returns Server that contains public information on the DVLS instance.
func (c *Client) GetPublicServerInfo() (Server, error) {
	return c.GetPublicServerInfoWithContext(context.Background())
}

// GetPublicServerInfoWithContext returns Server that contains public information on the DVLS instance.
// The provided context can be used to cancel the request.
func (c *Client) GetPublicServerInfoWithContext(ctx context.Context) (Server, error) {
	var server Server
	reqUrl, err := url.JoinPath(c.baseUri, serverPublicInfoEndpoint)
	if err != nil {
		return Server{}, fmt.Errorf("failed to build server info url. error: %w", err)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return Server{}, fmt.Errorf("error while fetching server info. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Server{}, err
	}

	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return Server{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return server, nil
}

// GetPrivateServerInfo returns Server that contains private information on the DVLS instance (need authentication).
func (c *Client) GetPrivateServerInfo() (Server, error) {
	return c.GetPrivateServerInfoWithContext(context.Background())
}

// GetPrivateServerInfoWithContext returns Server that contains private information on the DVLS instance (need authentication).
// The provided context can be used to cancel the request.
func (c *Client) GetPrivateServerInfoWithContext(ctx context.Context) (Server, error) {
	var server Server
	reqUrl, err := url.JoinPath(c.baseUri, serverPrivateInfoEndpoint)
	if err != nil {
		return Server{}, fmt.Errorf("failed to build server info url. error: %w", err)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return Server{}, fmt.Errorf("error while fetching server info. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Server{}, err
	}

	err = json.Unmarshal(resp.Response, &server)
	if err != nil {
		return Server{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return server, nil
}

// GetServerTimezones returns an array of Timezone that contains all of the available timezones on
// the DVLS instance.
func (c *Client) GetServerTimezones() ([]Timezone, error) {
	return c.GetServerTimezonesWithContext(context.Background())
}

// GetServerTimezonesWithContext returns an array of Timezone that contains all of the available timezones on
// the DVLS instance.
// The provided context can be used to cancel the request.
func (c *Client) GetServerTimezonesWithContext(ctx context.Context) ([]Timezone, error) {
	var timezones []Timezone
	reqUrl, err := url.JoinPath(c.baseUri, serverTimezonesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build timezone info url. error: %w", err)
	}

	resp, err := c.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching timezones. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return nil, err
	}

	raw := struct {
		Data []Timezone
	}{}
	err = json.Unmarshal(resp.Response, &raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	timezones = raw.Data

	return timezones, nil
}
