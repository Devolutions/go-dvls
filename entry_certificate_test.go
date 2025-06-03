package dvls

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	testCertificateFilePath     string
	testCertificateEntryId      string
	testNewCertificateEntryFile EntryCertificate
	testNewCertificateEntryURL  EntryCertificate
	testCertificateEntry        EntryCertificate = EntryCertificate{
		VaultId:               testVaultId,
		Name:                  "TestK8sCertificate",
		Password:              "TestK8sCertificatePassword",
		Tags:                  []string{"test", "k8s"},
		CertificateIdentifier: "test",
	}
)

func Test_EntryCertificate(t *testing.T) {
	testCertificateFilePath = os.Getenv("TEST_CERTIFICATE_FILE_PATH")
	testCertificateEntryId = os.Getenv("TEST_CERTIFICATE_ENTRY_ID")
	testCertificateEntry.ID = testCertificateEntryId
	testCertificateEntry.VaultId = testVaultId
	location, err := time.LoadLocation("America/Montreal")
	if err != nil {
		t.Fatal(err)
	}
	expiration := time.Date(2099, 1, 1, 0, 0, 0, 0, location)
	testCertificateEntry.Expiration = expiration

	t.Run("NewCertificateFile", test_NewCertificateEntryFile)
	t.Run("NewCertificateURL", test_NewCertificateEntryURL)
	t.Run("GetEntry", test_GetCertificateEntry)
	t.Run("UpdateEntry", test_UpdateCertificateEntry)
	t.Run("DeleteEntry", test_DeleteCertificateEntry)
}

func test_GetCertificateEntry(t *testing.T) {
	testGetEntry := testCertificateEntry

	entry, err := testClient.Entries.Certificate.Get(testGetEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	entry, err = testClient.Entries.Certificate.GetPassword(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.data = testGetEntry.data

	if !entry.Expiration.Equal(testGetEntry.Expiration) {
		t.Fatalf("fetched entry expiration did not match test entry. Expected %v, got %v", testGetEntry.Expiration, entry.Expiration)
	}

	entry.Expiration = testGetEntry.Expiration

	if !reflect.DeepEqual(entry, testGetEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testGetEntry, entry)
	}
}

func test_NewCertificateEntryFile(t *testing.T) {
	entry := testCertificateEntry
	entry.ID = ""
	file, err := os.Open(testCertificateFilePath)
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal("failed read file. error: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		t.Fatal("failed read file. error: %w", err)
	}

	entry.CertificateIdentifier = stat.Name()
	entry.UseDefaultCredentials = true

	newEntry, err := testClient.Entries.Certificate.NewFile(entry, fileBytes)
	if err != nil {
		t.Fatal(err)
	}

	returnedFileBytes, err := testClient.Entries.Certificate.GetFileContent(newEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fileBytes, returnedFileBytes) {
		t.Fatalf("fetched file content did not match test file content. Expected %#v, got %#v", fileBytes, returnedFileBytes)
	}

	entry.ID = newEntry.ID
	entry.data = newEntry.data
	newEntry, err = testClient.Entries.Certificate.GetPassword(newEntry)
	if err != nil {
		t.Fatal(err)
	}

	testNewCertificateEntryFile = newEntry

	if !entry.Expiration.Equal(newEntry.Expiration) {
		t.Fatalf("fetched entry expiration did not match test entry. Expected %v, got %v", entry.Expiration, newEntry.Expiration)
	}

	entry.Expiration = newEntry.Expiration

	if !reflect.DeepEqual(entry, newEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", entry, newEntry)
	}
}

func test_NewCertificateEntryURL(t *testing.T) {
	entry := testCertificateEntry
	entry.ID = ""
	entry.CertificateIdentifier = "https://devolutions.net/"

	newEntry, err := testClient.Entries.Certificate.NewURL(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.ID = newEntry.ID
	entry.data = newEntry.data
	newEntry, err = testClient.Entries.Certificate.GetPassword(newEntry)
	if err != nil {
		t.Fatal(err)
	}

	testNewCertificateEntryURL = newEntry

	if !entry.Expiration.Equal(newEntry.Expiration) {
		t.Fatalf("fetched entry expiration did not match test entry. Expected %v, got %v", entry.Expiration, newEntry.Expiration)
	}

	entry.Expiration = newEntry.Expiration

	if !reflect.DeepEqual(entry, newEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", entry, newEntry)
	}
	testCertificateEntry = entry
}

func test_UpdateCertificateEntry(t *testing.T) {
	testUpdatedEntry := testNewCertificateEntryURL
	testUpdatedEntry.Name = "TestK8sUpdatedEntry"

	entry, err := testClient.Entries.Certificate.Update(testUpdatedEntry)
	if err != nil {
		t.Fatal(err)
	}

	entry, err = testClient.Entries.Certificate.GetPassword(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.data = testUpdatedEntry.data

	if !reflect.DeepEqual(entry, testUpdatedEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testUpdatedEntry, entry)
	}

	testNewCertificateEntryURL = entry
}

func test_DeleteCertificateEntry(t *testing.T) {
	err := testClient.Entries.Certificate.Delete(testNewCertificateEntryURL.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = testClient.Entries.Certificate.Delete(testNewCertificateEntryFile.ID)
	if err != nil {
		t.Fatal(err)
	}
}
