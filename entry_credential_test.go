package dvls

import (
	"testing"
)

var (
	testCredentialDefaultEntryID *string
	testCredentialDefaultEntry   *Entry

	testCredentialAccessCodeEntryID *string
	testCredentialAccessCodeEntry   *Entry
)

func Test_EntryUserCredentials(t *testing.T) {
	if !t.Run("NewEntry", test_NewUserEntry) {
		t.Skip("Skipping subsequent tests due to failure in NewEntry")
		return
	}

	if !t.Run("GetEntry", test_GetUserEntry) {
		t.Skip("Skipping subsequent tests due to failure in GetEntry")
		return
	}

	if !t.Run("UpdateEntry", test_UpdateUserEntry) {
		t.Skip("Skipping subsequent tests due to failure in UpdateEntry")
		return
	}

	if !t.Run("DeleteEntry", test_DeleteUserEntry) {
		t.Skip("Skipping subsequent tests due to failure in DeleteEntry")
		return
	}
}

func test_NewUserEntry(t *testing.T) {
	// Credential/Default
	testCredentialDefaultEntry := Entry{
		VaultId:     testVaultId,
		Name:        "TestGoDvlsUsernamePassword",
		Path:        "go-dvls\\usernamepassword",
		Description: "Test Username/Password entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeDefault,
		Data: EntryCredentialDefaultData{
			Domain:   "www.example.com",
			Password: "abc-123",
			Username: "john.doe",
		},
		Tags: []string{"tag1", "tag2"},
	}

	newCredentialDefaultEntryID, err := testClient.Entries.Credential.New(testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to create new entry: %v", err)
	}

	if newCredentialDefaultEntryID == "" {
		t.Fatal("New entry ID is empty after creation.")
	}

	testCredentialDefaultEntryID = &newCredentialDefaultEntryID

	// Credential/AccessCode
	testCredentialAccessCodeEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsAccessCode",
		Path:        "go-dvls\\accesscode",
		Description: "Test Secret entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeAccessCode,
		Data: EntryCredentialAccessCodeData{
			Password: "abc-123",
		},
		Tags: []string{"testtag"},
	}

	newCredentialAccessCodeEntryID, err := testClient.Entries.Credential.New(testCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to create new entry: %v", err)
	}

	if newCredentialAccessCodeEntryID == "" {
		t.Fatal("New entry ID is empty after creation.")
	}

	testCredentialAccessCodeEntryID = &newCredentialAccessCodeEntryID
}

func test_GetUserEntry(t *testing.T) {
	// Credential/Default
	credentialDefaultEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialDefaultEntryID)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if credentialDefaultEntry.ID == "" {
		t.Fatalf("Entry ID is empty after GET: %v", credentialDefaultEntry)
	}

	testCredentialDefaultEntry = &credentialDefaultEntry

	// Credential/AccessCode
	credentialAccessCodeEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialAccessCodeEntryID)
	if err != nil {
		t.Fatalf("Failed to get entry: %v", err)
	}

	if credentialAccessCodeEntry.ID == "" {
		t.Fatalf("Entry ID is empty after GET: %v", credentialAccessCodeEntry)
	}

	testCredentialAccessCodeEntry = &credentialAccessCodeEntry
}

func test_UpdateUserEntry(t *testing.T) {
	// Credential/Default
	updatedCredentialDefaultEntry := *testCredentialDefaultEntry
	updatedCredentialDefaultEntry.Name = updatedCredentialDefaultEntry.Name + "Updated"
	updatedCredentialDefaultEntry.Path = updatedCredentialDefaultEntry.Path + "\\updated"
	updatedCredentialDefaultEntry.Description = updatedCredentialDefaultEntry.Description + " updated"
	updatedCredentialDefaultEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedData, ok := updatedCredentialDefaultEntry.GetCredentialDefaultData()
	if !ok {
		t.Fatalf("Failed to get credential default data from entry: %v", updatedCredentialDefaultEntry)
	}
	updatedData.Password = updatedData.Password + "-updated"
	updatedCredentialDefaultEntry.Data = updatedData

	updatedCredentialDefaultEntry, err := testClient.Entries.Credential.Update(updatedCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Credential/AccessCode
	updatedCredentialAccessCodeEntry := *testCredentialAccessCodeEntry
	updatedCredentialAccessCodeEntry.Name = updatedCredentialAccessCodeEntry.Name + "Updated"
	updatedCredentialAccessCodeEntry.Path = updatedCredentialAccessCodeEntry.Path + "\\updated"
	updatedCredentialAccessCodeEntry.Description = updatedCredentialAccessCodeEntry.Description + " updated"
	updatedCredentialAccessCodeEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedDataAccessCode, ok := updatedCredentialAccessCodeEntry.GetCredentialAccessCodeData()
	if !ok {
		t.Fatalf("Failed to get credential access code data from entry: %v", updatedCredentialAccessCodeEntry)
	}
	updatedDataAccessCode.Password = updatedDataAccessCode.Password + "-updated"
	updatedCredentialAccessCodeEntry.Data = updatedDataAccessCode

	updatedCredentialAccessCodeEntry, err = testClient.Entries.Credential.Update(updatedCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}
}

func test_DeleteUserEntry(t *testing.T) {
	// Credential/Default
	err := testClient.Entries.Credential.Delete(*testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialDefaultEntry)
	if err == nil {
		t.Fatalf("Entry still exists after deletion: %s", *testCredentialDefaultEntryID)
	}

	// Credential/AccessCode
	err = testClient.Entries.Credential.Delete(*testCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to delete entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialAccessCodeEntry)
	if err == nil {
		t.Fatalf("Entry still exists after deletion: %s", *testCredentialAccessCodeEntryID)
	}
}
