package dvls

import (
	"os"
	"reflect"
	"testing"
)

var (
	testWebsiteEntryId string
	testWebsiteEntry   EntryWebsite = EntryWebsite{
		Description:       "Test website description",
		EntryName:         "TestWebsite",
		ConnectionType:    ServerConnectionWebBrowser,
		ConnectionSubType: ServerConnectionSubTypeGoogleChrome,
		Tags:              []string{"Test tag 1", "Test tag 2", "web"},
	}
)

const (
	testWebsiteUsername string = "testuser"
	testWebsiteURL      string = "https://test.example.com"
	testWebsiteBrowser  string = "GoogleChrome"
)

var testWebsitePassword = "testpass123"

func Test_EntryWebsite(t *testing.T) {
	testWebsiteEntryId = os.Getenv("TEST_WEBSITE_ENTRY_ID")
	testWebsiteEntry.ID = testWebsiteEntryId
	testWebsiteEntry.VaultId = testVaultId
	testWebsiteEntry.WebsiteDetails = EntryWebsiteAuthDetails{
		Username:              testWebsiteUsername,
		URL:                   testWebsiteURL,
		WebBrowserApplication: 3,
	}
	testWebsiteEntry.ConnectionSubType = ServerConnectionSubTypeGoogleChrome

	t.Run("GetEntry", test_GetWebsiteEntry)
	t.Run("GetEntryWebsite", test_GetWebsiteDetails)
}

func test_GetWebsiteEntry(t *testing.T) {
	entry, err := testClient.Entries.Website.Get(testWebsiteEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testWebsiteEntry.ModifiedDate = entry.ModifiedDate
	if !reflect.DeepEqual(entry, testWebsiteEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testWebsiteEntry, entry)
	}
}

func test_GetWebsiteDetails(t *testing.T) {
	entry, err := testClient.Entries.Website.Get(testWebsiteEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	entryWithSensitiveData, err := testClient.Entries.Website.GetWebsiteDetails(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.WebsiteDetails.Password = entryWithSensitiveData.WebsiteDetails.Password

	expectedDetails := testWebsiteEntry.WebsiteDetails

	expectedDetails.Password = &testWebsitePassword

	if !reflect.DeepEqual(expectedDetails, entry.WebsiteDetails) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", expectedDetails, entry.WebsiteDetails)
	}
}
