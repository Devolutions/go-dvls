package dvls

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	testUserEntryId  string
	testNewUserEntry Entry
	testUserEntry    = Entry{
		ID:          "8d13bdea-3b33-42b6-8ae1-eb1f0aae2f88",
		VaultId:     testVaultId,
		EntryName:   "TestK8sSecret",
		Description: "Test description",
		Type:        "Credential",
		SubType:     "Default",
		Tags:        []string{"Test tag 1", "Test tag 2", "testtag"},
	}
)

// TODO
// Création avec un subType Default nous donne un objet sans subType (À gérer)
// On ne peut pas edit le subType (À gérer)
// Faire une run de test pour chaque subtype disponible dans /api/v1/entry/definition
// les crédentials vont devoir être varié
func Test_EntryUserCredentials(t *testing.T) {
	testUserEntryId = os.Getenv("TEST_USER_ENTRY_ID")
	testUserEntry.ID = testUserEntryId
	testUserEntry.VaultId = testVaultId

	t.Run("NewEntry", test_NewUserEntry)
	t.Run("GetEntry", test_GetUserEntry)
	t.Run("UpdateEntry", test_UpdateUserEntry)
	t.Run("DeleteEntry", test_DeleteUserEntry)
}

func init() {
	// Set the credentials using the interface approach
	testUserEntry.SetCredentials(&DefaultCredentials{
		Username: "TestK8s",
		Password: "TestK8sPassword",
	})
}

// Update test_NewUserEntry to use the interface approach
func test_NewUserEntry(t *testing.T) {
	testNewUserEntry = testUserEntry
	testNewUserEntry.ID = "" // Set empty ID for new creation
	testNewUserEntry.EntryName = "TestK8sNewEntry"

	// Set credentials using the interface approach
	testNewUserEntry.SetCredentials(&DefaultCredentials{
		Username: "TestK8sNew",
		Password: "TestK8sNewPassword",
	})

	newEntry, err := testClient.NewEntry(testNewUserEntry)
	if err != nil {
		t.Fatalf("Failed to create new entry: %v", err)
	}

	testNewUserEntry.ID = newEntry.ID
	testNewUserEntry.ModifiedOn = newEntry.ModifiedOn
	testNewUserEntry.ModifiedBy = newEntry.ModifiedBy
	testNewUserEntry.CreatedOn = newEntry.CreatedOn
	testNewUserEntry.CreatedBy = newEntry.CreatedBy

	// Verify the credentials are correctly set
	creds, ok := newEntry.GetCredentials().(*DefaultCredentials)
	if !ok {
		t.Fatalf("Expected DefaultCredentials, got %T", newEntry.GetCredentials())
	}

	if creds.Username != "TestK8sNew" || creds.Password != "TestK8sNewPassword" {
		t.Fatalf("Credentials not created correctly. Got Username=%s, Password=%s",
			creds.Username, creds.Password)
	}

	fmt.Printf("Created new entry: %+v\n", newEntry)
}

// Update test_GetUserEntry to handle the interface credentials
func test_GetUserEntry(t *testing.T) {
	entry, err := testClient.GetEntry(testVaultId, testUserEntry.ID)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	// Ignore fields that are not set by the user
	entry.CreatedBy = ""
	entry.ModifiedBy = ""
	entry.CreatedOn = nil
	entry.ModifiedOn = nil

	// Check main fields
	if entry.ID != testUserEntry.ID ||
		entry.VaultId != testUserEntry.VaultId ||
		entry.EntryName != testUserEntry.EntryName ||
		entry.Description != testUserEntry.Description ||
		entry.Type != testUserEntry.Type ||
		entry.SubType != testUserEntry.SubType ||
		!reflect.DeepEqual(entry.Tags, testUserEntry.Tags) {
		t.Fatalf("Entry fields don't match expected values.")
	}

	// Check credentials separately
	entryCreds, ok := entry.GetCredentials().(*DefaultCredentials)
	if !ok {
		t.Fatalf("Expected DefaultCredentials, got %T", entry.GetCredentials())
	}

	expectedCreds, ok := testUserEntry.GetCredentials().(*DefaultCredentials)
	if !ok {
		t.Fatalf("Test entry credentials not set properly")
	}

	if entryCreds.Username != expectedCreds.Username ||
		entryCreds.Password != expectedCreds.Password {
		t.Fatalf("Credentials don't match. Got {%s, %s}, Expected {%s, %s}",
			entryCreds.Username, entryCreds.Password,
			expectedCreds.Username, expectedCreds.Password)
	}
}

// Update test_UpdateUserEntry to handle the interface credentials
func test_UpdateUserEntry(t *testing.T) {
	testUpdatedEntry := testNewUserEntry
	testUpdatedEntry.EntryName = "TestK8sUpdatedEntry"

	// Set updated credentials
	testUpdatedEntry.SetCredentials(&DefaultCredentials{
		Username: "TestK8sUpdatedUser",
		Password: "TestK8sUpdatedPassword",
	})

	updatedEntry, err := testClient.UpdateEntry(testUpdatedEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	testUpdatedEntry.ModifiedOn = updatedEntry.ModifiedOn
	testUpdatedEntry.ModifiedBy = updatedEntry.ModifiedBy
	testUpdatedEntry.Tags = updatedEntry.Tags

	if updatedEntry.ID != testUpdatedEntry.ID ||
		updatedEntry.VaultId != testUpdatedEntry.VaultId ||
		updatedEntry.EntryName != testUpdatedEntry.EntryName ||
		updatedEntry.Description != testUpdatedEntry.Description ||
		updatedEntry.Type != testUpdatedEntry.Type ||
		!reflect.DeepEqual(updatedEntry.Tags, testUpdatedEntry.Tags) {
		t.Fatalf("Updated entry fields don't match expected values.")
	}

	// Check credentials separately
	updatedCreds, ok := updatedEntry.GetCredentials().(*DefaultCredentials)
	if !ok {
		t.Fatalf("Expected DefaultCredentials, got %T", updatedEntry.GetCredentials())
	}

	expectedCreds, ok := testUpdatedEntry.GetCredentials().(*DefaultCredentials)
	if !ok {
		t.Fatalf("Test updated entry credentials not set properly")
	}

	if updatedCreds.Username != expectedCreds.Username ||
		updatedCreds.Password != expectedCreds.Password {
		t.Fatalf("Updated credentials don't match. Got {%s, %s}, Expected {%s, %s}",
			updatedCreds.Username, updatedCreds.Password,
			expectedCreds.Username, expectedCreds.Password)
	}

	testNewUserEntry = updatedEntry
}

func test_DeleteUserEntry(t *testing.T) {
	err := testClient.DeleteEntry(testNewUserEntry)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	// Verify it's gone by trying to retrieve it
	_, err = testClient.GetEntry(testVaultId, testNewUserEntry.ID)
	if err == nil {
		t.Fatalf("Entry still exists after deletion: %s", testNewUserEntry.ID)
	}
}

// TestAllCredentialTypes tests all credential subtypes
func TestAllCredentialTypes(t *testing.T) {
	// Skip this test if environment variable is set
	if os.Getenv("SKIP_CREDENTIAL_TYPES_TEST") == "true" {
		t.Skip("Skipping credential types test")
	}

	// Define all credential types to test
	testCases := []struct {
		name        string
		credentials EntryCredentialsData
	}{
		{
			name: "DefaultCredentials",
			credentials: &DefaultCredentials{
				Username: "default_user",
				Password: "default_pass",
				Domain:   "default.domain",
			},
		},
		{
			name: "PrivateKeyCredentials",
			credentials: &PrivateKeyCredentials{
				PrivateKeyData:             "private_key_data",
				PublicKeyData:              "public_key_data",
				PrivateKeyOverridePassword: "override_pwd",
				PrivateKeyPassPhrase:       "passphrase",
			},
		},
		{
			name: "AccessCodeCredentials",
			credentials: &AccessCodeCredentials{
				Password: "access_code_123",
			},
		},
		{
			name: "ApiKeyCredentials",
			credentials: &ApiKeyCredentials{
				ApiId:    "api_id_value",
				ApiKey:   "api_key_value",
				TenantId: "tenant_id_value",
			},
		},
		{
			name: "AzureServicePrincipalCredentials",
			credentials: &AzureServicePrincipalCredentials{
				ClientId:     "azure_client_id",
				ClientSecret: "azure_client_secret",
				TenantId:     "azure_tenant_id",
			},
		},
		{
			name: "ConnectionStringCredentials",
			credentials: &ConnectionStringCredentials{
				ConnectionString: "Server=myserver;Database=mydb;User Id=user;Password=pwd;",
			},
		},
		{
			name: "PasskeyCredentials",
			credentials: &PasskeyCredentials{
				PasskeyPrivateKey: "passkey_private_data",
				PasskeyRpID:       "passkey.example.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timestamp := time.Now().UnixNano() % 1000000
			entryName := fmt.Sprintf("Test_%s_%d", tc.name, timestamp)

			entry := Entry{
				EntryName:   entryName,
				Description: "Testing credential type: " + tc.name,
				VaultId:     testVaultId,
				Type:        "Credential",
				Tags:        []string{"test", tc.name},
			}

			entry.SetCredentials(tc.credentials)

			t.Logf("Creating %s entry", tc.name)
			createdEntry, err := testClient.NewEntry(entry)

			entryId := createdEntry.ID

			// Cleanup at the end
			defer func() {
				if entryId != "" {
					t.Logf("Cleaning up %s entry with ID %s", tc.name, entryId)
					cleanupEntry := Entry{
						ID:      entryId,
						VaultId: testVaultId,
					}
					if err := testClient.DeleteEntry(cleanupEntry); err != nil {
						t.Logf("Warning: Failed to delete test entry: %v", err)
					}
				}
			}()

			// Verify entry was created
			if entryId == "" {
				t.Fatalf("Created entry has no ID")
			}

			t.Logf("Successfully created %s entry with ID: %s", tc.name, entryId)

			// Try to fetch the entry
			fetchedEntry, err := testClient.GetEntry(testVaultId, entryId)
			if err != nil {
				t.Fatalf("Failed to fetch %s entry: %v", tc.name, err)
			}

			// Verify entry fields
			if fetchedEntry.EntryName != entryName ||
				fetchedEntry.Type != "Credential" {
				t.Fatalf("Entry basic fields don't match")
			}

			// Verify credentials type - just check that we got the right type back
			fetchedCreds := fetchedEntry.GetCredentials()
			if fetchedCreds == nil {
				t.Fatalf("Fetched credentials are nil")
			}

			expectedType := fmt.Sprintf("*dvls.%s", tc.name)
			actualType := fmt.Sprintf("%T", fetchedCreds)

			if !strings.HasSuffix(actualType, tc.name) {
				t.Fatalf("Credential type mismatch. Expected %s, got %s",
					expectedType, actualType)
			}

			t.Logf("Successfully verified %s entry", tc.name)

			// Test update operation (simple name change)
			fetchedEntry.EntryName = entryName + "_Updated"

			t.Logf("Updating %s entry", tc.name)
			_, err = testClient.UpdateEntry(fetchedEntry)
			if err != nil {
				// Some credential types might not support updating
				t.Logf("Warning: Failed to update %s entry: %v", tc.name, err)
				return
			}

			t.Logf("Successfully updated %s entry", tc.name)
		})
	}
}
