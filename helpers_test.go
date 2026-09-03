//go:build integration

package dvls

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// waitUntil polls probe every 200ms until it returns true, failing the test
// after the timeout. The probe runs once before the first tick so an already
// satisfied condition returns immediately.
func waitUntil(t *testing.T, timeout time.Duration, description string, probe func() bool) {
	t.Helper()

	deadline := time.After(timeout)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		if probe() {
			return
		}

		select {
		case <-deadline:
			t.Fatalf("timeout waiting for %s", description)
		case <-ticker.C:
		}
	}
}

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

	waitUntil(t, 5*time.Second, fmt.Sprintf("vault %s to be indexed", vault.Id), func() bool {
		_, err := testClient.Vaults.Get(vault.Id)
		return err == nil
	})

	return vault
}

// createdTestFolderPaths tracks folders already verified or created, keyed by
// vault id and folder path, to avoid re-checking shared parents across tests.
var createdTestFolderPaths = map[string]bool{}

// createTestFolderPath creates every missing folder of a backslash-separated
// path in the vault, parents first. The server rejects entry creation into a
// non-existent folder, so tests must create the folder hierarchy before
// creating entries at a nested path. Each created folder is polled until it
// is indexed (max 10s timeout). Folders are deleted with the test vault.
func createTestFolderPath(t *testing.T, vaultId string, path string) {
	t.Helper()

	parentPath := ""
	for segment := range strings.SplitSeq(path, "\\") {
		// A folder entry's path is its own full path, not its parent's.
		folderPath := segment
		if parentPath != "" {
			folderPath = parentPath + "\\" + segment
		}

		if createdTestFolderPaths[vaultId+"|"+folderPath] {
			parentPath = folderPath
			continue
		}

		_, err := testClient.Entries.Folder.GetByName(vaultId, segment, GetByNameOptions{Path: &folderPath})
		if errors.Is(err, ErrEntryNotFound) {
			_, err = testClient.Entries.Folder.New(Entry{
				VaultId: vaultId,
				Name:    segment,
				Path:    parentPath,
				Type:    EntryFolderType,
				SubType: EntryFolderSubTypeFolder,
				Data:    &EntryFolderData{},
			})
			require.NoError(t, err, "failed to create folder %q under %q", segment, parentPath)

			waitUntil(t, 10*time.Second, fmt.Sprintf("folder %q to be indexed", folderPath), func() bool {
				_, err := testClient.Entries.Folder.GetByName(vaultId, segment, GetByNameOptions{Path: &folderPath})
				return err == nil
			})
		} else {
			require.NoError(t, err, "failed to look up folder %q under %q", segment, parentPath)
		}

		createdTestFolderPaths[vaultId+"|"+folderPath] = true
		parentPath = folderPath
	}
}
