//go:build integration

package dvls

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testClient  Client
	testVaultId string // Vault holding the pre-created host and website entries
)

func TestMain(m *testing.M) {
	testVaultId = os.Getenv("TEST_VAULT_ID")

	err := setupTestClient()
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_Client(t *testing.T) {
	t.Run("ReloginAfterLogout", test_ReloginAfterLogout)
}

// test_ReloginAfterLogout ends the client's session server-side, as a token expiry does, and checks that the next
// call logs in again.
func test_ReloginAfterLogout(t *testing.T) {
	if testClient.credential.useApiKey() {
		t.Skip("re-login does not apply to API key authentication")
	}

	staleToken := testClient.credential.currentToken()
	logoutUrl, err := testClient.v2Url(nil, "logout")
	require.NoError(t, err)
	_, err = testClient.rawRequestWithContext(t.Context(), logoutUrl, http.MethodPost, defaultContentType, staleToken, nil)
	require.NoError(t, err)

	_, err = testClient.GetServerInfo()
	require.NoError(t, err)
	assert.NotEqual(t, staleToken, testClient.credential.currentToken())
}

func setupTestClient() error {
	var c Client
	var err error

	if apiKey := os.Getenv("TEST_API_KEY"); apiKey != "" {
		c, err = NewClientWithApiKey(apiKey, os.Getenv("TEST_INSTANCE"))
	} else {
		c, err = NewClient(os.Getenv("TEST_USER"), os.Getenv("TEST_PASSWORD"), os.Getenv("TEST_INSTANCE"))
	}
	if err != nil {
		return err
	}

	testClient = c

	return nil
}
