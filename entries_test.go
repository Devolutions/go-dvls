package dvls

import (
	"reflect"
	"testing"
)

var (
	testNewEntry Entry
	testEntry    Entry = Entry{
		Description: "Test description",
		EntryName:   "TestK8sSecret",
		Type:        "Credential",
		Tags:        []string{"Test tag 1", "Test tag 2", "testtag"},
		Credentials: EntryCredentials{"TestK8s", "TestK8sPassword"},
	}
)

func Test_Entries(t *testing.T) {
	testEntry.ID = testEntryId
	testEntry.VaultId = testVaultId

	t.Run("GetEntry", test_GetEntry)
	t.Run("NewEntry", test_NewEntry)
	t.Run("UpdateEntry", test_UpdateEntry)
	t.Run("DeleteEntry", test_DeleteEntry)
}

func test_GetEntry(t *testing.T) {
	testGetEntry := testEntry

	entry, err := testClient.GetEntry(testVaultId, testGetEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testGetEntry.ModifiedOn = entry.ModifiedOn
	testGetEntry.ModifiedBy = entry.ModifiedBy
	testGetEntry.CreatedOn = entry.CreatedOn
	testGetEntry.CreatedBy = entry.CreatedBy

	if !reflect.DeepEqual(entry, testGetEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testGetEntry, entry)
	}
}

func test_NewEntry(t *testing.T) {
	testNewEntry = testEntry

	testNewEntry.EntryName = "TestK8sNewEntry"

	entry, err := testClient.NewEntry(testNewEntry)
	if err != nil {
		t.Fatal(err)
	}

	testNewEntry.ID = entry.ID
	testNewEntry.ModifiedOn = entry.ModifiedOn
	testNewEntry.Tags = entry.Tags

	if !reflect.DeepEqual(entry, testNewEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testNewEntry, entry)
	}

	testNewEntry = entry
}

func test_UpdateEntry(t *testing.T) {
	testUpdatedEntry := testNewEntry
	testUpdatedEntry.EntryName = "TestK8sUpdatedEntry"
	testUpdatedEntry.Credentials = EntryCredentials{"TestK8sUpdatedUser", "TestK8sUpdatedPassword"}

	entry, err := testClient.UpdateEntry(testUpdatedEntry)
	if err != nil {
		t.Fatal(err)
	}

	testUpdatedEntry.ModifiedOn = entry.ModifiedOn
	testUpdatedEntry.Tags = entry.Tags

	if !reflect.DeepEqual(entry, testUpdatedEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testUpdatedEntry, entry)
	}

	testNewEntry = entry
}

func test_DeleteEntry(t *testing.T) {
	err := testClient.DeleteEntry(testNewEntry.ID)
	if err != nil {
		t.Fatal(err)
	}
}
