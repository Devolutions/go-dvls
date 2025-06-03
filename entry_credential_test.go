package dvls

import (
	"testing"
)

var (
	testCredentialAccessCodeEntryID *string
	testCredentialAccessCodeEntry   *Entry

	testCredentialApiKeyEntryID *string
	testCredentialApiKeyEntry   *Entry

	testCredentialAzureServicePrincipalEntryID *string
	testCredentialAzureServicePrincipalEntry   *Entry

	testCredentialConnectionStringEntryID *string
	testCredentialConnectionStringEntry   *Entry

	testCredentialDefaultEntryID *string
	testCredentialDefaultEntry   *Entry

	testCredentialPrivateKeyEntryID *string
	testCredentialPrivateKeyEntry   *Entry
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
	// Notes: all entries values are random and for testing purposes only.

	// Credential/AccessCode
	testCredentialAccessCodeEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsAccessCode",
		Path:        "go-dvls\\accesscode",
		Description: "Test AccessCode entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeAccessCode,
		Data: EntryCredentialAccessCodeData{
			Password: "abc-123",
		},
		Tags: []string{"accesscode"},
	}

	newCredentialAccessCodeEntryID, err := testClient.Entries.Credential.New(testCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to create new AccessCode entry: %v", err)
	}

	if newCredentialAccessCodeEntryID == "" {
		t.Fatal("New AccessCode entry ID is empty after creation.")
	}

	testCredentialAccessCodeEntryID = &newCredentialAccessCodeEntryID

	// Credential/ApiKey
	testCredentialApiKeyEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsApiKey",
		Path:        "go-dvls\\apikey",
		Description: "Test ApiKey entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeApiKey,
		Data: EntryCredentialApiKeyData{
			ApiID:    "abcd1234-abcd-1234-abcd-1234abcd1234",
			ApiKey:   "123-abc",
			TenantID: "00000000-aaaa-bbbb-cccc-000000000000",
		},
		Tags: []string{"apikey"},
	}

	newCredentialApiKeyEntryID, err := testClient.Entries.Credential.New(testCredentialApiKeyEntry)
	if err != nil {
		t.Fatalf("Failed to create new ApiKey entry: %v", err)
	}

	if newCredentialApiKeyEntryID == "" {
		t.Fatal("New ApiKey entry ID is empty after creation.")
	}

	testCredentialApiKeyEntryID = &newCredentialApiKeyEntryID

	// Credential/AzureServicePrincipal
	testCredentialAzureServicePrincipalEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsAzureServicePrincipal",
		Path:        "go-dvls\\azureserviceprincipal",
		Description: "Test AzureServicePrincipal entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeAzureServicePrincipal,
		Data: EntryCredentialAzureServicePrincipalData{
			ClientID:     "abcd1234-abcd-1234-abcd-1234abcd1234",
			ClientSecret: "123-abc",
			TenantID:     "00000000-aaaa-bbbb-cccc-000000000000",
		},
		Tags: []string{"azureserviceprincipal"},
	}

	newCredentialAzureServicePrincipalEntryID, err := testClient.Entries.Credential.New(testCredentialAzureServicePrincipalEntry)
	if err != nil {
		t.Fatalf("Failed to create new AzureServicePrincipal entry: %v", err)
	}

	if newCredentialAzureServicePrincipalEntryID == "" {
		t.Fatal("New AzureServicePrincipal entry ID is empty after creation.")
	}

	testCredentialAzureServicePrincipalEntryID = &newCredentialAzureServicePrincipalEntryID

	// Credential/ConnectionString
	testCredentialConnectionStringEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsConnectionString",
		Path:        "go-dvls\\connectionstring",
		Description: "Test ConnectionString entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeConnectionString,
		Data: EntryCredentialConnectionStringData{
			ConnectionString: "Server=tcp:example.database.windows.net,1433;Initial Catalog=exampledb;Persist Security Info=False;User ID=exampleuser;Password=examplepassword;",
		},
		Tags: []string{"connectionstring"},
	}

	newCredentialConnectionStringEntryID, err := testClient.Entries.Credential.New(testCredentialConnectionStringEntry)
	if err != nil {
		t.Fatalf("Failed to create new ConnectionString entry: %v", err)
	}

	if newCredentialConnectionStringEntryID == "" {
		t.Fatal("New ConnectionString entry ID is empty after creation.")
	}

	testCredentialConnectionStringEntryID = &newCredentialConnectionStringEntryID

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
		Tags: []string{"usernamepassword"},
	}

	newCredentialDefaultEntryID, err := testClient.Entries.Credential.New(testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to create new Default entry: %v", err)
	}

	if newCredentialDefaultEntryID == "" {
		t.Fatal("New Default entry ID is empty after creation.")
	}

	testCredentialDefaultEntryID = &newCredentialDefaultEntryID

	// Credential/PrivateKey
	testCredentialPrivateKeyEntry := Entry{
		ID:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsPrivateKey",
		Path:        "go-dvls\\privatekey",
		Description: "Test Secret entry",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypePrivateKey,
		Data: EntryCredentialPrivateKeyData{
			PrivateKey:       "-----BEGIN PRIVATE KEY-----\abcdefghijklmnopqrstuvwxyz1234567890...\n-----END PRIVATE",
			PublicKey:        "-----BEGIN PUBLIC KEY-----\abcdefghijklmnopqrstuvwxyz...\n-----END PUBLIC KEY-----",
			OverridePassword: "override-password",
			Passphrase:       "passphrase",
		},
		Tags: []string{"testtag"},
	}

	newCredentialPrivateKeyEntryID, err := testClient.Entries.Credential.New(testCredentialPrivateKeyEntry)
	if err != nil {
		t.Fatalf("Failed to create new PrivateKey entry: %v", err)
	}

	if newCredentialPrivateKeyEntryID == "" {
		t.Fatal("New PrivateKey entry ID is empty after creation.")
	}

	testCredentialPrivateKeyEntryID = &newCredentialPrivateKeyEntryID
}

func test_GetUserEntry(t *testing.T) {
	// Credential/AccessCode
	credentialAccessCodeEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialAccessCodeEntryID)
	if err != nil {
		t.Fatalf("Failed to get AccessCode entry: %v", err)
	}

	if credentialAccessCodeEntry.ID == "" {
		t.Fatalf("AccessCode entry ID is empty after GET: %v", credentialAccessCodeEntry)
	}

	testCredentialAccessCodeEntry = &credentialAccessCodeEntry

	// Credential/ApiKey
	credentialApiKeyEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialApiKeyEntryID)
	if err != nil {
		t.Fatalf("Failed to get ApiKey entry: %v", err)
	}

	if credentialApiKeyEntry.ID == "" {
		t.Fatalf("ApiKey entry ID is empty after GET: %v", credentialApiKeyEntry)
	}

	testCredentialApiKeyEntry = &credentialApiKeyEntry

	// Credential/AzureServicePrincipal
	credentialAzureServicePrincipalEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialAzureServicePrincipalEntryID)
	if err != nil {
		t.Fatalf("Failed to get AzureServicePrincipal entry: %v", err)
	}

	if credentialAzureServicePrincipalEntry.ID == "" {
		t.Fatalf("AzureServicePrincipal entry ID is empty after GET: %v", credentialAzureServicePrincipalEntry)
	}

	testCredentialAzureServicePrincipalEntry = &credentialAzureServicePrincipalEntry

	// Credential/ConnectionString
	credentialConnectionStringEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialConnectionStringEntryID)
	if err != nil {
		t.Fatalf("Failed to get ConnectionString entry: %v", err)
	}

	if credentialConnectionStringEntry.ID == "" {
		t.Fatalf("ConnectionString entry ID is empty after GET: %v", credentialConnectionStringEntry)
	}

	testCredentialConnectionStringEntry = &credentialConnectionStringEntry

	// Credential/Default
	credentialDefaultEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialDefaultEntryID)
	if err != nil {
		t.Fatalf("Failed to get Default entry: %v", err)
	}

	if credentialDefaultEntry.ID == "" {
		t.Fatalf("Default entry ID is empty after GET: %v", credentialDefaultEntry)
	}

	testCredentialDefaultEntry = &credentialDefaultEntry

	// Credential/PrivateKey
	credentialPrivateKeyEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialPrivateKeyEntryID)
	if err != nil {
		t.Fatalf("Failed to get PrivateKey entry: %v", err)
	}

	if credentialPrivateKeyEntry.ID == "" {
		t.Fatalf("PrivateKey entry ID is empty after GET: %v", credentialPrivateKeyEntry)
	}

	testCredentialPrivateKeyEntry = &credentialPrivateKeyEntry
}

func test_UpdateUserEntry(t *testing.T) {
	// Credential/AccessCode
	updatedCredentialAccessCodeEntry := *testCredentialAccessCodeEntry
	updatedCredentialAccessCodeEntry.Name = updatedCredentialAccessCodeEntry.Name + "Updated"
	updatedCredentialAccessCodeEntry.Path = updatedCredentialAccessCodeEntry.Path + "\\updated"
	updatedCredentialAccessCodeEntry.Description = updatedCredentialAccessCodeEntry.Description + " updated"
	updatedCredentialAccessCodeEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedAccessCodeData, ok := updatedCredentialAccessCodeEntry.GetCredentialAccessCodeData()
	if !ok {
		t.Fatalf("Failed to get credential AccessCode data from entry: %v", updatedCredentialAccessCodeEntry)
	}
	updatedAccessCodeData.Password = updatedAccessCodeData.Password + "-updated"
	updatedCredentialAccessCodeEntry.Data = updatedAccessCodeData

	updatedCredentialAccessCodeEntry, err := testClient.Entries.Credential.Update(updatedCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to update AccessCode entry: %v", err)
	}

	// Credential/ApiKey
	updatedCredentialApiKeyEntry := *testCredentialApiKeyEntry
	updatedCredentialApiKeyEntry.Name = updatedCredentialApiKeyEntry.Name + "Updated"
	updatedCredentialApiKeyEntry.Path = updatedCredentialApiKeyEntry.Path + "\\updated"
	updatedCredentialApiKeyEntry.Description = updatedCredentialApiKeyEntry.Description + " updated"
	updatedCredentialApiKeyEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedApiKeyData, ok := updatedCredentialApiKeyEntry.GetCredentialApiKeyData()
	if !ok {
		t.Fatalf("Failed to get credential ApiKey data from entry: %v", updatedCredentialApiKeyEntry)
	}

	updatedApiKeyData.ApiKey = updatedApiKeyData.ApiKey + "-updated"
	updatedCredentialApiKeyEntry.Data = updatedApiKeyData

	updatedCredentialApiKeyEntry, err = testClient.Entries.Credential.Update(updatedCredentialApiKeyEntry)
	if err != nil {
		t.Fatalf("Failed to update ApiKey entry: %v", err)
	}

	// Credential/AzureServicePrincipal
	updatedCredentialAzureServicePrincipalEntry := *testCredentialAzureServicePrincipalEntry
	updatedCredentialAzureServicePrincipalEntry.Name = updatedCredentialAzureServicePrincipalEntry.Name + "Updated"
	updatedCredentialAzureServicePrincipalEntry.Path = updatedCredentialAzureServicePrincipalEntry.Path + "\\updated"
	updatedCredentialAzureServicePrincipalEntry.Description = updatedCredentialAzureServicePrincipalEntry.Description + " updated"
	updatedCredentialAzureServicePrincipalEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedAzureServicePrincipalData, ok := updatedCredentialAzureServicePrincipalEntry.GetCredentialAzureServicePrincipalData()
	if !ok {
		t.Fatalf("Failed to get credential AzureServicePrincipal data from entry: %v", updatedCredentialAzureServicePrincipalEntry)
	}

	updatedAzureServicePrincipalData.ClientSecret = updatedAzureServicePrincipalData.ClientSecret + "-updated"
	updatedCredentialAzureServicePrincipalEntry.Data = updatedAzureServicePrincipalData

	updatedCredentialAzureServicePrincipalEntry, err = testClient.Entries.Credential.Update(updatedCredentialAzureServicePrincipalEntry)
	if err != nil {
		t.Fatalf("Failed to update AzureServicePrincipal entry: %v", err)
	}

	// Credential/ConnectionString
	updatedCredentialConnectionStringEntry := *testCredentialConnectionStringEntry
	updatedCredentialConnectionStringEntry.Name = updatedCredentialConnectionStringEntry.Name + "Updated"
	updatedCredentialConnectionStringEntry.Path = updatedCredentialConnectionStringEntry.Path + "\\updated"
	updatedCredentialConnectionStringEntry.Description = updatedCredentialConnectionStringEntry.Description + " updated"
	updatedCredentialConnectionStringEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedConnectionStringData, ok := updatedCredentialConnectionStringEntry.GetCredentialConnectionStringData()
	if !ok {
		t.Fatalf("Failed to get credential ConnectionString data from entry: %v", updatedCredentialConnectionStringEntry)
	}

	updatedConnectionStringData.ConnectionString = updatedConnectionStringData.ConnectionString + "MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;"
	updatedCredentialConnectionStringEntry.Data = updatedConnectionStringData

	updatedCredentialConnectionStringEntry, err = testClient.Entries.Credential.Update(updatedCredentialConnectionStringEntry)
	if err != nil {
		t.Fatalf("Failed to update ConnectionString entry: %v", err)
	}

	// Credential/Default
	updatedCredentialDefaultEntry := *testCredentialDefaultEntry
	updatedCredentialDefaultEntry.Name = updatedCredentialDefaultEntry.Name + "Updated"
	updatedCredentialDefaultEntry.Path = updatedCredentialDefaultEntry.Path + "\\updated"
	updatedCredentialDefaultEntry.Description = updatedCredentialDefaultEntry.Description + " updated"
	updatedCredentialDefaultEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedDefaultData, ok := updatedCredentialDefaultEntry.GetCredentialDefaultData()
	if !ok {
		t.Fatalf("Failed to get credential default data from entry: %v", updatedCredentialDefaultEntry)
	}
	updatedDefaultData.Password = updatedDefaultData.Password + "-updated"
	updatedCredentialDefaultEntry.Data = updatedDefaultData

	updatedCredentialDefaultEntry, err = testClient.Entries.Credential.Update(updatedCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}

	// Credential/PrivateKey
	updatedCredentialPrivateKeyEntry := *testCredentialPrivateKeyEntry
	updatedCredentialPrivateKeyEntry.Name = updatedCredentialPrivateKeyEntry.Name + "Updated"
	updatedCredentialPrivateKeyEntry.Path = updatedCredentialPrivateKeyEntry.Path + "\\updated"
	updatedCredentialPrivateKeyEntry.Description = updatedCredentialPrivateKeyEntry.Description + " updated"
	updatedCredentialPrivateKeyEntry.Tags = []string{"tag one", "tag two"} // testing multi-word tags

	updatedPrivateKeyData, ok := updatedCredentialPrivateKeyEntry.GetCredentialPrivayeKey()
	if !ok {
		t.Fatalf("Failed to get credential access code data from entry: %v", updatedCredentialAccessCodeEntry)
	}
	updatedPrivateKeyData.Passphrase = updatedPrivateKeyData.Passphrase + "-updated"
	updatedPrivateKeyData.OverridePassword = updatedPrivateKeyData.OverridePassword + "-updated"
	updatedCredentialPrivateKeyEntry.Data = updatedPrivateKeyData

	updatedCredentialPrivateKeyEntry, err = testClient.Entries.Credential.Update(updatedCredentialPrivateKeyEntry)
	if err != nil {
		t.Fatalf("Failed to update entry: %v", err)
	}
}

func test_DeleteUserEntry(t *testing.T) {
	// Credential/AccessCode
	err := testClient.Entries.Credential.Delete(*testCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to delete AccessCode entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialAccessCodeEntry)
	if err == nil {
		t.Fatalf("AccessCode entry still exists after deletion: %s", *testCredentialAccessCodeEntryID)
	}

	// Credential/ApiKey
	err = testClient.Entries.Credential.Delete(*testCredentialApiKeyEntry)
	if err != nil {
		t.Fatalf("Failed to delete ApiKey entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialApiKeyEntry)
	if err == nil {
		t.Fatalf("ApiKey entry still exists after deletion: %s", *testCredentialApiKeyEntryID)
	}

	// Credential/AzureServicePrincipal
	err = testClient.Entries.Credential.Delete(*testCredentialAzureServicePrincipalEntry)
	if err != nil {
		t.Fatalf("Failed to delete AzureServicePrincipal entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialAzureServicePrincipalEntry)
	if err == nil {
		t.Fatalf("AzureServicePrincipal entry still exists after deletion: %s", *testCredentialAzureServicePrincipalEntryID)
	}

	// Credential/ConnectionString
	err = testClient.Entries.Credential.Delete(*testCredentialConnectionStringEntry)
	if err != nil {
		t.Fatalf("Failed to delete ConnectionString entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialConnectionStringEntry)
	if err == nil {
		t.Fatalf("ConnectionString entry still exists after deletion: %s", *testCredentialConnectionStringEntryID)
	}

	// Credential/Default
	err = testClient.Entries.Credential.Delete(*testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to delete Default entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialDefaultEntry)
	if err == nil {
		t.Fatalf("Default entry still exists after deletion: %s", *testCredentialDefaultEntryID)
	}

	// Credential/PrivateKey
	err = testClient.Entries.Credential.Delete(*testCredentialPrivateKeyEntry)
	if err != nil {
		t.Fatalf("Failed to delete PrivateKey entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialPrivateKeyEntry)
	if err == nil {
		t.Fatalf("PrivateKey entry still exists after deletion: %s", *testCredentialPrivateKeyEntryID)
	}
}
