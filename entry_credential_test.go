package dvls

import (
	"testing"
)

var (
	testCredentialAccessCodeEntryId *string
	testCredentialAccessCodeEntry   *Entry

	testCredentialApiKeyEntryId *string
	testCredentialApiKeyEntry   *Entry

	testCredentialAzureServicePrincipalEntryId *string
	testCredentialAzureServicePrincipalEntry   *Entry

	testCredentialConnectionStringEntryId *string
	testCredentialConnectionStringEntry   *Entry

	testCredentialDefaultEntryId *string
	testCredentialDefaultEntry   *Entry

	testCredentialPrivateKeyEntryId *string
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
		Id:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsAccessCode",
		Path:        "go-dvls\\accesscode",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeAccessCode,
		Description: "Test AccessCode entry",
		Tags:        []string{"accesscode"},

		Data: EntryCredentialAccessCodeData{
			Password: "abc-123",
		},
	}

	newCredentialAccessCodeEntryId, err := testClient.Entries.Credential.New(testCredentialAccessCodeEntry)
	if err != nil {
		t.Fatalf("Failed to create new AccessCode entry: %v", err)
	}

	if newCredentialAccessCodeEntryId == "" {
		t.Fatal("New AccessCode entry Id is empty after creation.")
	}

	testCredentialAccessCodeEntryId = &newCredentialAccessCodeEntryId

	// Credential/ApiKey
	testCredentialApiKeyEntry := Entry{
		Id:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsApiKey",
		Path:        "go-dvls\\apikey",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeApiKey,
		Description: "Test ApiKey entry",
		Tags:        []string{"apikey"},

		Data: EntryCredentialApiKeyData{
			ApiId:    "abcd1234-abcd-1234-abcd-1234abcd1234",
			ApiKey:   "123-abc",
			TenantId: "00000000-aaaa-bbbb-cccc-000000000000",
		},
	}

	newCredentialApiKeyEntryId, err := testClient.Entries.Credential.New(testCredentialApiKeyEntry)
	if err != nil {
		t.Fatalf("Failed to create new ApiKey entry: %v", err)
	}

	if newCredentialApiKeyEntryId == "" {
		t.Fatal("New ApiKey entry Id is empty after creation.")
	}

	testCredentialApiKeyEntryId = &newCredentialApiKeyEntryId

	// Credential/AzureServicePrincipal
	testCredentialAzureServicePrincipalEntry := Entry{
		Id:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsAzureServicePrincipal",
		Path:        "go-dvls\\azureserviceprincipal",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeAzureServicePrincipal,
		Description: "Test AzureServicePrincipal entry",
		Tags:        []string{"azureserviceprincipal"},

		Data: EntryCredentialAzureServicePrincipalData{
			ClientId:     "abcd1234-abcd-1234-abcd-1234abcd1234",
			ClientSecret: "123-abc",
			TenantId:     "00000000-aaaa-bbbb-cccc-000000000000",
		},
	}

	newCredentialAzureServicePrincipalEntryId, err := testClient.Entries.Credential.New(testCredentialAzureServicePrincipalEntry)
	if err != nil {
		t.Fatalf("Failed to create new AzureServicePrincipal entry: %v", err)
	}

	if newCredentialAzureServicePrincipalEntryId == "" {
		t.Fatal("New AzureServicePrincipal entry Id is empty after creation.")
	}

	testCredentialAzureServicePrincipalEntryId = &newCredentialAzureServicePrincipalEntryId

	// Credential/ConnectionString
	testCredentialConnectionStringEntry := Entry{
		Id:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsConnectionString",
		Path:        "go-dvls\\connectionstring",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeConnectionString,
		Description: "Test ConnectionString entry",
		Tags:        []string{"connectionstring"},

		Data: EntryCredentialConnectionStringData{
			ConnectionString: "Server=tcp:example.database.windows.net,1433;Initial Catalog=exampledb;Persist Security Info=False;User ID=exampleuser;Password=examplepassword;",
		},
	}

	newCredentialConnectionStringEntryId, err := testClient.Entries.Credential.New(testCredentialConnectionStringEntry)
	if err != nil {
		t.Fatalf("Failed to create new ConnectionString entry: %v", err)
	}

	if newCredentialConnectionStringEntryId == "" {
		t.Fatal("New ConnectionString entry Id is empty after creation.")
	}

	testCredentialConnectionStringEntryId = &newCredentialConnectionStringEntryId

	// Credential/Default
	testCredentialDefaultEntry := Entry{
		VaultId:     testVaultId,
		Name:        "TestGoDvlsUsernamePassword",
		Path:        "go-dvls\\usernamepassword",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypeDefault,
		Description: "Test Username/Password entry",
		Tags:        []string{"usernamepassword"},

		Data: EntryCredentialDefaultData{
			Domain:   "www.example.com",
			Password: "abc-123",
			Username: "john.doe",
		},
	}

	newCredentialDefaultEntryId, err := testClient.Entries.Credential.New(testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to create new Default entry: %v", err)
	}

	if newCredentialDefaultEntryId == "" {
		t.Fatal("New Default entry Id is empty after creation.")
	}

	testCredentialDefaultEntryId = &newCredentialDefaultEntryId

	// Credential/PrivateKey
	testCredentialPrivateKeyEntry := Entry{
		Id:          "",
		VaultId:     testVaultId,
		Name:        "TestGoDvlsPrivateKey",
		Path:        "go-dvls\\privatekey",
		Type:        EntryCredentialType,
		SubType:     EntryCredentialSubTypePrivateKey,
		Description: "Test Secret entry",
		Tags:        []string{"testtag"},

		Data: EntryCredentialPrivateKeyData{
			PrivateKey:       "-----BEGIN PRIVATE KEY-----\abcdefghijklmnopqrstuvwxyz1234567890...\n-----END PRIVATE",
			PublicKey:        "-----BEGIN PUBLIC KEY-----\abcdefghijklmnopqrstuvwxyz...\n-----END PUBLIC KEY-----",
			OverridePassword: "override-password",
			Passphrase:       "passphrase",
		},
	}

	newCredentialPrivateKeyEntryId, err := testClient.Entries.Credential.New(testCredentialPrivateKeyEntry)
	if err != nil {
		t.Fatalf("Failed to create new PrivateKey entry: %v", err)
	}

	if newCredentialPrivateKeyEntryId == "" {
		t.Fatal("New PrivateKey entry Id is empty after creation.")
	}

	testCredentialPrivateKeyEntryId = &newCredentialPrivateKeyEntryId
}

func test_GetUserEntry(t *testing.T) {
	// Credential/AccessCode
	credentialAccessCodeEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialAccessCodeEntryId)
	if err != nil {
		t.Fatalf("Failed to get AccessCode entry: %v", err)
	}

	if credentialAccessCodeEntry.Id == "" {
		t.Fatalf("AccessCode entry Id is empty after GET: %v", credentialAccessCodeEntry)
	}

	testCredentialAccessCodeEntry = &credentialAccessCodeEntry

	// Credential/ApiKey
	credentialApiKeyEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialApiKeyEntryId)
	if err != nil {
		t.Fatalf("Failed to get ApiKey entry: %v", err)
	}

	if credentialApiKeyEntry.Id == "" {
		t.Fatalf("ApiKey entry Id is empty after GET: %v", credentialApiKeyEntry)
	}

	testCredentialApiKeyEntry = &credentialApiKeyEntry

	// Credential/AzureServicePrincipal
	credentialAzureServicePrincipalEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialAzureServicePrincipalEntryId)
	if err != nil {
		t.Fatalf("Failed to get AzureServicePrincipal entry: %v", err)
	}

	if credentialAzureServicePrincipalEntry.Id == "" {
		t.Fatalf("AzureServicePrincipal entry Id is empty after GET: %v", credentialAzureServicePrincipalEntry)
	}

	testCredentialAzureServicePrincipalEntry = &credentialAzureServicePrincipalEntry

	// Credential/ConnectionString
	credentialConnectionStringEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialConnectionStringEntryId)
	if err != nil {
		t.Fatalf("Failed to get ConnectionString entry: %v", err)
	}

	if credentialConnectionStringEntry.Id == "" {
		t.Fatalf("ConnectionString entry Id is empty after GET: %v", credentialConnectionStringEntry)
	}

	testCredentialConnectionStringEntry = &credentialConnectionStringEntry

	// Credential/Default
	credentialDefaultEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialDefaultEntryId)
	if err != nil {
		t.Fatalf("Failed to get Default entry: %v", err)
	}

	if credentialDefaultEntry.Id == "" {
		t.Fatalf("Default entry Id is empty after GET: %v", credentialDefaultEntry)
	}

	testCredentialDefaultEntry = &credentialDefaultEntry

	// Credential/PrivateKey
	credentialPrivateKeyEntry, err := testClient.Entries.Credential.GetById(testVaultId, *testCredentialPrivateKeyEntryId)
	if err != nil {
		t.Fatalf("Failed to get PrivateKey entry: %v", err)
	}

	if credentialPrivateKeyEntry.Id == "" {
		t.Fatalf("PrivateKey entry Id is empty after GET: %v", credentialPrivateKeyEntry)
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
		t.Fatalf("AccessCode entry still exists after deletion: %s", *testCredentialAccessCodeEntryId)
	}

	// Credential/ApiKey
	err = testClient.Entries.Credential.Delete(*testCredentialApiKeyEntry)
	if err != nil {
		t.Fatalf("Failed to delete ApiKey entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialApiKeyEntry)
	if err == nil {
		t.Fatalf("ApiKey entry still exists after deletion: %s", *testCredentialApiKeyEntryId)
	}

	// Credential/AzureServicePrincipal
	err = testClient.Entries.Credential.Delete(*testCredentialAzureServicePrincipalEntry)
	if err != nil {
		t.Fatalf("Failed to delete AzureServicePrincipal entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialAzureServicePrincipalEntry)
	if err == nil {
		t.Fatalf("AzureServicePrincipal entry still exists after deletion: %s", *testCredentialAzureServicePrincipalEntryId)
	}

	// Credential/ConnectionString
	err = testClient.Entries.Credential.Delete(*testCredentialConnectionStringEntry)
	if err != nil {
		t.Fatalf("Failed to delete ConnectionString entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialConnectionStringEntry)
	if err == nil {
		t.Fatalf("ConnectionString entry still exists after deletion: %s", *testCredentialConnectionStringEntryId)
	}

	// Credential/Default
	err = testClient.Entries.Credential.Delete(*testCredentialDefaultEntry)
	if err != nil {
		t.Fatalf("Failed to delete Default entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialDefaultEntry)
	if err == nil {
		t.Fatalf("Default entry still exists after deletion: %s", *testCredentialDefaultEntryId)
	}

	// Credential/PrivateKey
	err = testClient.Entries.Credential.Delete(*testCredentialPrivateKeyEntry)
	if err != nil {
		t.Fatalf("Failed to delete PrivateKey entry: %v", err)
	}

	_, err = testClient.Entries.Credential.Get(*testCredentialPrivateKeyEntry)
	if err == nil {
		t.Fatalf("PrivateKey entry still exists after deletion: %s", *testCredentialPrivateKeyEntryId)
	}
}
