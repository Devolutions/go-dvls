package dvls

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var (
	testSSHEntryId  string
	testSSHPassword          = "testpass123"
	testSSHEntry    EntrySSH = EntrySSH{
		Description:    "Test SSH description",
		EntryName:      "TestSSH",
		ConnectionType: ServerConnectionSSHShell,
		Tags:           []string{"Test tag 1", "Test tag 2", "ssh"},
	}
	testSSHKeyEntryId string
	// the following is a dummy private key for testing purposes
	testSSHPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jYmMAAAAGYmNyeXB0AAAAGAAAABA9Z33hXL
3PvlqyP2/m1oVeAAAAQAAAAAEAAAAAAAADwIPtw5xp00RQniA43BoCxMakzlaWGBXu4Baq
s2bFgslANxdX4g8GQioUvQiIW5XVLHAaUbyJRp3hugxfiOO+xRHcBIGE14CbTCPSdiq8Ng
VzoUO7zyhDk7wksQ9nIKi6tosU+RtSNHWUgYt3KWDB87Jb6LIvEVXML7tl8RvDZYu/Skec
g+muWGMxBltWt5g62HTCUcO6QB67Ipz2qoUp/cRkG790Q28vNn+n12nnDIHBH4CqeGU1nO
5pjgIiMY07JbVmPlG9v7EwL2XHykx/kc3x2cFK9VIRSxGHHdu2JPz0CVK5Cs2wZSLttd+8
UA4DwhqXfIpdjl6QyDlQge1njSvg3Q1P1sfURON/0FfraKkWGYB69JxhBNcyoETftwj0Z7
3JjCvBkM+T/9gHwJnbvt29x1v37wsvZiWAajzN+61Wpo9+IeBX6d9Jf2WSTwb247645mEA
fofyQE2OudsnWvFYVnWD8JfN6ILI5q/EMeYjjaCbdWZItF0ZDfar0Q2o9LO6DpI519DRZG
z4bfx6f6QiiA2jPfOXB0VhZD0GgkZNbFOJxh7g4CpvDIhOPXN/Jd3vMAyFa8lpEPe13Z+0
JYdfDjY/sx9VRpSMeIdTQLycNESzYKtQmEO/50kTlAnReC62HjLWmVJLtCpFw9pmuCPsdn
1hVv8W3TTO3mZD7ouGQfa+B3NmKZT+pIA0A68EJ2DZZwW3sIpal52gGnw0TU7jpAQjirzJ
zQIIGRUdW8l7ZqMEjvO9d23EFHLaZ1PRoePyijoh6TD/6Uh8rWOtOa0WPBRWeKnIntYQM/
fXZTVAvOqVnjlwakOfOzG0pAXmMxlXO5iS9sNg6o/dxBdedJELzyU0DOAdM6L5PRjkySvR
8Wo7qlNhl5xM+FyH047ox19+w0moRe3/Wz6E/U5Eqo1u73igkb1XZ0ENL5STT2uXopQXb7
wtEVyv/IzGlw09e5rmyrEoFd03s7KCLuqb71xwOD4fvKXlwZpqX8dcpQKDTWbTkyqSBTFH
QSdAfbvfzcaUBVLW+b8TLXUBOMUvyiK91mJz1GFIcoh9L8fyWqq1gSKNhm/FH7eIahU5Tv
EahidkUcp1pPdqrtHk3D75Naial6v2g6R9+hD299akroK2wwoiNkARjPyaYGQ1Ck6vjDtw
qO6xETt9SAkNgkpaPV9hO1ldnGtdncO26LyInjnaKRT5ZugeGbWr3Ef3rxmeimi6egpBv9
tl7MnBKy3YywVoKhcrkosPVz6D2eB2E5ti0AfFCajg3qEAEUL8eji1Cg==
-----END OPENSSH PRIVATE KEY-----`
	testSSHKeyPassphrase          = "testpassphrase"
	testSSHKeyEntry      EntrySSH = EntrySSH{
		Description:    "Test SSH Key description",
		EntryName:      "TestSSHKey",
		ConnectionType: ServerConnectionSSHShell,
		Tags:           []string{"Test tag 1", "Test tag 2", "ssh"},
	}
)

const (
	testSSHUsername string = "testuser"
	testSSHHost     string = "myhost"
	testSSHPort     int    = 22
	testSSHKeyHost  string = "myotherhost"
	testSSHKeyPort  int    = 23
)

func Test_EntrySSH(t *testing.T) {
	testSSHEntryId = os.Getenv("TEST_SSH_ENTRY_ID")
	testSSHEntry.ID = testSSHEntryId
	testSSHEntry.VaultId = testVaultId
	testSSHEntry.SSHDetails = EntrySSHAuthDetails{
		Username: testSSHUsername,
		Host:     testSSHHost,
		Port:     testSSHPort,
	}

	testSSHKeyEntryId = os.Getenv("TEST_SSH_KEY_ENTRY_ID")
	testSSHKeyEntry.ID = testSSHKeyEntryId
	testSSHKeyEntry.VaultId = testVaultId
	testSSHKeyEntry.SSHDetails = EntrySSHAuthDetails{
		Host:       testSSHKeyHost,
		Port:       testSSHKeyPort,
		Passphrase: &testSSHKeyPassphrase,
		PrivateKey: &testSSHPrivateKey,
	}

	t.Run("GetEntry", test_GetSSHEntry)
	t.Run("GetEntrySSH", test_GetSSHDetails)
	t.Run("GetKeyEntry", test_GetSSHKeyEntry)
	t.Run("GetKeyDetails", test_GetSSHKeyDetails)
}

func test_GetSSHEntry(t *testing.T) {
	entry, err := testClient.Entries.SSH.Get(testSSHEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	testSSHEntry.ModifiedDate = entry.ModifiedDate
	if !reflect.DeepEqual(entry, testSSHEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testSSHEntry, entry)
	}
}

func test_GetSSHDetails(t *testing.T) {
	entry, err := testClient.Entries.SSH.Get(testSSHEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	entryWithSensitiveData, err := testClient.Entries.SSH.GetSSHDetails(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.SSHDetails.Password = entryWithSensitiveData.SSHDetails.Password

	expectedDetails := testSSHEntry.SSHDetails

	expectedDetails.Password = &testSSHPassword

	if !reflect.DeepEqual(expectedDetails, entry.SSHDetails) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", expectedDetails, entry.SSHDetails)
	}
}

func test_GetSSHKeyEntry(t *testing.T) {
	entry, err := testClient.Entries.SSH.Get(testSSHKeyEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch sensitive data
	entryWithSensitiveData, err := testClient.Entries.SSH.GetSSHDetails(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.SSHDetails = entryWithSensitiveData.SSHDetails
	testSSHKeyEntry.ModifiedDate = entry.ModifiedDate

	fmt.Printf("Expected Entry: %+v\n", testSSHKeyEntry)
	fmt.Printf("Actual Entry: %+v\n", entry)

	if !reflect.DeepEqual(entry, testSSHKeyEntry) {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testSSHKeyEntry, entry)
	}
}

func test_GetSSHKeyDetails(t *testing.T) {
	entry, err := testClient.Entries.SSH.Get(testSSHKeyEntry.ID)
	if err != nil {
		t.Fatal(err)
	}

	entryWithSensitiveData, err := testClient.Entries.SSH.GetSSHDetails(entry)
	if err != nil {
		t.Fatal(err)
	}

	entry.SSHDetails.Password = entryWithSensitiveData.SSHDetails.Password
	entry.SSHDetails.PrivateKey = entryWithSensitiveData.SSHDetails.PrivateKey
	entry.SSHDetails.Passphrase = entryWithSensitiveData.SSHDetails.Passphrase

	expectedDetails := testSSHKeyEntry.SSHDetails

	expectedDetails.Password = &testSSHPrivateKey

	if !reflect.DeepEqual(expectedDetails, entry.SSHDetails) {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", expectedDetails, entry.SSHDetails)
	}
}
