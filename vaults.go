package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Vaults service

// Vault represents a DVLS vault. Contains relevant vault information.
type Vault struct {
	Id            string
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
		Id:            raw.Data.Id,
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
	raw.Id = v.Id
	raw.IdString = v.Id
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

// Get returns a single Vault based on vaultId.
func (c *Vaults) Get(vaultId string) (Vault, error) {
	return c.GetWithContext(context.Background(), vaultId)
}

// GetWithContext returns a single Vault based on vaultId.
// The provided context can be used to cancel the request.
func (c *Vaults) GetWithContext(ctx context.Context, vaultId string) (Vault, error) {
	var vault Vault
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint, vaultId)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url. error: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return Vault{}, fmt.Errorf("error while fetching vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return Vault{}, err
	}

	err = json.Unmarshal(resp.Response, &vault)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to unmarshal response body. error: %w", err)
	}

	return vault, nil
}

// New creates a new Vault based on vault.
func (c *Vaults) New(vault Vault, options *VaultOptions) error {
	return c.NewWithContext(context.Background(), vault, options)
}

// NewWithContext creates a new Vault based on vault.
// The provided context can be used to cancel the request.
func (c *Vaults) NewWithContext(ctx context.Context, vault Vault, options *VaultOptions) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint)
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
		return fmt.Errorf("failed to marshal body. error: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPut, bytes.NewBuffer(vaultJson))
	if err != nil {
		return fmt.Errorf("error while creating vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// Update updates a Vault based on vault.
func (c *Vaults) Update(vault Vault, options *VaultOptions) error {
	return c.UpdateWithContext(context.Background(), vault, options)
}

// UpdateWithContext updates a Vault based on vault.
// The provided context can be used to cancel the request.
func (c *Vaults) UpdateWithContext(ctx context.Context, vault Vault, options *VaultOptions) error {
	_, err := c.client.Vaults.GetWithContext(ctx, vault.Id)
	if err != nil {
		return fmt.Errorf("error while fetching vault. error: %w", err)
	}

	err = c.client.Vaults.NewWithContext(ctx, vault, options)
	if err != nil {
		return fmt.Errorf("error while updating vault. error: %w", err)
	}

	return nil
}

// Delete deletes a Vault based on vaultId.
func (c *Vaults) Delete(vaultId string) error {
	return c.DeleteWithContext(context.Background(), vaultId)
}

// DeleteWithContext deletes a Vault based on vaultId.
// The provided context can be used to cancel the request.
func (c *Vaults) DeleteWithContext(ctx context.Context, vaultId string) error {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint, vaultId)
	if err != nil {
		return fmt.Errorf("failed to delete vault url. error: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting vault. error: %w", err)
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return err
	}

	return nil
}

// ValidatePassword validates a Vault password based on vaultId and password.
func (c *Vaults) ValidatePassword(vaultId string, password string) (bool, error) {
	return c.ValidatePasswordWithContext(context.Background(), vaultId, password)
}

// ValidatePasswordWithContext validates a Vault password based on vaultId and password.
// The provided context can be used to cancel the request.
func (c *Vaults) ValidatePasswordWithContext(ctx context.Context, vaultId string, password string) (bool, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint, vaultId, "login")
	if err != nil {
		return false, fmt.Errorf("failed to build vault url. error: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBufferString(fmt.Sprintf("\"%s\"", password)))
	if err != nil {
		return false, fmt.Errorf("error while fetching vault. error: %w", err)
	} else if resp.Result == uint8(SaveResultAccessDenied) {
		return false, nil
	} else if err = resp.CheckRespSaveResult(); err != nil {
		return false, err
	}

	return true, nil
}
