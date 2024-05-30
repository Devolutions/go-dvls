package dvls

import (
	"os"
	"reflect"
	"testing"
)

var (
	testUserEntryId  string
	testNewUserEntry EntryUserCredential
	testUserEntry    EntryUserCredential = EntryUserCredential{
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

func Test_EntryUserCredentials(t *testing.T) {
	testUserEntryId = os.Getenv("TEST_USER_ENTRY_ID")
	testUserEntry.ID = testUserEntryId
	testUserEntry.VaultId = testVaultId
	testUserEntry.Credentials = testClient.Entries.UserCredential.NewUserAuthDetails(testEntryUsername, testEntryPassword)

	t.Run("GetEntry", test_GetUserEntry)
	t.Run("NewEntry", test_NewUserEntry)
	t.Run("GetEntryCredentialsPassword", test_GetEntryCredentialsPassword)

	t.Run("UpdateEntry", test_UpdateUserEntry)
	t.Run("DeleteEntry", test_DeleteUserEntry)
}

func test_GetUserEntry(t *testing.T) {
	testGetEntry := testUserEntry

	testGetEntry.Credentials = EntryUserAuthDetails{
		Username: testUserEntry.Credentials.Username,
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
	testSecret := testUserEntry.Credentials
	secret, err := testClient.Entries.UserCredential.GetUserAuthDetails(testUserEntry)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(testSecret, secret.Credentials) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", testSecret, secret.Credentials)
	}
}

func test_NewUserEntry(t *testing.T) {
	testNewUserEntry = testUserEntry

	testNewUserEntry.EntryName = "TestK8sNewEntry"

	entry, err := testClient.Entries.UserCredential.New(testNewUserEntry)
	if err != nil {
		t.Fatal(err)
	}

	testNewUserEntry.ID = entry.ID
	testNewUserEntry.ModifiedDate = entry.ModifiedDate
	testNewUserEntry.Tags = entry.Tags

	if !reflect.DeepEqual(entry, testNewUserEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testNewUserEntry, entry)
	}

	testNewUserEntry = entry
}

func test_UpdateUserEntry(t *testing.T) {
	testUpdatedEntry := testNewUserEntry
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

	testNewUserEntry = entry
}

func test_DeleteUserEntry(t *testing.T) {
	err := testClient.Entries.UserCredential.Delete(testNewUserEntry.ID)
	if err != nil {
		t.Fatal(err)
	}
}
