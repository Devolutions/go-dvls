package dvls

import (
	"log"
	"os"
	"testing"
)

const testVaultId string = "e0f4f35d-8cb5-40d9-8b2b-35c96ea1c9b5"

var testClient Client

func TestMain(m *testing.M) {
	err := setupTestClient()
	if err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Client(t *testing.T) {
	t.Run("isLogged", test_isLogged)
	t.Run("GetServerInfo", test_GetServerInfo)
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

func test_GetServerInfo(t *testing.T) {
	info, err := testClient.GetServerInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("server info: %#v", info)
}

func setupTestClient() error {
	c, err := NewClient(os.Getenv("TEST_USER"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_INSTANCE"))
	if err != nil {
		return err
	}

	testClient = c

	return nil
}
