package dvls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsiteGetAndDetails(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("includePasswords") == "false" {
			w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"w1","path":"","type":"WebBrowser","subType":"GoogleChrome","tags":[],"data":{"url":"https://x","username":"wu"}}`))
			return
		}
		w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"w1","path":"","type":"WebBrowser","subType":"GoogleChrome","tags":[],"data":{"url":"https://x","username":"wu","password":"wp"}}`))
	})

	client := newTestClient(t, mux)

	site, err := client.Entries.Website.Get(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, ServerConnectionSubTypeGoogleChrome, site.ConnectionSubType)
	assert.Equal(t, "https://x", site.WebsiteDetails.URL)
	assert.Equal(t, "wu", site.WebsiteDetails.Username)
	assert.Nil(t, site.WebsiteDetails.Password)

	site, err = client.Entries.Website.GetWebsiteDetails(site)
	require.NoError(t, err)
	require.NotNil(t, site.WebsiteDetails.Password)
	assert.Equal(t, "wp", *site.WebsiteDetails.Password)
}
