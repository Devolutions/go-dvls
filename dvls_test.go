package dvls

import (
	"os"
	"reflect"
	"testing"
)

const testEntryId string = "76a4fcf6-fec1-4297-bc1e-a327841055ad"
const testVaultId string = "e0f4f35d-8cb5-40d9-8b2b-35c96ea1c9b5"

var testPassword string = "TestK8sPassword"
var testEntry DvlsEntry = DvlsEntry{
	ID:                testEntryId,
	VaultId:           testVaultId,
	Description:       "Test description",
	EntryName:         "TestK8sSecret",
	ConnectionType:    ServerConnectionCredential,
	ConnectionSubType: ServerConnectionSubTypeDefault,
	Tags:              []string{"Test tag 1", "Test tag 2", "testtag"},
	Credentials: DvlsEntryCredentials{
		Username: "TestK8s",
		Password: &testPassword,
	},
}

var testClient Client

func Test_NewClient(t *testing.T) {
	c, user, err := NewClient(os.Getenv("TEST_USER"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_INSTANCE"))
	if err != nil {
		t.Fatal(err)
	}
	if user.UserType != UserAuthenticationApplication {
		t.Fatalf("user credentials is not an Application. User type %s", user.UserType)
	}

	testClient = c

	t.Run("isLogged", test_isLogged)
	t.Run("GetEntryCredentialsPassword", test_GetEntryCredentialsPassword)
	t.Run("GetEntry", test_GetEntry)
	t.Run("NewEntry", test_NewEntry)
	t.Run("GetServerInfo", test_GetServerInfo)
}

func test_GetEntryCredentialsPassword(t *testing.T) {
	testSecret := DvlsEntryCredentials{
		Username: "TestK8s",
		Password: &testPassword,
	}
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

	testGetEntry.Credentials = DvlsEntryCredentials{
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
	testNewEntry := testEntry

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
}

func test_isLogged(t *testing.T) {
	islogged, err := testClient.isLogged()
	if err != nil {
		t.Fatal(err)
	}
	if !islogged {
		t.Fatalf("expected token to be valid but isLogged returned %t", islogged)
	}

	invalidClient := testClient
	invalidClient.credential.token = "placeholder"
	islogged, err = invalidClient.isLogged()
	if err != nil {
		t.Fatal(err)
	}
	if islogged {
		t.Fatalf("expected token to be invalid but isLogged returned %t", islogged)
	}
}

func test_GetServerInfo(t *testing.T) {
	info, err := testClient.GetServerInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("server info: %#v", info)
}
