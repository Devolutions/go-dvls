package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostGetAndDetails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("includePasswords") == "false" {
			w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"h1","path":"F","type":"Host","tags":["t"],"modifiedOn":"2026-11-21T17:38:36Z","data":{"host":"10.0.0.1","username":"hu","domain":"","credentialMode":"Custom"}}`))
			return
		}
		w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"h1","path":"F","type":"Host","tags":["t"],"data":{"host":"10.0.0.1","username":"hu","password":"hp"}}`))
	})

	client := newTestClient(t, mux)

	host, err := client.Entries.Host.Get(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, "h1", host.EntryName)
	assert.Equal(t, "F", host.EntryFolderPath)
	assert.Equal(t, "10.0.0.1", host.HostDetails.Host)
	assert.Equal(t, "hu", host.HostDetails.Username)
	assert.Nil(t, host.HostDetails.Password)
	assert.Equal(t, EntryCredentialModeCustom, host.HostDetails.CredentialMode)
	require.NotNil(t, host.ModifiedDate)

	host, err = client.Entries.Host.GetHostDetails(host)
	require.NoError(t, err)
	require.NotNil(t, host.HostDetails.Password)
	assert.Equal(t, "hp", *host.HostDetails.Password)
}

func TestHostGetHostDetails_LinkedCredential(t *testing.T) {
	const credentialEntryID = "bbbbbbbb-cccc-dddd-eeee-ffffffffffff"

	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"h1","path":"","type":"Host","tags":[],"data":{"host":"10.0.0.1","username":"","password":"","credentialMode":"Linked","credentialEntryId":"` + credentialEntryID + `"}}`))
	})

	client := newTestClient(t, mux)

	host, err := client.Entries.Host.GetHostDetails(EntryHost{Id: testEntryID, VaultId: testVaultID})
	require.NoError(t, err)
	assert.Equal(t, EntryCredentialModeLinked, host.HostDetails.CredentialMode)
	assert.Equal(t, credentialEntryID, host.HostDetails.CredentialEntryId)
	require.NotNil(t, host.HostDetails.Password)
	assert.Empty(t, *host.HostDetails.Password, "the stored password of a linked entry is not the credential it uses")
}

func TestHostGet_WrongType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testCredentialJSON))
	})

	client := newTestClient(t, mux)

	_, err := client.Entries.Host.Get(testVaultID, testEntryID)
	assert.True(t, IsUnsupportedEntryType(err))
}
