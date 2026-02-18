package dvls

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createTestVault creates a vault for testing and registers cleanup.
// The vault name reflects the test being performed.
// Polls until the vault is indexed and ready to use (max 5s timeout).
func createTestVault(t *testing.T, name string) Vault {
	t.Helper()
	vault, err := testClient.Vaults.New(Vault{
		Name:          fmt.Sprintf("test-%s", name),
		Description:   "Auto-created test vault",
		ContentType:   VaultContentTypeEverything,
		SecurityLevel: VaultSecurityLevelStandard,
		Visibility:    VaultVisibilityDefault,
	})

	require.NoError(t, err)

	// Register cleanup immediately after creation to ensure deletion even if polling times out
	t.Cleanup(func() {
		testClient.Vaults.Delete(vault.Id)
	})

	// Wait for vault to be indexed by polling
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for vault %s to be indexed", vault.Id)
		case <-ticker.C:
			_, err := testClient.Vaults.Get(vault.Id)
			if err == nil {
				// Vault is indexed and ready
				return vault
			}
		}
	}
}
