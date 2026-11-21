package dvls

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type EntryHostService service

const entryHostType string = "Host"

// EntryHost represents a host entry in DVLS
type EntryHost struct {
	Id              string
	VaultId         string
	EntryName       string
	Description     string
	EntryFolderPath string
	ModifiedDate    *ServerTime
	Tags            []string

	HostDetails EntryHostAuthDetails
}

// EntryHostAuthDetails represents host-specific fields
type EntryHostAuthDetails struct {
	Username string
	Password *string
	Host     string
}

type entryHostData struct {
	Host     string  `json:"host"`
	Username string  `json:"username"`
	Password *string `json:"password"`
}

func hostFromRaw(raw rawEntry) (EntryHost, error) {
	if raw.Type != entryHostType {
		return EntryHost{}, ErrUnsupportedEntryType{Type: raw.Type, SubType: raw.SubType}
	}

	var data entryHostData
	if len(raw.Data) > 0 {
		err := json.Unmarshal(raw.Data, &data)
		if err != nil {
			return EntryHost{}, fmt.Errorf("failed to unmarshal host details: %w", err)
		}
	}

	return EntryHost{
		Id:              raw.Id,
		VaultId:         raw.VaultId,
		EntryName:       raw.Name,
		Description:     raw.Description,
		EntryFolderPath: raw.Path,
		ModifiedDate:    raw.ModifiedOn,
		Tags:            raw.Tags,
		HostDetails:     EntryHostAuthDetails{Username: data.Username, Password: data.Password, Host: data.Host},
	}, nil
}

// Get returns a single host entry, without its password, specified by vaultId and entryId. Call GetHostDetails with
// the returned entry to fetch the password.
func (s *EntryHostService) Get(vaultId string, entryId string) (EntryHost, error) {
	return s.GetWithContext(context.Background(), vaultId, entryId)
}

// GetWithContext returns a single host entry, without its password, specified by vaultId and entryId. Call
// GetHostDetailsWithContext with the returned entry to fetch the password.
// The provided context can be used to cancel the request.
func (s *EntryHostService) GetWithContext(ctx context.Context, vaultId string, entryId string) (EntryHost, error) {
	raw, err := s.client.getRawEntry(ctx, vaultId, entryId, url.Values{"includePasswords": {"false"}})
	if err != nil {
		return EntryHost{}, err
	}

	return hostFromRaw(raw)
}

// GetHostDetails returns entry with the entry.HostDetails.Password field.
func (s *EntryHostService) GetHostDetails(entry EntryHost) (EntryHost, error) {
	return s.GetHostDetailsWithContext(context.Background(), entry)
}

// GetHostDetailsWithContext returns entry with the entry.HostDetails.Password field.
// The provided context can be used to cancel the request.
func (s *EntryHostService) GetHostDetailsWithContext(ctx context.Context, entry EntryHost) (EntryHost, error) {
	raw, err := s.client.getRawEntry(ctx, entry.VaultId, entry.Id, nil)
	if err != nil {
		return EntryHost{}, err
	}

	stored, err := hostFromRaw(raw)
	if err != nil {
		return EntryHost{}, err
	}

	entry.HostDetails.Password = stored.HostDetails.Password

	return entry, nil
}
