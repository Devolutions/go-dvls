package dvls

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All folder subtypes to test
var folderSubTypes = []string{
	EntryFolderSubTypeCompany,
	EntryFolderSubTypeCredentials,
	EntryFolderSubTypeCustomer,
	EntryFolderSubTypeDatabase,
	EntryFolderSubTypeDevice,
	EntryFolderSubTypeDomain,
	EntryFolderSubTypeFolder,
	EntryFolderSubTypeIdentity,
	EntryFolderSubTypeMacroScriptTools,
	EntryFolderSubTypePrinter,
	EntryFolderSubTypeServer,
	EntryFolderSubTypeSite,
	EntryFolderSubTypeSmartFolder,
	EntryFolderSubTypeSoftware,
	EntryFolderSubTypeTeam,
	EntryFolderSubTypeWorkstation,
}

func Test_FolderCRUD(t *testing.T) {
	vault := createTestVault(t, "folders")

	for _, subType := range folderSubTypes {
		t.Run(subType, func(t *testing.T) {
			testPath := ""
			entryName := fmt.Sprintf("Test %s Folder", subType)
			description := fmt.Sprintf("Test %s folder entry", strings.ToLower(subType))

			// Initial data with domain and username
			initialDomain := fmt.Sprintf("%s.local", strings.ToLower(subType))
			initialUsername := fmt.Sprintf("%s-user", strings.ToLower(subType))

			// Create entry
			t.Logf("Creating %s folder with domain=%q, username=%q", subType, initialDomain, initialUsername)
			entry := Entry{
				VaultId:     vault.Id,
				Name:        entryName,
				Path:        testPath,
				Type:        EntryFolderType,
				SubType:     subType,
				Description: description,
				Tags:        []string{"test", strings.ToLower(subType)},
				Data: &EntryFolderData{
					Domain:   initialDomain,
					Username: initialUsername,
				},
			}

			id, err := testClient.Entries.Folder.New(entry)
			require.NoError(t, err, "Failed to create %s folder", subType)
			require.NotEmpty(t, id, "Entry ID should not be empty after creation")
			t.Logf("Created folder with ID: %s", id)

			// Get entry and verify domain/username
			t.Logf("Fetching folder %s", id)
			fetched, err := testClient.Entries.Folder.GetById(vault.Id, id)
			require.NoError(t, err, "Failed to get %s folder", subType)
			assert.Equal(t, entry.Name, fetched.Name)
			assert.Equal(t, entry.Description, fetched.Description)
			assert.Equal(t, EntryFolderType, fetched.Type, "Type should be Folder")
			assert.Equal(t, subType, fetched.SubType, "SubType should match")
			t.Logf("Verified type=%q, subType=%q", fetched.Type, fetched.SubType)

			// Verify data fields after creation
			data, ok := fetched.GetFolderData()
			require.True(t, ok, "Expected EntryFolderData type")
			assert.Equal(t, initialDomain, data.Domain, "Domain should match after creation")
			assert.Equal(t, initialUsername, data.Username, "Username should match after creation")
			t.Logf("Verified data: domain=%q, username=%q", data.Domain, data.Username)

			// Update entry with new domain and username
			updatedDomain := fmt.Sprintf("updated.%s.local", strings.ToLower(subType))
			updatedUsername := fmt.Sprintf("updated-%s-user", strings.ToLower(subType))
			newName := entryName + " (Updated)"
			newDescription := description + " - modified"

			t.Logf("Updating folder: domain=%q->%q, username=%q->%q", initialDomain, updatedDomain, initialUsername, updatedUsername)
			fetched.Name = newName
			fetched.Description = newDescription
			fetched.Tags = []string{"test", "updated"}
			fetched.Data = &EntryFolderData{
				Domain:   updatedDomain,
				Username: updatedUsername,
			}

			updated, err := testClient.Entries.Folder.Update(fetched)
			require.NoError(t, err, "Failed to update %s folder", subType)
			assert.Equal(t, newName, updated.Name)
			assert.Equal(t, newDescription, updated.Description)

			// Verify data fields after update
			updatedData, ok := updated.GetFolderData()
			require.True(t, ok, "Expected EntryFolderData type after update")
			assert.Equal(t, updatedDomain, updatedData.Domain, "Domain should match after update")
			assert.Equal(t, updatedUsername, updatedData.Username, "Username should match after update")
			t.Logf("Verified updated data: domain=%q, username=%q", updatedData.Domain, updatedData.Username)

			// Delete entry
			err = testClient.Entries.Folder.DeleteById(vault.Id, id)
			require.NoError(t, err, "Failed to delete %s folder", subType)

			// Verify deletion
			_, err = testClient.Entries.Folder.GetById(vault.Id, id)
			require.Error(t, err, "Entry should no longer exist after deletion")
		})
	}
}

func Test_NestedFolders(t *testing.T) {
	vault := createTestVault(t, "nested-folders")

	// Create parent folder at root
	parentEntry := Entry{
		VaultId:     vault.Id,
		Name:        "Parent Folder",
		Path:        "",
		Type:        EntryFolderType,
		SubType:     EntryFolderSubTypeFolder,
		Description: "Parent folder",
		Data:        &EntryFolderData{},
	}

	parentId, err := testClient.Entries.Folder.New(parentEntry)
	require.NoError(t, err, "Failed to create parent folder")
	t.Logf("Created parent folder with ID: %s", parentId)

	// Fetch parent
	parent, err := testClient.Entries.Folder.GetById(vault.Id, parentId)
	require.NoError(t, err, "Failed to fetch parent folder")
	t.Logf("Parent folder: Name=%q, Path=%q", parent.Name, parent.Path)

	// Create child folder inside parent
	childEntry := Entry{
		VaultId:     vault.Id,
		Name:        "Child Folder",
		Path:        parent.Name,
		Type:        EntryFolderType,
		SubType:     EntryFolderSubTypeServer,
		Description: "Child folder inside parent",
		Data:        &EntryFolderData{},
	}

	childId, err := testClient.Entries.Folder.New(childEntry)
	require.NoError(t, err, "Failed to create child folder")
	t.Logf("Created child folder with ID: %s", childId)

	// Fetch child and verify
	child, err := testClient.Entries.Folder.GetById(vault.Id, childId)
	require.NoError(t, err, "Failed to fetch child folder")
	t.Logf("Child folder: Name=%q, Path=%q, SubType=%q", child.Name, child.Path, child.SubType)

	assert.Equal(t, "Child Folder", child.Name)
	assert.Equal(t, EntryFolderSubTypeServer, child.SubType)

	// Create grandchild folder inside child
	grandchildEntry := Entry{
		VaultId:     vault.Id,
		Name:        "Grandchild Folder",
		Path:        fmt.Sprintf("%s\\%s", parent.Name, child.Name),
		Type:        EntryFolderType,
		SubType:     EntryFolderSubTypeDatabase,
		Description: "Grandchild folder inside child",
		Data:        &EntryFolderData{},
	}

	grandchildId, err := testClient.Entries.Folder.New(grandchildEntry)
	require.NoError(t, err, "Failed to create grandchild folder")
	t.Logf("Created grandchild folder with ID: %s", grandchildId)

	// Fetch grandchild and verify
	grandchild, err := testClient.Entries.Folder.GetById(vault.Id, grandchildId)
	require.NoError(t, err, "Failed to fetch grandchild folder")
	t.Logf("Grandchild folder: Name=%q, Path=%q, SubType=%q", grandchild.Name, grandchild.Path, grandchild.SubType)

	assert.Equal(t, "Grandchild Folder", grandchild.Name)
	assert.Equal(t, EntryFolderSubTypeDatabase, grandchild.SubType)

	// Delete entries (in reverse order)
	err = testClient.Entries.Folder.DeleteById(vault.Id, grandchildId)
	require.NoError(t, err, "Failed to delete grandchild folder")
	err = testClient.Entries.Folder.DeleteById(vault.Id, childId)
	require.NoError(t, err, "Failed to delete child folder")
	err = testClient.Entries.Folder.DeleteById(vault.Id, parentId)
	require.NoError(t, err, "Failed to delete parent folder")
}

func Test_GetFolderByName(t *testing.T) {
	vault := createTestVault(t, "folder-getbyname")
	testPath := "go-dvls\\folder-getbyname"

	entry := Entry{
		VaultId: vault.Id,
		Name:    "MyFolder",
		Path:    testPath,
		Type:    EntryFolderType,
		SubType: EntryFolderSubTypeFolder,
		Data:    &EntryFolderData{Domain: "test.local", Username: "testuser"},
	}

	id, err := testClient.Entries.Folder.New(entry)
	require.NoError(t, err, "Failed to create folder entry")
	t.Cleanup(func() {
		_ = testClient.Entries.Folder.DeleteById(vault.Id, id)
	})

	// GetByName with name only
	got, err := testClient.Entries.Folder.GetByName(vault.Id, "MyFolder", GetByNameOptions{Path: &testPath})
	require.NoError(t, err)
	assert.Equal(t, id, got.Id)
	assert.Equal(t, "MyFolder", got.Name)
	assert.Equal(t, testPath+`\MyFolder`, got.Path)

	// GetByName with non-existent name returns ErrEntryNotFound
	_, err = testClient.Entries.Folder.GetByName(vault.Id, "NonExistentFolder", GetByNameOptions{Path: &testPath})
	assert.ErrorIs(t, err, ErrEntryNotFound)

	// GetByName without path filter also finds the entry
	got, err = testClient.Entries.Folder.GetByName(vault.Id, "MyFolder", GetByNameOptions{})
	require.NoError(t, err)
	assert.Equal(t, id, got.Id)

	// Root-level folder: path returned by API is the folder name itself
	rootEntry := Entry{
		VaultId: vault.Id,
		Name:    "MyRootFolder",
		Path:    "",
		Type:    EntryFolderType,
		SubType: EntryFolderSubTypeFolder,
		Data:    &EntryFolderData{},
	}
	rootId, err := testClient.Entries.Folder.New(rootEntry)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testClient.Entries.Folder.DeleteById(vault.Id, rootId)
	})

	root, err := testClient.Entries.Folder.GetByName(vault.Id, "MyRootFolder", GetByNameOptions{})
	require.NoError(t, err)
	assert.Equal(t, rootId, root.Id)
	assert.Equal(t, "MyRootFolder", root.Path)
}

func Test_GetFolderEntries_Filters(t *testing.T) {
	vault := createTestVault(t, "folder-getentries")
	testPath := "go-dvls\\folder-getentries"

	// Create 3 test folder entries - "Database" is exact match, others contain "Database" in name
	entriesToCreate := []Entry{
		{
			VaultId:     vault.Id,
			Name:        "Database",
			Path:        testPath,
			Type:        EntryFolderType,
			SubType:     EntryFolderSubTypeDatabase,
			Description: "Exact match folder",
			Data:        &EntryFolderData{Domain: "db.local", Username: "dbuser"},
		},
		{
			VaultId:     vault.Id,
			Name:        "Database Backup",
			Path:        testPath,
			Type:        EntryFolderType,
			SubType:     EntryFolderSubTypeDatabase,
			Description: "Contains Database in name",
			Data:        &EntryFolderData{Domain: "backup.local", Username: "backupuser"},
		},
		{
			VaultId:     vault.Id,
			Name:        "Database Production",
			Path:        testPath,
			Type:        EntryFolderType,
			SubType:     EntryFolderSubTypeDatabase,
			Description: "Contains Database in name",
			Data:        &EntryFolderData{Domain: "prod.local", Username: "produser"},
		},
	}

	// Create test entries
	t.Log("Creating test folder entries for GetEntries")
	var createdIds []string
	for _, entry := range entriesToCreate {
		id, err := testClient.Entries.Folder.New(entry)
		require.NoError(t, err, "Failed to create folder entry %s", entry.Name)
		createdIds = append(createdIds, id)
		t.Logf("Created folder entry %q with ID: %s", entry.Name, id)
	}

	databaseName := "Database"
	databaseBackupName := "Database Backup"
	nonExistentName := "Non Existent Folder"

	// Test 1: GetEntries with path filter should return at least our 3 folders
	// Note: DVLS may auto-create parent folders, so we check for >= 3
	t.Log("Test 1: GetEntries with path filter")
	entries, err := testClient.Entries.Folder.GetEntries(vault.Id, GetEntriesOptions{Path: &testPath})
	require.NoError(t, err, "GetEntries failed")
	assert.GreaterOrEqual(t, len(entries), 3, "Expected at least 3 folder entries with path filter")

	// Verify our 3 folders are present
	foundNames := make(map[string]bool)
	for _, e := range entries {
		foundNames[e.Name] = true
	}
	assert.True(t, foundNames["Database"], "Expected to find 'Database' folder")
	assert.True(t, foundNames["Database Backup"], "Expected to find 'Database Backup' folder")
	assert.True(t, foundNames["Database Production"], "Expected to find 'Database Production' folder")
	t.Logf("Found %d folder entries in path %q (including auto-created parent folders)", len(entries), testPath)

	// Test 2: GetEntries with exact name match - should return only "Database"
	t.Log("Test 2: GetEntries with exact name match")
	entries, err = testClient.Entries.Folder.GetEntries(vault.Id, GetEntriesOptions{Name: &databaseName})
	require.NoError(t, err, "GetEntries with exact name failed")
	assert.Len(t, entries, 1, "Expected 1 folder entry with exact name match")
	if len(entries) > 0 {
		assert.Equal(t, "Database", entries[0].Name)
		t.Logf("Found exact match: %q", entries[0].Name)
	}

	// Test 3: GetEntries with name and path filter
	t.Log("Test 3: GetEntries with name and path filter")
	entries, err = testClient.Entries.Folder.GetEntries(vault.Id, GetEntriesOptions{Name: &databaseBackupName, Path: &testPath})
	require.NoError(t, err, "GetEntries with name and path filter failed")
	assert.Len(t, entries, 1, "Expected 1 folder entry with name and path filter")
	t.Logf("Found %d folder entry with combined filters", len(entries))

	// Test 4: GetEntries with non-existent name should return empty
	t.Log("Test 4: GetEntries with non-existent name")
	entries, err = testClient.Entries.Folder.GetEntries(vault.Id, GetEntriesOptions{Name: &nonExistentName, Path: &testPath})
	require.NoError(t, err, "GetEntries with non-existent name failed")
	assert.Empty(t, entries, "Expected 0 folder entries for non-existent name")
	t.Logf("Correctly returned %d folder entries for non-existent name", len(entries))

	// Cleanup test entries
	t.Log("Cleaning up test folder entries")
	for _, id := range createdIds {
		err := testClient.Entries.Folder.DeleteById(vault.Id, id)
		require.NoError(t, err, "Failed to delete folder entry %s", id)
	}
	t.Log("Cleanup complete")
}
