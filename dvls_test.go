package dvls

import (
	"os"
	"testing"
)

const testEntryId string = "76a4fcf6-fec1-4297-bc1e-a327841055ad"

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
	t.Run("GetSecret", test_GetSecret)
	t.Run("GetEntry", test_GetEntry)
}

func test_GetSecret(t *testing.T) {
	testSecret := DvlsSecret{
		ID:       testEntryId,
		Username: "TestK8s",
		Password: "TestK8sPassword",
	}
	secret, err := testClient.GetSecret(testEntryId)
	if err != nil {
		t.Fatal(err)
	}

	if secret != testSecret {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", testSecret, secret)
	}
}

func test_GetEntry(t *testing.T) {
	testEntry := DvlsEntry{
		ID:                testEntryId,
		Name:              "TestK8sSecret",
		ConnectionType:    ServerConnectionCredential,
		ConnectionSubType: ServerConnectionSubTypeDefault,
	}
	entry, err := testClient.GetEntry(testEntryId)
	if err != nil {
		t.Fatal(err)
	}

	if entry != testEntry {
		t.Fatalf("fetched entry did not match test entry. Expected %#v, got %#v", testEntry, entry)
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
