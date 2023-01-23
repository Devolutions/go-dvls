package dvls

import (
	"os"
	"testing"
)

const testEntry string = "76a4fcf6-fec1-4297-bc1e-a327841055ad"

var testClient Client

func Test_NewClient(t *testing.T) {
	c, err := NewClient(os.Getenv("TEST_USER"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_INSTANCE"))
	if err != nil {
		t.Fatal(err)
	}

	testClient = c

	t.Run("isLogged", test_isLogged)
	t.Run("GetSecret", test_GetSecret)
}

func test_GetSecret(t *testing.T) {
	testSecret := DvlsSecret{
		ID:       testEntry,
		Username: "TestK8s",
		Password: "TestK8sPassword",
	}
	secret, err := testClient.GetSecret(testEntry)
	if err != nil {
		t.Fatal(err)
	}

	if secret != testSecret {
		t.Fatalf("fetched secret did not match test secret. Expected %#v, got %#v", testSecret, secret)
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

	validToken := testClient.credential.token
	testClient.credential.token = "placeholder"
	islogged, err = testClient.isLogged()
	if err != nil {
		t.Fatal(err)
	}
	if islogged {
		t.Fatalf("expected token to be invalid but isLogged returned %t", islogged)
	}
	testClient.credential.token = validToken
}
