package dvls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

type Timezone struct {
	Id                         string
	DisplayName                string
	StandardName               string
	DaylightName               string
	BaseUtcOffset              string
	AdjustmentRules            []TimezoneAdjustmentRule
	SupportsDaylightSavingTime bool
}

type TimezoneAdjustmentRule struct {
	DateStart               ServerTime
	DateEnd                 ServerTime
	DaylightDelta           string
	DaylightTransitionStart TimezoneAdjustmentRuleTransitionTime
	DaylightTransitionEnd   TimezoneAdjustmentRuleTransitionTime
	BaseUtcOffsetDelta      string
	NoDaylightTransitions   bool
}

type TimezoneAdjustmentRuleTransitionTime struct {
	TimeOfDay       ServerTime
	Month           int
	Week            int
	Day             int
	DayOfWeek       int
	IsFixedDateRule bool
}

type ServerTime struct {
	time.Time
}

func (z *ServerTime) UnmarshalJSON(d []byte) error {
	s := strings.Trim(string(d), "\"")
	if s == "null" {
		return nil
	}

	dateParsed, err := time.Parse(serverTimeLayout, s)
	if err != nil {
		return err
	}

	z.Time = dateParsed
	return nil
}

const (
	serverInfoEndpoint      string = "/api/server-information"
	serverTimezonesEndpoint string = "/api/configuration/timezones"
	serverTimeLayout        string = "2006-01-02T15:04:05"
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

func (c *Client) GetServerTimezones() ([]Timezone, error) {
	var timezones []Timezone
	reqUrl, err := url.JoinPath(c.baseUri, serverTimezonesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build timezone info url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching timezones. error: %w", err)
	} else if resp.Result != 1 {
		return nil, fmt.Errorf("unexpected result code %d", resp.Result)
	}

	raw := struct {
		Data []Timezone
	}{}
	err = json.Unmarshal(resp.Response, &raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	timezones = raw.Data

	return timezones, nil
}
