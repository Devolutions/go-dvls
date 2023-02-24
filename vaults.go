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
	ID           string
	Name         string
	Description  string
	CreationDate *ServerTime
	ModifiedDate *ServerTime
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vault) UnmarshalJSON(b []byte) error {
	type rawVault Vault
	var raw struct {
		Data rawVault
	}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	*v = Vault(raw.Data)

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v Vault) MarshalJSON() ([]byte, error) {
	raw := struct {
		Description            string `json:"description"`
		HasPasswordChanged     bool   `json:"hasPasswordChanged"`
		Id                     string `json:"id"`
		IdString               string `json:"idString"`
		Image                  string `json:"image"`
		ImageBytes             string `json:"imageBytes"`
		ImageName              string `json:"imageName"`
		IsAllowedOffline       bool   `json:"isAllowedOffline"`
		IsLocked               bool   `json:"isLocked"`
		IsPrivate              bool   `json:"isPrivate"`
		ModifiedLoggedUserName string `json:"modifiedLoggedUserName"`
		ModifiedUserName       string `json:"modifiedUserName"`
		Name                   string `json:"name"`
		RepositorySettings     struct {
			QuickAddEntries    [0]struct{} `json:"quickAddEntries"`
			MasterPasswordHash *string     `json:"masterPasswordHash"`
		} `json:"repositorySettings"`
		Selected bool `json:"selected"`
	}{}

	raw.Name = v.Name
	raw.Description = v.Description
	raw.Id = v.ID
	raw.IdString = v.ID

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
func (c *Client) NewVault(vault Vault) error {
	reqUrl, err := url.JoinPath(c.baseUri, vaultEndpoint)
	if err != nil {
		return fmt.Errorf("failed to build vault url. error: %w", err)
	}

	vault.CreationDate = nil
	vault.ModifiedDate = nil

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
func (c *Client) UpdateVault(vault Vault) error {
	_, err := c.GetVault(vault.ID)
	if err != nil {
		return fmt.Errorf("error while fetching vault. error: %w", err)
	}
	err = c.NewVault(vault)
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
