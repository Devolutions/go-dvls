package dvls

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type EntryWebsiteService service

const entryWebsiteType string = "WebBrowser"

// EntryWebsite represents a website entry in DVLS. ConnectionSubType is the browser that opens it.
type EntryWebsite struct {
	Id                string
	VaultId           string
	EntryName         string
	Description       string
	EntryFolderPath   string
	ModifiedDate      *ServerTime
	ConnectionSubType ServerConnectionSubType
	Tags              []string

	WebsiteDetails EntryWebsiteAuthDetails
}

// EntryWebsiteAuthDetails represents website-specific fields
type EntryWebsiteAuthDetails struct {
	Username string
	Password *string
	URL      string
}

type entryWebsiteData struct {
	Url      string  `json:"url"`
	Username string  `json:"username"`
	Password *string `json:"password"`
}

func websiteFromRaw(raw rawEntry) (EntryWebsite, error) {
	if raw.Type != entryWebsiteType {
		return EntryWebsite{}, ErrUnsupportedEntryType{Type: raw.Type, SubType: raw.SubType}
	}

	var data entryWebsiteData
	if len(raw.Data) > 0 {
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return EntryWebsite{}, fmt.Errorf("failed to unmarshal website details: %w", err)
		}
	}

	return EntryWebsite{
		Id:                raw.Id,
		VaultId:           raw.VaultId,
		EntryName:         raw.Name,
		Description:       raw.Description,
		EntryFolderPath:   raw.Path,
		ModifiedDate:      raw.ModifiedOn,
		ConnectionSubType: ServerConnectionSubType(raw.SubType),
		Tags:              raw.Tags,
		WebsiteDetails:    EntryWebsiteAuthDetails{Username: data.Username, Password: data.Password, URL: data.Url},
	}, nil
}

// Get returns a single website entry, without its password, specified by vaultId and entryId. Call
// GetWebsiteDetails with the returned entry to fetch the password.
func (s *EntryWebsiteService) Get(vaultId string, entryId string) (EntryWebsite, error) {
	return s.GetWithContext(context.Background(), vaultId, entryId)
}

// GetWithContext returns a single website entry, without its password, specified by vaultId and entryId. Call
// GetWebsiteDetailsWithContext with the returned entry to fetch the password.
// The provided context can be used to cancel the request.
func (s *EntryWebsiteService) GetWithContext(ctx context.Context, vaultId string, entryId string) (EntryWebsite, error) {
	raw, err := s.client.getRawEntry(ctx, vaultId, entryId, url.Values{"includePasswords": {"false"}})
	if err != nil {
		return EntryWebsite{}, err
	}

	return websiteFromRaw(raw)
}

// GetWebsiteDetails returns entry with the entry.WebsiteDetails.Password field.
func (c *EntryWebsiteService) GetWebsiteDetails(entry EntryWebsite) (EntryWebsite, error) {
	return c.GetWebsiteDetailsWithContext(context.Background(), entry)
}

// GetWebsiteDetailsWithContext returns entry with the entry.WebsiteDetails.Password field.
// The provided context can be used to cancel the request.
func (c *EntryWebsiteService) GetWebsiteDetailsWithContext(ctx context.Context, entry EntryWebsite) (EntryWebsite, error) {
	raw, err := c.client.getRawEntry(ctx, entry.VaultId, entry.Id, nil)
	if err != nil {
		return EntryWebsite{}, err
	}

	stored, err := websiteFromRaw(raw)
	if err != nil {
		return EntryWebsite{}, err
	}

	entry.WebsiteDetails.Password = stored.WebsiteDetails.Password

	return entry, nil
}
