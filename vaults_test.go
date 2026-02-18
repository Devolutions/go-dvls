package dvls

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Vaults(t *testing.T) {
	t.Run("ListVaults", test_ListVaults)
	t.Run("GetVault", test_GetVault)
	t.Run("GetVaultByName", test_GetVaultByName)
	t.Run("GetVaultByName_NotFound", test_GetVaultByName_NotFound)
	t.Run("NewVault", test_NewVault)
	t.Run("UpdateVault", test_UpdateVault)
	t.Run("ContentType_DefaultEquivalence", test_ContentType_DefaultEquivalence)
}

func test_ListVaults(t *testing.T) {
	vault := createTestVault(t, "list-vaults")

	vaults, err := testClient.Vaults.List()
	require.NoError(t, err)
	assert.NotEmpty(t, vaults)

	found := false
	for _, v := range vaults {
		if v.Id == vault.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "expected test vault to be in the list")
}

func test_GetVault(t *testing.T) {
	vault := createTestVault(t, "get-vault")

	fetchedVault, err := testClient.Vaults.Get(vault.Id)
	require.NoError(t, err)
	assert.Equal(t, vault.Id, fetchedVault.Id)
	assert.Equal(t, vault.Name, fetchedVault.Name)
}

func test_GetVaultByName(t *testing.T) {
	vault := createTestVault(t, "get-by-name")

	// Test GetByName with the created vault's name
	foundVault, err := testClient.Vaults.GetByName(vault.Name)
	require.NoError(t, err)
	assert.Equal(t, vault.Id, foundVault.Id)
	assert.Equal(t, vault.Name, foundVault.Name)
}

func test_GetVaultByName_NotFound(t *testing.T) {
	_, err := testClient.Vaults.GetByName("nonexistent-vault-name-12345")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrVaultNotFound))
}

func test_NewVault(t *testing.T) {
	tests := []struct {
		name  string
		vault Vault
	}{
		{
			name: "Standard/Default/Default",
			vault: Vault{
				Name:          "test-standard-default",
				Description:   "Test vault",
				ContentType:   VaultContentTypeEverything,
				SecurityLevel: VaultSecurityLevelStandard,
				Visibility:    VaultVisibilityDefault,
			},
		},
		{
			name: "High/Everyone/Secrets",
			vault: Vault{
				Name:          "test-high-everyone",
				Description:   "High security public vault",
				ContentType:   VaultContentTypeSecrets,
				SecurityLevel: VaultSecurityLevelHigh,
				Visibility:    VaultVisibilityPublic,
			},
		},
		{
			name: "Standard/Never/Credentials",
			vault: Vault{
				Name:          "test-credentials",
				Description:   "Credentials vault",
				ContentType:   VaultContentTypeCredentials,
				SecurityLevel: VaultSecurityLevelStandard,
				Visibility:    VaultVisibilityPrivate,
			},
		},
		{
			name: "High/Never/BusinessInformation",
			vault: Vault{
				Name:          "test-business",
				Description:   "Business info vault",
				ContentType:   VaultContentTypeBusinessInformation,
				SecurityLevel: VaultSecurityLevelHigh,
				Visibility:    VaultVisibilityPrivate,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := testClient.Vaults.New(tt.vault)
			require.NoError(t, err)
			require.NotEmpty(t, created.Id)

			// Register cleanup to ensure vault deletion even if test fails
			t.Cleanup(func() {
				testClient.Vaults.Delete(created.Id)
			})

			fetched, err := testClient.Vaults.Get(created.Id)
			require.NoError(t, err)
			assert.Equal(t, tt.vault.Name, fetched.Name)
			assert.Equal(t, tt.vault.Description, fetched.Description)
			assert.Equal(t, tt.vault.ContentType, fetched.ContentType)
			assert.Equal(t, tt.vault.SecurityLevel, fetched.SecurityLevel)
			assert.Equal(t, tt.vault.Visibility, fetched.Visibility)
		})
	}
}

func test_UpdateVault(t *testing.T) {
	originalVault := Vault{
		Name:          "test-update-vault",
		Description:   "Original description",
		ContentType:   VaultContentTypeEverything,
		SecurityLevel: VaultSecurityLevelStandard,
		Visibility:    VaultVisibilityDefault,
	}

	created, err := testClient.Vaults.New(originalVault)
	require.NoError(t, err)

	// Register cleanup to ensure vault deletion even if test fails
	t.Cleanup(func() {
		testClient.Vaults.Delete(created.Id)
	})

	tests := []struct {
		name   string
		update func(v *Vault)
		verify func(t *testing.T, v Vault)
	}{
		{
			name: "UpdateName",
			update: func(v *Vault) {
				v.Name = "test-update-vault-renamed"
			},
			verify: func(t *testing.T, v Vault) {
				assert.Equal(t, "test-update-vault-renamed", v.Name)
			},
		},
		{
			name: "UpdateDescription",
			update: func(v *Vault) {
				v.Description = "Updated description"
			},
			verify: func(t *testing.T, v Vault) {
				assert.Equal(t, "Updated description", v.Description)
			},
		},
		{
			name: "UpdateSecurityLevel",
			update: func(v *Vault) {
				v.SecurityLevel = VaultSecurityLevelHigh
			},
			verify: func(t *testing.T, v Vault) {
				assert.Equal(t, VaultSecurityLevelHigh, v.SecurityLevel)
			},
		},
		{
			name: "UpdateVisibility",
			update: func(v *Vault) {
				v.Visibility = VaultVisibilityPublic
			},
			verify: func(t *testing.T, v Vault) {
				assert.Equal(t, VaultVisibilityPublic, v.Visibility)
			},
		},
		{
			name: "UpdateContentType",
			update: func(v *Vault) {
				v.ContentType = VaultContentTypeSecrets
			},
			verify: func(t *testing.T, v Vault) {
				assert.Equal(t, VaultContentTypeSecrets, v.ContentType)
			},
		},
	}

	currentVault := created
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.update(&currentVault)
			updated, err := testClient.Vaults.Update(currentVault)
			require.NoError(t, err)

			fetched, err := testClient.Vaults.Get(updated.Id)
			require.NoError(t, err)
			tt.verify(t, fetched)

			currentVault = fetched
		})
	}
}

// test_ContentType_DefaultEquivalence verifies that:
// 1. System vaults (Default, User vault) return VaultContentTypeDefault ("Default")
// 2. VaultContentTypeDefault is automatically converted to VaultContentTypeEverything on creation
func test_ContentType_DefaultEquivalence(t *testing.T) {
	// Verify system vaults use "Default"
	vault, err := testClient.Vaults.Get(testVaultId)
	require.NoError(t, err)
	assert.True(t,
		vault.ContentType == VaultContentTypeEverything || vault.ContentType == VaultContentTypeDefault,
		"expected ContentType to be 'Everything' or 'Default', got %q", vault.ContentType)

	// Verify that using VaultContentTypeDefault in New() works (converted to Everything)
	newVault := Vault{
		Name:          "test-default-conversion",
		Description:   "Test Default to Everything conversion",
		ContentType:   VaultContentTypeDefault,
		SecurityLevel: VaultSecurityLevelStandard,
		Visibility:    VaultVisibilityDefault,
	}

	created, err := testClient.Vaults.New(newVault)
	require.NoError(t, err, "creating vault with VaultContentTypeDefault should work")
	assert.Equal(t, VaultContentTypeEverything, created.ContentType)

	err = testClient.Vaults.Delete(created.Id)
	require.NoError(t, err)
}
