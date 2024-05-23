package dvls

import (
	"reflect"
	"testing"
)

var (
	testNewEntry EntryUserCredential
	testEntry    EntryUserCredential = EntryUserCredential{
		Description:       "Test description",
		EntryName:         "TestK8sSecret",
		ConnectionType:    ServerConnectionCredential,
		ConnectionSubType: ServerConnectionSubTypeDefault,
		Tags:              []string{"Test tag 1", "Test tag 2", "testtag"},
	}
)

const (
	testEntryUsername string = "TestK8s"
	testEntryPassword string = "TestK8sPassword"
)

func Test_Entries(t *testing.T) {
	testEntry.ID = testEntryId
	testEntry.VaultId = testVaultId
	testEntry.Credentials = testClient.Entries.UserCredential.NewUserAuthDetails(testEntryUsername, testEntryPassword)

	t.Run("GetEntry", test_GetEntry)
	t.Run("GetEntryCredentialsPassword", test_GetEntryCredentialsPassword)
	t.Run("NewEntry", test_NewEntry)
	t.Run("UpdateEntry", test_UpdateEntry)
	t.Run("DeleteEntry", test_DeleteEntry)
}

func test_GetEntry(t *testing.T) {
	testGetEntry := testEntry

	testGetEntry.Credentials = EntryUserAuthDetails{
		Username: testEntry.Credentials.Username,
	}
	entry, err := testClient.Entries.UserCredential.Get(testGetEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testClient.Entries.UserCredential.Get(testGetEntry.ID)
	testGetEntry.ModifiedDate = entry.ModifiedDate

	if !reflect.DeepEqual(entry, testGetEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testGetEntry, entry)
	}
}

func test_GetEntryCredentialsPassword(t *testing.T) {
	testSecret := testEntry.Credentials
	secret, err := testClient.Entries.UserCredential.GetUserAuthDetails(testEntry)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testSecret, secret.Credentials) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", testSecret, secret.Credentials)
	}
}

func test_NewEntry(t *testing.T) {
	testNewEntry = testEntry

	testNewEntry.EntryName = "TestK8sNewEntry"

	entry, err := testClient.Entries.UserCredential.New(testNewEntry)
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
	testUpdatedEntry.Credentials = testClient.Entries.UserCredential.NewUserAuthDetails("TestK8sUpdatedUser", "TestK8sUpdatedPassword")

	entry, err := testClient.Entries.UserCredential.Update(testUpdatedEntry)
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
	err := testClient.Entries.UserCredential.Delete(testNewEntry.ID)
	if err != nil {
		t.Fatal(err)
	}
}
