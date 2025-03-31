package dvls

import (
	"reflect"
	"testing"
)

var (
	testUserEntry = EntryUserCredential{
		ID:          "",
		VaultId:     testVaultId,
		EntryName:   "TestGoDvlsSecret",
		Description: "Test description",
		Type:        "Credential",
		SubType:     "Default",
		Tags:        []string{"testtag"},
		Credentials: EntryCredentials{
			Username: testEntryUsername,
			Password: testEntryPassword,
		},
	}
)

const (
	testEntryUsername string = "Test"
	testEntryPassword string = "TestPassword"
)

func Test_EntryUserCredentials(t *testing.T) {
	t.Run("NewEntry", test_NewUserEntry)
	t.Run("GetEntry", test_GetUserEntry)
	t.Run("UpdateEntry", test_UpdateUserEntry)
	t.Run("DeleteEntry", test_DeleteUserEntry)
}

// Allowing accurate comparison by ignoring fields that differ due to server assignment.
func NormalizeEntry(source, target *EntryUserCredential) {
	target.ID = source.ID
	target.ModifiedOn = source.ModifiedOn
	target.ModifiedBy = source.ModifiedBy
	target.CreatedOn = source.CreatedOn
	target.CreatedBy = source.CreatedBy
}

func test_NewUserEntry(t *testing.T) {
	testUserEntry.VaultId = testVaultId
	newEntry, err := testClient.Entries.UserCredential.New(testUserEntry)
	if err != nil {
		t.Fatalf("Failed to create new entry: %v", err)
	}

	NormalizeEntry(&newEntry, &testUserEntry)

	if !reflect.DeepEqual(&newEntry, &testUserEntry) {
		t.Fatalf("Entries differ.\nGot:      %+v\nExpected: %+v", &newEntry, &testUserEntry)
	}
}

func test_GetUserEntry(t *testing.T) {
	entry, err := testClient.Entries.UserCredential.Get(testVaultId, testUserEntry.ID)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	NormalizeEntry(&entry, &testUserEntry)

	if !reflect.DeepEqual(&entry, &testUserEntry) {
		t.Fatalf("Entries differ.\nGot:      %+v\nExpected: %+v", &entry, &testUserEntry)
	}
}

func test_UpdateUserEntry(t *testing.T) {
	testUpdatedEntry := testUserEntry
	testUpdatedEntry.EntryName = "TestGoDvlsSecretUpdated"
	testUpdatedEntry.Description = "Test description updated"
	testUpdatedEntry.Credentials = EntryCredentials{
		Username: "TestK8sUpdatedUser",
		Password: "TestK8sUpdatedPassword",
	}

	updatedEntry, err := testClient.Entries.UserCredential.Update(testUpdatedEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}
	NormalizeEntry(&updatedEntry, &testUpdatedEntry)

	if !reflect.DeepEqual(&updatedEntry, &testUpdatedEntry) {
		t.Fatalf("Entries differ.\nGot:      %+v\nExpected: %+v", &updatedEntry, &testUpdatedEntry)
	}
	testUserEntry = updatedEntry
}

func test_DeleteUserEntry(t *testing.T) {
	err := testClient.Entries.UserCredential.Delete(testUserEntry)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	// Verify it's gone by trying to retrieve it
	_, err = testClient.Entries.UserCredential.Get(testVaultId, testUserEntry.ID)
	if err == nil {
		t.Fatalf("Entry still exists after deletion: %s", testUserEntry.ID)
	}
}
