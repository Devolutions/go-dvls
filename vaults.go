package dvls

import (
	"context"
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

// vaultCreateRequest leaves out the unset settings, which take their server default.
type vaultCreateRequest struct {
	Name          string             `json:"name"`
	Description   string             `json:"description,omitempty"`
	ContentType   VaultContentType   `json:"contentType,omitempty"`
	SecurityLevel VaultSecurityLevel `json:"securityLevel,omitempty"`
	Visibility    VaultVisibility    `json:"visibility,omitempty"`
}

type vaultUpdateRequest struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ContentType   VaultContentType   `json:"contentType"`
	SecurityLevel VaultSecurityLevel `json:"securityLevel"`
	Visibility    VaultVisibility    `json:"visibility"`
}

const vaultEndpoint string = "vault"

var ErrVaultNotFound = fmt.Errorf("vault not found")
var ErrMultipleVaultsFound = fmt.Errorf("multiple vaults found")

// List returns all vaults.
func (c *Vaults) List() ([]Vault, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all vaults.
// This function handles pagination automatically and returns all vaults across all pages.
// The provided context can be used to cancel the request.
func (c *Vaults) ListWithContext(ctx context.Context) ([]Vault, error) {
	return c.list(ctx, nil)
}

func (c *Vaults) list(ctx context.Context, query url.Values) ([]Vault, error) {
	vaults, err := listAll[Vault](ctx, c.client, query, vaultEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching vaults: %w", err)
	}

	return vaults, nil
}

// Get returns a single Vault based on vaultId.
func (c *Vaults) Get(vaultId string) (Vault, error) {
	return c.GetWithContext(context.Background(), vaultId)
}

// GetWithContext returns a single Vault based on vaultId.
// The provided context can be used to cancel the request.
func (c *Vaults) GetWithContext(ctx context.Context, vaultId string) (Vault, error) {
	reqUrl, err := c.client.v2Url(nil, vaultEndpoint, vaultId)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	var vault Vault
	err = c.client.getJSON(ctx, reqUrl, &vault)
	if err != nil {
		return Vault{}, fmt.Errorf("error while fetching vault: %w", err)
	}

	return vault, nil
}

// GetByName returns a single Vault based on name. The name matches ignoring case.
// Returns ErrVaultNotFound if no vault is found.
// Returns ErrMultipleVaultsFound if more than one vault matches the name.
func (c *Vaults) GetByName(name string) (Vault, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single Vault based on name. The name matches ignoring case.
// Returns ErrVaultNotFound if no vault is found.
// Returns ErrMultipleVaultsFound if more than one vault matches the name.
// The provided context can be used to cancel the request.
func (c *Vaults) GetByNameWithContext(ctx context.Context, name string) (Vault, error) {
	vaults, err := c.list(ctx, url.Values{"name": {name}})
	if err != nil {
		return Vault{}, err
	}

	return singleByName(vaults, name, func(v Vault) string { return v.Name }, ErrVaultNotFound, ErrMultipleVaultsFound)
}

// New creates a new Vault and returns the created vault. Settings left empty take their server default.
func (c *Vaults) New(vault Vault) (Vault, error) {
	return c.NewWithContext(context.Background(), vault)
}

// NewWithContext creates a new Vault and returns the created vault. Settings left empty take their server default.
// The provided context can be used to cancel the request.
func (c *Vaults) NewWithContext(ctx context.Context, vault Vault) (Vault, error) {
	reqUrl, err := c.client.v2Url(nil, vaultEndpoint)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	reqBody := vaultCreateRequest{
		Name:          vault.Name,
		Description:   vault.Description,
		ContentType:   vault.ContentType,
		SecurityLevel: vault.SecurityLevel,
		Visibility:    vault.Visibility,
	}

	var createdVault Vault
	err = c.client.sendJSON(ctx, http.MethodPost, reqUrl, reqBody, &createdVault)
	if err != nil {
		return Vault{}, fmt.Errorf("error while creating vault: %w", err)
	}

	return createdVault, nil
}

// Update replaces the settings of an existing Vault and returns the updated vault.
// ContentType, SecurityLevel and Visibility are required.
func (c *Vaults) Update(vault Vault) (Vault, error) {
	return c.UpdateWithContext(context.Background(), vault)
}

// UpdateWithContext replaces the settings of an existing Vault and returns the updated vault.
// ContentType, SecurityLevel and Visibility are required.
// The provided context can be used to cancel the request.
func (c *Vaults) UpdateWithContext(ctx context.Context, vault Vault) (Vault, error) {
	reqUrl, err := c.client.v2Url(nil, vaultEndpoint, vault.Id)
	if err != nil {
		return Vault{}, fmt.Errorf("failed to build vault url: %w", err)
	}

	reqBody := vaultUpdateRequest{
		Name:          vault.Name,
		Description:   vault.Description,
		ContentType:   vault.ContentType,
		SecurityLevel: vault.SecurityLevel,
		Visibility:    vault.Visibility,
	}

	var updatedVault Vault
	err = c.client.sendJSON(ctx, http.MethodPut, reqUrl, reqBody, &updatedVault)
	if err != nil {
		return Vault{}, fmt.Errorf("error while updating vault: %w", err)
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
	reqUrl, err := c.client.v2Url(nil, vaultEndpoint, vaultId)
	if err != nil {
		return fmt.Errorf("failed to build vault url: %w", err)
	}

	err = c.client.sendJSON(ctx, http.MethodDelete, reqUrl, nil, nil)
	if err != nil {
		return fmt.Errorf("error while deleting vault: %w", err)
	}

	return nil
}

// ContentReadable reports whether the caller can read the entries of the vault vaultId. The entry routes answer the
// same 404 for a vault the caller cannot open as for an entry that does not exist: when an entry read answers 404, a
// readable vault means the entry is gone, an unreadable one that it may still exist.
// Returns false without an error when the server answers 403 or 404, and an error for any other failure.
func (c *Vaults) ContentReadable(vaultId string) (bool, error) {
	return c.ContentReadableWithContext(context.Background(), vaultId)
}

// ContentReadableWithContext reports whether the caller can read the entries of the vault vaultId. The entry routes
// answer the same 404 for a vault the caller cannot open as for an entry that does not exist: when an entry read
// answers 404, a readable vault means the entry is gone, an unreadable one that it may still exist.
// Returns false without an error when the server answers 403 or 404, and an error for any other failure.
// The provided context can be used to cancel the request.
func (c *Vaults) ContentReadableWithContext(ctx context.Context, vaultId string) (bool, error) {
	if vaultId == "" {
		return false, fmt.Errorf("vaultId is required")
	}

	reqUrl, err := entryUrl(c.client, url.Values{"pageSize": {"1"}}, vaultId)
	if err != nil {
		return false, err
	}

	_, err = c.client.RequestWithContext(ctx, reqUrl, http.MethodGet, nil)
	switch {
	case err == nil:
		return true, nil
	case IsNotFound(err) || hasStatusCode(err, http.StatusForbidden):
		return false, nil
	default:
		return false, fmt.Errorf("error while checking access to the vault content: %w", err)
	}
}
