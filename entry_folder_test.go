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
