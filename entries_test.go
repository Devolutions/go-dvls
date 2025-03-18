package dvls

import (
	"os"
	"reflect"
	"testing"
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

const (
	testEntryUsername string = "TestK8s"
	testEntryPassword string = "TestK8sPassword"
)

func Test_EntryUserCredentials(t *testing.T) {
	testUserEntryId = os.Getenv("TEST_USER_ENTRY_ID")
	testUserEntry.ID = testUserEntryId
	testUserEntry.VaultId = testVaultId

	t.Run("GetEntry", test_GetUserEntry)
	t.Run("NewEntry", test_NewUserEntry)
	t.Run("UpdateEntry", test_UpdateUserEntry)
	t.Run("DeleteEntry", test_DeleteUserEntry)
}

func test_GetUserEntry(t *testing.T) {
	entry, err := testClient.GetEntry(testVaultId, testUserEntry.ID)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if entry.ID != testUserEntry.ID ||
		entry.VaultId != testUserEntry.VaultId ||
		entry.EntryName != testUserEntry.EntryName ||
		entry.Description != testUserEntry.Description ||
		entry.Type != testUserEntry.Type ||
		!reflect.DeepEqual(entry.Tags, testUserEntry.Tags) ||
		entry.Credentials.Username != "TestK8s" ||
		entry.Credentials.Password != "TestK8sPassword" {

		t.Fatalf("Entry fields don't match expected values. \nGot: %+v\nExpected ID: %s, VaultId: %s, EntryName: %s, Description: %s, Type: %s, Tags: %v, Credentials: {Username: TestK8s, Password: TestK8sPassword}",
			entry,
			testUserEntry.ID,
			testUserEntry.VaultId,
			testUserEntry.EntryName,
			testUserEntry.Description,
			testUserEntry.Type,
			testUserEntry.Tags)
	}
}

func test_NewUserEntry(t *testing.T) {
	testNewUserEntry = testUserEntry
	testNewUserEntry.ID = "" // Set empty ID for new creation
	testNewUserEntry.EntryName = "TestK8sNewEntry"
	testNewUserEntry.Credentials = EntryCredentials{
		Username: "TestK8sNew",
		Password: "TestK8sNewPassword",
	}

	newEntry, err := testClient.NewEntry(testNewUserEntry)
	if err != nil {
		t.Fatalf("Failed to create new entry: %v", err)
	}

	testNewUserEntry.ID = newEntry.ID
	testNewUserEntry.ModifiedOn = newEntry.ModifiedOn
	testNewUserEntry.ModifiedBy = newEntry.ModifiedBy
	testNewUserEntry.CreatedOn = newEntry.CreatedOn
	testNewUserEntry.CreatedBy = newEntry.CreatedBy
}

func test_UpdateUserEntry(t *testing.T) {
	testUpdatedEntry := testNewUserEntry
	testUpdatedEntry.EntryName = "TestK8sUpdatedEntry"
	testUpdatedEntry.Credentials = EntryCredentials{
		Username: "TestK8sUpdatedUser",
		Password: "TestK8sUpdatedPassword",
	}

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
		!reflect.DeepEqual(updatedEntry.Tags, testUpdatedEntry.Tags) ||
		updatedEntry.Credentials.Username != testUpdatedEntry.Credentials.Username ||
		updatedEntry.Credentials.Password != testUpdatedEntry.Credentials.Password {

		t.Fatalf("Updated entry fields don't match expected values. \nGot: %+v\nExpected ID: %s, VaultId: %s, EntryName: %s, Description: %s, Type: %s, Tags: %v, Credentials: {Username: %s, Password: %s}",
			updatedEntry,
			testUpdatedEntry.ID,
			testUpdatedEntry.VaultId,
			testUpdatedEntry.EntryName,
			testUpdatedEntry.Description,
			testUpdatedEntry.Type,
			testUpdatedEntry.Tags,
			testUpdatedEntry.Credentials.Username,
			testUpdatedEntry.Credentials.Password)
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
