package dvls

import (
	"os"
	"reflect"
	"testing"
)

var (
	testHostEntryId  string
	testHostPassword           = "testpass123"
	testHostEntry    EntryHost = EntryHost{
		Description:    "Test host description",
		EntryName:      "TestHost",
		ConnectionType: ServerConnectionHost,
		Tags:           []string{"Test tag 1", "Test tag 2", "host"},
	}
)

const (
	testHostUsername string = "testuser"
	testHost         string = "host1234"
)

func Test_EntryHost(t *testing.T) {
	testHostEntryId = os.Getenv("TEST_HOST_ENTRY_ID")
	testHostEntry.ID = testHostEntryId
	testHostEntry.VaultId = testVaultId
	testHostEntry.HostDetails = EntryHostAuthDetails{
		Username: testHostUsername,
		Host:     testHost,
	}

	t.Run("GetEntry", test_GetHostEntry)
	t.Run("GetEntryHost", test_GetHostDetails)
}

func test_GetHostEntry(t *testing.T) {
	entry, err := testClient.Entries.Host.Get(testHostEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testHostEntry.ModifiedDate = entry.ModifiedDate
	if !reflect.DeepEqual(entry, testHostEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testHostEntry, entry)
	}
}

func test_GetHostDetails(t *testing.T) {
	entry, err := testClient.Entries.Host.Get(testHostEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	entryWithSensitiveData, err := testClient.Entries.Host.GetHostDetails(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.HostDetails.Password = entryWithSensitiveData.HostDetails.Password

	expectedDetails := testHostEntry.HostDetails

	expectedDetails.Password = &testHostPassword

	if !reflect.DeepEqual(expectedDetails, entry.HostDetails) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", expectedDetails, entry.HostDetails)
	}
}
