package dvls

import (
	"reflect"
	"testing"
)

const testEntryId string = "76a4fcf6-fec1-4297-bc1e-a327841055ad"

var (
	testNewEntry Entry
	testEntry    Entry = Entry{
		ID:                testEntryId,
		VaultId:           testVaultId,
		Description:       "Test description",
		EntryName:         "TestK8sSecret",
		ConnectionType:    ServerConnectionCredential,
		ConnectionSubType: ServerConnectionSubTypeDefault,
		Tags:              []string{"Test tag 1", "Test tag 2", "testtag"},
		Credentials:       NewEntryCredentials("TestK8s", "TestK8sPassword"),
	}
)

func Test_Entries(t *testing.T) {
	t.Run("GetEntryCredentialsPassword", test_GetEntryCredentialsPassword)
	t.Run("GetEntry", test_GetEntry)
	t.Run("NewEntry", test_NewEntry)
	t.Run("UpdateEntry", test_UpdateEntry)
	t.Run("DeleteEntry", test_DeleteEntry)
}

func test_GetEntryCredentialsPassword(t *testing.T) {
	testSecret := testEntry.Credentials
	secret, err := testClient.GetEntryCredentialsPassword(testEntry)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testSecret, secret.Credentials) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", testSecret, secret.Credentials)
	}
}

func test_GetEntry(t *testing.T) {
	testGetEntry := testEntry

	testGetEntry.Credentials = EntryCredentials{
		Username: testEntry.Credentials.Username,
	}
	entry, err := testClient.GetEntry(testGetEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testGetEntry.ModifiedDate = entry.ModifiedDate

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
	testNewEntry.ModifiedDate = entry.ModifiedDate
	testNewEntry.Tags = entry.Tags

	if !reflect.DeepEqual(entry, testNewEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testNewEntry, entry)
	}

	testNewEntry = entry
}

func test_UpdateEntry(t *testing.T) {
	testUpdatedEntry := testNewEntry
	testUpdatedEntry.EntryName = "TestK8sUpdatedEntry"
	testUpdatedEntry.Credentials = NewEntryCredentials("TestK8sUpdatedUser", "TestK8sUpdatedPassword")

	entry, err := testClient.UpdateEntry(testUpdatedEntry)
	if err != nil {
		t.Fatal(err)
	}

	testUpdatedEntry.ModifiedDate = entry.ModifiedDate
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
