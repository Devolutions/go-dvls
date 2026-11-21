//go:build integration

package dvls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EntryCertificate(t *testing.T) {
	vault := createTestVault(t, "certificates")

	entry := EntryCertificate{
		VaultId:     vault.Id,
		Name:        "TestK8sCertificate",
		Description: "Test certificate",
		Password:    "TestK8sCertificatePassword",
		Tags:        []string{"test", "k8s"},
	}

	t.Run("URL", func(t *testing.T) {
		urlEntry := entry
		urlEntry.CertificateIdentifier = "https://devolutions.net/"

		created, err := testClient.Entries.Certificate.NewURL(urlEntry)
		require.NoError(t, err)
		require.NotEmpty(t, created.Id)
		t.Cleanup(func() { _ = testClient.Entries.Certificate.Delete(vault.Id, created.Id) })

		fetched, err := testClient.Entries.Certificate.Get(vault.Id, created.Id)
		require.NoError(t, err)
		assert.Equal(t, EntryCertificateDataModeURL, fetched.GetDataMode())
		assert.Equal(t, urlEntry.CertificateIdentifier, fetched.CertificateIdentifier)
		assert.Equal(t, urlEntry.Password, fetched.Password)
		assert.ElementsMatch(t, urlEntry.Tags, fetched.Tags)

		withoutPassword := fetched
		withoutPassword.Password = ""
		withPassword, err := testClient.Entries.Certificate.GetPassword(withoutPassword)
		require.NoError(t, err)
		assert.Equal(t, urlEntry.Password, withPassword.Password)

		fetched.Name = "TestK8sUpdatedEntry"
		fetched.CertificateIdentifier = "https://devolutions.net/updated"
		_, err = testClient.Entries.Certificate.Update(fetched)
		require.NoError(t, err)

		updated, err := testClient.Entries.Certificate.Get(vault.Id, created.Id)
		require.NoError(t, err)
		assert.Equal(t, "TestK8sUpdatedEntry", updated.Name)
		assert.Equal(t, "https://devolutions.net/updated", updated.CertificateIdentifier)
		assert.Equal(t, EntryCertificateDataModeURL, updated.GetDataMode())
	})

	t.Run("File", func(t *testing.T) {
		fileBytes := []byte("go-dvls test certificate content")
		fileName := "test.p12"
		if path := os.Getenv("TEST_CERTIFICATE_FILE_PATH"); path != "" {
			var err error
			fileBytes, err = os.ReadFile(path)
			require.NoError(t, err)
			fileName = filepath.Base(path)
		}

		fileEntry := entry
		fileEntry.CertificateIdentifier = fileName

		created, err := testClient.Entries.Certificate.NewFile(fileEntry, fileBytes)
		require.NoError(t, err)
		t.Cleanup(func() { _ = testClient.Entries.Certificate.Delete(vault.Id, created.Id) })

		content, err := testClient.Entries.Certificate.GetFileContent(vault.Id, created.Id)
		require.NoError(t, err)
		assert.Equal(t, fileBytes, content)

		fetched, err := testClient.Entries.Certificate.Get(vault.Id, created.Id)
		require.NoError(t, err)
		assert.Equal(t, EntryCertificateDataModeFile, fetched.GetDataMode())
		assert.Equal(t, fileName, fetched.CertificateIdentifier)

		fetched.Description = "updated"
		_, err = testClient.Entries.Certificate.Update(fetched)
		require.NoError(t, err)

		content, err = testClient.Entries.Certificate.GetFileContent(vault.Id, created.Id)
		require.NoError(t, err)
		assert.Equal(t, fileBytes, content, "an update without content keeps the stored file")
	})
}
