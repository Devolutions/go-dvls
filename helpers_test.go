package dvls

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createTestVault creates a vault for testing and registers cleanup.
// The vault name reflects the test being performed.
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

	// Wait for vault to be fully indexed
	time.Sleep(2 * time.Second)

	t.Cleanup(func() {
		testClient.Vaults.Delete(vault.Id)
	})

	return vault
}
