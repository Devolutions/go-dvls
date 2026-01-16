package dvls

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// credentialTestCase defines a test case for credential CRUD operations.
type credentialTestCase struct {
	name        string
	entryName   string
	description string
	subType     string
	data        EntryData
	updateData  func(entry *Entry)
}

var credentialTestCases = []credentialTestCase{
	{
		name:        "AccessCode",
		entryName:   "Test Access Code",
		description: "Test access code entry",
		subType:     EntryCredentialSubTypeAccessCode,
		data:        &EntryCredentialAccessCodeData{Password: "1234"},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialAccessCodeData(); ok {
				data.Password = "5678"
				entry.Data = data
			}
		},
	},
	{
		name:        "ApiKey",
		entryName:   "Test API Key",
		description: "Test API key entry",
		subType:     EntryCredentialSubTypeApiKey,
		data: &EntryCredentialApiKeyData{
			ApiId:    "test-api-id",
			ApiKey:   "test-api-key",
			TenantId: "test-tenant",
		},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialApiKeyData(); ok {
				data.ApiKey = "test-api-key-updated"
				entry.Data = data
			}
		},
	},
	{
		name:        "AzureServicePrincipal",
		entryName:   "Test Azure Service Principal",
		description: "Test Azure service principal entry",
		subType:     EntryCredentialSubTypeAzureServicePrincipal,
		data: &EntryCredentialAzureServicePrincipalData{
			ClientId:     "test-client-id",
			ClientSecret: "test-client-secret",
			TenantId:     "test-tenant-id",
		},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialAzureServicePrincipalData(); ok {
				data.ClientSecret = "test-client-secret-updated"
				entry.Data = data
			}
		},
	},
	{
		name:        "ConnectionString",
		entryName:   "Test Connection String",
		description: "Test connection string entry",
		subType:     EntryCredentialSubTypeConnectionString,
		data: &EntryCredentialConnectionStringData{
			ConnectionString: "Server=localhost;Database=testdb;",
		},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialConnectionStringData(); ok {
				data.ConnectionString = "Server=localhost;Database=testdb;Encrypt=True;"
				entry.Data = data
			}
		},
	},
	{
		name:        "Default",
		entryName:   "Test Username Password",
		description: "Test username/password entry",
		subType:     EntryCredentialSubTypeDefault,
		data: &EntryCredentialDefaultData{
			Domain:   "example.com",
			Username: "testuser",
			Password: "testpass",
		},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialDefaultData(); ok {
				data.Password = "testpass-updated"
				entry.Data = data
			}
		},
	},
	{
		name:        "PrivateKey",
		entryName:   "Test Private Key",
		description: "Test private key entry",
		subType:     EntryCredentialSubTypePrivateKey,
		data: &EntryCredentialPrivateKeyData{
			Username:   "testuser",
			Password:   "testpass",
			PrivateKey: "-----BEGIN PRIVATE KEY-----\ntestkey\n-----END PRIVATE KEY-----",
			PublicKey:  "-----BEGIN PUBLIC KEY-----\ntestkey\n-----END PUBLIC KEY-----",
			Passphrase: "testpassphrase",
		},
		updateData: func(entry *Entry) {
			if data, ok := entry.GetCredentialPrivateKeyData(); ok {
				data.Passphrase = "testpassphrase-updated"
				entry.Data = data
			}
		},
	},
}

func Test_CredentialCRUD(t *testing.T) {
	for _, tc := range credentialTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testPath := "go-dvls\\credentials\\" + strings.ToLower(tc.name)

			// Create entry
			t.Logf("Creating %s entry: %q", tc.subType, tc.entryName)
			entry := Entry{
				VaultId:     testVaultId,
				Name:        tc.entryName,
				Path:        testPath,
				Type:        EntryCredentialType,
				SubType:     tc.subType,
				Description: tc.description,
				Tags:        []string{"test", strings.ToLower(tc.name)},
				Data:        tc.data,
			}

			id, err := testClient.Entries.Credential.New(entry)
			require.NoError(t, err, "Failed to create %s entry", tc.name)
			require.NotEmpty(t, id, "Entry ID should not be empty after creation")
			t.Logf("Created entry with ID: %s", id)

			// Get entry
			t.Logf("Fetching entry %s", id)
			fetched, err := testClient.Entries.Credential.GetById(testVaultId, id)
			require.NoError(t, err, "Failed to get %s entry", tc.name)
			assert.Equal(t, entry.Name, fetched.Name)
			assert.Equal(t, entry.Description, fetched.Description)
			t.Logf("Fetched entry: Name=%q, Path=%q", fetched.Name, fetched.Path)

			// Update entry
			newName := tc.entryName + " (Updated)"
			newDescription := tc.description + " - modified"
			t.Logf("Updating entry: %q -> %q", fetched.Name, newName)
			fetched.Name = newName
			fetched.Description = newDescription
			fetched.Tags = []string{"test", "updated"}
			tc.updateData(&fetched)

			updated, err := testClient.Entries.Credential.Update(fetched)
			require.NoError(t, err, "Failed to update %s entry", tc.name)
			assert.Equal(t, newName, updated.Name)
			assert.Equal(t, newDescription, updated.Description)
			t.Logf("Updated entry successfully")

			// Delete entry
			t.Logf("Deleting entry %s", id)
			err = testClient.Entries.Credential.DeleteById(testVaultId, id)
			require.NoError(t, err, "Failed to delete %s entry", tc.name)

			// Verify deletion
			_, err = testClient.Entries.Credential.GetById(testVaultId, id)
			assert.Error(t, err, "Entry should not exist after deletion")
			t.Logf("Entry deleted and verified")
		})
	}
}

func Test_GetEntries(t *testing.T) {
	testPath := "go-dvls\\getentries"

	// Create 3 test entries - "Server" is exact match, others contain "Server" in name
	entriesToCreate := []Entry{
		{
			VaultId:     testVaultId,
			Name:        "Server",
			Path:        testPath,
			Type:        EntryCredentialType,
			SubType:     EntryCredentialSubTypeDefault,
			Description: "Exact match entry",
			Data:        &EntryCredentialDefaultData{Username: "testuser", Password: "testpass"},
		},
		{
			VaultId:     testVaultId,
			Name:        "Server Backup",
			Path:        testPath,
			Type:        EntryCredentialType,
			SubType:     EntryCredentialSubTypeDefault,
			Description: "Contains Server in name",
			Data:        &EntryCredentialDefaultData{Username: "testuser", Password: "testpass"},
		},
		{
			VaultId:     testVaultId,
			Name:        "Server Production",
			Path:        testPath,
			Type:        EntryCredentialType,
			SubType:     EntryCredentialSubTypeDefault,
			Description: "Contains Server in name",
			Data:        &EntryCredentialDefaultData{Username: "testuser", Password: "testpass"},
		},
	}

	// Create test entries
	t.Log("Creating test entries for GetEntries")
	var createdIds []string
	for _, entry := range entriesToCreate {
		id, err := testClient.Entries.Credential.New(entry)
		require.NoError(t, err, "Failed to create entry %s", entry.Name)
		createdIds = append(createdIds, id)
		t.Logf("Created entry %q with ID: %s", entry.Name, id)
	}

	// Test 1: GetEntries with path filter should return all 3 entries
	t.Log("Test 1: GetEntries with path filter")
	entries, err := testClient.Entries.Credential.GetEntries(testVaultId, "", testPath)
	require.NoError(t, err, "GetEntries failed")
	assert.Len(t, entries, 3, "Expected 3 entries with path filter")
	t.Logf("Found %d entries in path %q", len(entries), testPath)

	// Test 2: GetEntries with exact name match - should return only "Server", not "Server Backup" or "Server Production"
	t.Log("Test 2: GetEntries with exact name match")
	entries, err = testClient.Entries.Credential.GetEntries(testVaultId, "Server", "")
	require.NoError(t, err, "GetEntries with exact name failed")
	assert.Len(t, entries, 1, "Expected 1 entry with exact name match")
	if len(entries) > 0 {
		assert.Equal(t, "Server", entries[0].Name)
		t.Logf("Found exact match: %q", entries[0].Name)
	}

	// Test 3: GetEntries with name and path filter
	t.Log("Test 3: GetEntries with name and path filter")
	entries, err = testClient.Entries.Credential.GetEntries(testVaultId, "Server Backup", testPath)
	require.NoError(t, err, "GetEntries with name and path filter failed")
	assert.Len(t, entries, 1, "Expected 1 entry with name and path filter")
	t.Logf("Found %d entry with combined filters", len(entries))

	// Test 4: GetEntries with non-existent name should return empty
	t.Log("Test 4: GetEntries with non-existent name")
	entries, err = testClient.Entries.Credential.GetEntries(testVaultId, "Non Existent Entry", testPath)
	require.NoError(t, err, "GetEntries with non-existent name failed")
	assert.Empty(t, entries, "Expected 0 entries for non-existent name")
	t.Logf("Correctly returned %d entries for non-existent name", len(entries))

	// Cleanup test entries
	t.Log("Cleaning up test entries")
	for _, id := range createdIds {
		err := testClient.Entries.Credential.DeleteById(testVaultId, id)
		require.NoError(t, err, "Failed to delete entry %s", id)
	}
	t.Log("Cleanup complete")
}
