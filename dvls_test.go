package dvls

import (
	"log"
	"os"
	"testing"
)

var (
	testClient  Client
	testVaultId string // Used by legacy tests (certificate, host, website)
)

func TestMain(m *testing.M) {
	testVaultId = os.Getenv("TEST_VAULT_ID") // Optional, only for legacy tests

	err := setupTestClient()
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Client(t *testing.T) {
	t.Run("isLogged", test_isLogged)
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

func setupTestClient() error {
	c, err := NewClient(os.Getenv("TEST_USER"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_INSTANCE"))
	if err != nil {
		return err
	}

	testClient = c

	return nil
}
