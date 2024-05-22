package dvls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Vault represents a DVLS vault. Contains relevant vault information.
type Vault struct {
	ID            string
	Name          string
	Description   string
	SecurityLevel VaultSecurityLevel
	Visibility    VaultVisibility
	CreationDate  *ServerTime
	ModifiedDate  *ServerTime
	password      *string
}

type VaultOptions struct {
	Password *string
}

type rawVault struct {
	Description            string  `json:"description"`
	Id                     string  `json:"id"`
	IdString               string  `json:"idString"`
	Image                  string  `json:"image"`
	ImageBytes             string  `json:"imageBytes"`
	ImageName              string  `json:"imageName"`
	IsAllowedOffline       bool    `json:"isAllowedOffline"`
	IsLocked               bool    `json:"isLocked"`
	IsPrivate              bool    `json:"isPrivate"`
	Password               *string `json:"password,omitempty"`
	HasPasswordChanged     *bool   `json:"hasPasswordChanged,omitempty"`
	ModifiedLoggedUserName string  `json:"modifiedLoggedUserName"`
	ModifiedUserName       string  `json:"modifiedUserName"`
	Name                   string  `json:"name"`
	RepositorySettings     struct {
		QuickAddEntries             [0]struct{} `json:"quickAddEntries"`
		IsPasswordProtected         bool        `json:"isPasswordProtected"`
		MasterPasswordHash          *string     `json:"masterPasswordHash,omitempty"`
		VaultSecurityLevel          *int        `json:"vaultSecurityLevel,omitempty"`
		VaultAllowAccessRequestRole int         `json:"vaultAllowAccessRequestRole"`
		VaultType                   int         `json:"vaultType"`
	} `json:"repositorySettings"`
	Selected bool `json:"selected"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vault) UnmarshalJSON(b []byte) error {
	var raw struct {
		Data rawVault
	}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var securityLevel VaultSecurityLevel

	if raw.Data.RepositorySettings.VaultSecurityLevel != nil {
		securityLevel = VaultSecurityLevel(*raw.Data.RepositorySettings.VaultSecurityLevel)
	}

	vault := Vault{
		ID:            raw.Data.Id,
		Name:          raw.Data.Name,
		Description:   raw.Data.Description,
		SecurityLevel: securityLevel,
		Visibility:    VaultVisibility(raw.Data.RepositorySettings.VaultAllowAccessRequestRole),
	}

	*v = vault

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v Vault) MarshalJSON() ([]byte, error) {
	var raw rawVault

	securityLevel := 1

	if v.SecurityLevel == VaultSecurityLevelHigh {
		securityLevel = 0
		raw.RepositorySettings.VaultType = 1
	}

	if v.password != nil {
		raw.Password = v.password
		hasPasswordChanged := true
		raw.HasPasswordChanged = &hasPasswordChanged
	}

	raw.Name = v.Name
	raw.Description = v.Description
	raw.Id = v.ID
	raw.IdString = v.ID
	raw.RepositorySettings.VaultSecurityLevel = &securityLevel
	raw.RepositorySettings.VaultAllowAccessRequestRole = int(v.Visibility)

	if v.SecurityLevel == VaultSecurityLevelStandard {
		raw.IsAllowedOffline = true
	}

	json, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return json, nil
}

const (
	vaultEndpoint string = "/api/security/repositories"
)

// GetVault returns a single Vault based on vaultId.
func (c *Client) GetVault(vaultId string) (Vault, error) {
	var vault Vault
	reqUrl, err := url.JoinPath(c.baseUri, vaultEndpoint, vaultId)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodGet, nil)
	if err != nil {
		return Vault{}, fmt.Errorf("error while fetching vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Vault{}, err
	}

	err = json.Unmarshal(resp.Response, &vault)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to unmarshall response body. error: %w", err)
	}

	return vault, nil
}

// NewVault creates a new Vault based on vault.
func (c *Client) NewVault(vault Vault, options *VaultOptions) error {
	reqUrl, err := url.JoinPath(c.baseUri, vaultEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build vault url. error: %w", err)
	}

	vault.CreationDate = nil
	vault.ModifiedDate = nil

	if options != nil {
		vault.password = options.Password
	}

	vaultJson, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to marshall body. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPut, bytes.NewBuffer(vaultJson))
	if err != nil {
		return fmt.Errorf("error while creating vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// UpdateVault updates a Vault based on vault.
func (c *Client) UpdateVault(vault Vault, options *VaultOptions) error {
	_, err := c.GetVault(vault.ID)
	if err != nil {
		return fmt.Errorf("error while fetching vault. error: %w", err)
	}

	err = c.NewVault(vault, options)
	if err != nil {
		return fmt.Errorf("error while updating vault. error: %w", err)
	}

	return nil
}

// DeleteVault deletes a Vault based on vaultId.
func (c *Client) DeleteVault(vaultId string) error {
	reqUrl, err := url.JoinPath(c.baseUri, vaultEndpoint, vaultId)
	if err != nil {
		return fmt.Errorf("failed to delete vault url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// ValidateVaultPassword validates a Vault password based on vaultId and password.
func (c *Client) ValidateVaultPassword(vaultId string, password string) (bool, error) {
	reqUrl, err := url.JoinPath(c.baseUri, vaultEndpoint, vaultId, "login")
	if err != nil {
		return false, fmt.Errorf("failed to build vault url. error: %w", err)
	}

	resp, err := c.Request(reqUrl, http.MethodPost, bytes.NewBufferString(fmt.Sprintf("\"%s\"", password)))
	if err != nil {
		return false, fmt.Errorf("error while fetching vault. error: %w", err)
	} else if resp.Result == uint8(SaveResultAccessDenied) {
		return false, nil
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return false, err
	}

	return true, nil
}
