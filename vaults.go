package dvls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type VaultVisibility string

const (
	VaultVisibilityDefault VaultVisibility = "Default"
	VaultVisibilityPrivate VaultVisibility = "Never"
	VaultVisibilityPublic  VaultVisibility = "Everyone"
)

type VaultSecurityLevel string

const (
	VaultSecurityLevelStandard VaultSecurityLevel = "Standard"
	VaultSecurityLevelHigh     VaultSecurityLevel = "High"
)

type VaultContentType string

const (
	VaultContentTypeEverything          VaultContentType = "Everything"
	VaultContentTypeDefault             VaultContentType = "Default" // Equivalent to Everything, used by system vaults (Default, User vault)
	VaultContentTypeSecrets             VaultContentType = "Secrets"
	VaultContentTypeCredentials         VaultContentType = "Credentials"
	VaultContentTypeBusinessInformation VaultContentType = "BusinessInformation"
)

type Vaults service

// Vault represents a DVLS vault.
type Vault struct {
	Id            string             `json:"id,omitempty"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ContentType   VaultContentType   `json:"contentType"`
	Type          string             `json:"type,omitempty"`
	SecurityLevel VaultSecurityLevel `json:"securityLevel"`
	Visibility    VaultVisibility    `json:"visibility"`
}

// vaultListResponse represents the paginated response from the vault list endpoint.
type vaultListResponse struct {
	Data        []Vault `json:"data"`
	CurrentPage int     `json:"currentPage"`
	PageSize    int     `json:"pageSize"`
	TotalCount  int     `json:"totalCount"`
	TotalPage   int     `json:"totalPage"`
}

// vaultRequest represents the request body for create/update operations.
type vaultRequest struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ContentType   VaultContentType   `json:"contentType"`
	SecurityLevel VaultSecurityLevel `json:"securityLevel"`
	Visibility    VaultVisibility    `json:"visibility"`
}

const (
	vaultEndpoint string = "/api/v1/vault"
)

var ErrVaultNotFound = fmt.Errorf("vault not found")
var ErrMultipleVaultsFound = fmt.Errorf("multiple vaults found")

// List returns all vaults.
func (c *Vaults) List() ([]Vault, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all vaults.
// The provided context can be used to cancel the request.
func (c *Vaults) ListWithContext(ctx context.Context) ([]Vault, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to build vault url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching vaults: %w", err)
	}

	var listResp vaultListResponse
	err = json.Unmarshal(resp.Response, &listResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return listResp.Data, nil
}

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
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	if err != nil {
		return Vault{}, fmt.Errorf("error while fetching vault: %w", err)
	}

	err = json.Unmarshal(resp.Response, &vault)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return vault, nil
}

// GetByName returns a single Vault based on name.
// Returns ErrVaultNotFound if no vault is found.
// Returns ErrMultipleVaultsFound if more than one vault matches the name.
func (c *Vaults) GetByName(name string) (Vault, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single Vault based on name.
// Returns ErrVaultNotFound if no vault is found.
// Returns ErrMultipleVaultsFound if more than one vault matches the name.
// The provided context can be used to cancel the request.
func (c *Vaults) GetByNameWithContext(ctx context.Context, name string) (Vault, error) {
	vaults, err := c.ListWithContext(ctx)
	if err != nil {
		return Vault{}, err
	}

	var matches []Vault
	for _, v := range vaults {
		if v.Name == name {
			matches = append(matches, v)
		}
	}

	if len(matches) == 0 {
		return Vault{}, ErrVaultNotFound
	}

	if len(matches) > 1 {
		return Vault{}, ErrMultipleVaultsFound
	}

	return matches[0], nil
}

// New creates a new Vault and returns the created vault.
func (c *Vaults) New(vault Vault) (Vault, error) {
	return c.NewWithContext(context.Background(), vault)
}

// NewWithContext creates a new Vault and returns the created vault.
// The provided context can be used to cancel the request.
func (c *Vaults) NewWithContext(ctx context.Context, vault Vault) (Vault, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	// Convert Default to Everything (API rejects "Default" for creation)
	contentType := vault.ContentType
	if contentType == VaultContentTypeDefault {
		contentType = VaultContentTypeEverything
	}

	reqBody := vaultRequest{
		Name:          vault.Name,
		Description:   vault.Description,
		ContentType:   contentType,
		SecurityLevel: vault.SecurityLevel,
		Visibility:    vault.Visibility,
	}

	vaultJson, err := json.Marshal(reqBody)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPost, bytes.NewBuffer(vaultJson))
	if err != nil {
		return Vault{}, fmt.Errorf("error while creating vault: %w", err)
	}

	var createdVault Vault
	err = json.Unmarshal(resp.Response, &createdVault)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return createdVault, nil
}

// Update updates an existing Vault and returns the updated vault.
func (c *Vaults) Update(vault Vault) (Vault, error) {
	return c.UpdateWithContext(context.Background(), vault)
}

// UpdateWithContext updates an existing Vault and returns the updated vault.
// The provided context can be used to cancel the request.
func (c *Vaults) UpdateWithContext(ctx context.Context, vault Vault) (Vault, error) {
	reqUrl, err := url.JoinPath(c.client.baseUri, vaultEndpoint, vault.Id)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	reqBody := vaultRequest{
		Name:          vault.Name,
		Description:   vault.Description,
		ContentType:   vault.ContentType,
		SecurityLevel: vault.SecurityLevel,
		Visibility:    vault.Visibility,
	}

	vaultJson, err := json.Marshal(reqBody)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	resp, err := c.client.RequestWithContext(ctx, reqUrl, http.MethodPut, bytes.NewBuffer(vaultJson))
	if err != nil {
		return Vault{}, fmt.Errorf("error while updating vault: %w", err)
	}

	var updatedVault Vault
	err = json.Unmarshal(resp.Response, &updatedVault)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return updatedVault, nil
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
		return fmt.Errorf("failed to build vault url: %w", err)
	}

	_, err = c.client.RequestWithContext(ctx, reqUrl, http.MethodDelete, nil)
	if err != nil {
		return fmt.Errorf("error while deleting vault: %w", err)
	}

	return nil
}
