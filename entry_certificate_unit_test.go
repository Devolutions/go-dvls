package dvls

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCertificateUrlJSON  = `{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"cert","path":"F","type":"Document","subType":"Certificate","description":"desc","tags":["t"],"data":{"dataMode":"Url","url":"https://example.com/c.pfx","password":"pw","expiresOn":"2099-01-01T05:00:00Z","useDefaultCredentials":true}}`
	testCertificateFileJSON = `{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"cert","path":"","type":"Document","subType":"Certificate","description":"","tags":[],"data":{"dataMode":"File","fileName":"c.p12","password":"pw","expiresOn":null,"useDefaultCredentials":false}}`
)

func TestCertificateGet(t *testing.T) {
	tests := []struct {
		name               string
		response           string
		wantMode           EntryCertificateDataMode
		wantIdentifier     string
		wantExpiration     time.Time
		wantDefaultCredUse bool
	}{
		{name: "Url", response: testCertificateUrlJSON, wantMode: EntryCertificateDataModeURL, wantIdentifier: "https://example.com/c.pfx", wantExpiration: time.Date(2099, 1, 1, 5, 0, 0, 0, time.UTC), wantDefaultCredUse: true},
		{name: "File", response: testCertificateFileJSON, wantMode: EntryCertificateDataModeFile, wantIdentifier: "c.p12"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.response))
			})

			client := newTestClient(t, mux)

			cert, err := client.Entries.Certificate.Get(testVaultID, testEntryID)
			require.NoError(t, err)
			assert.Equal(t, tt.wantMode, cert.GetDataMode())
			assert.Equal(t, tt.wantIdentifier, cert.CertificateIdentifier)
			require.NotNil(t, cert.Password)
			assert.Equal(t, "pw", *cert.Password)
			assert.Equal(t, testVaultID, cert.VaultId)
			assert.True(t, tt.wantExpiration.Equal(cert.Expiration), "expiration %v", cert.Expiration)
			assert.Equal(t, tt.wantDefaultCredUse, cert.UseDefaultCredentials)
		})
	}
}

func TestCertificateGet_WrongType(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testCredentialJSON))
	})

	client := newTestClient(t, mux)

	_, err := client.Entries.Certificate.Get(testVaultID, testEntryID)
	assert.True(t, IsUnsupportedEntryType(err))
}

func TestCertificateGetFileContent(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath+"/content", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte{0x30, 0x82, 0x00})
	})

	client := newTestClient(t, mux)

	content, err := client.Entries.Certificate.GetFileContent(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x30, 0x82, 0x00}, content)
}

func TestCertificateGetPassword(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testCertificateUrlJSON))
	})

	client := newTestClient(t, mux)

	cert, err := client.Entries.Certificate.GetPassword(EntryCertificate{Id: testEntryID, VaultId: testVaultID, Name: "kept"})
	require.NoError(t, err)
	require.NotNil(t, cert.Password)
	assert.Equal(t, "pw", *cert.Password)
	assert.Equal(t, "kept", cert.Name)
}

func TestCertificateNew(t *testing.T) {
	tests := []struct {
		name     string
		create   func(*Client, EntryCertificate) (EntryCertificate, error)
		wantData map[string]any
		response string
	}{
		{
			name:     "Url",
			create:   func(c *Client, e EntryCertificate) (EntryCertificate, error) { return c.Entries.Certificate.NewURL(e) },
			wantData: map[string]any{"dataMode": "Url", "url": "https://example.com/c.pfx", "password": "pw", "expiresOn": "2099-01-01T05:00:00Z", "useDefaultCredentials": true},
			response: testCertificateUrlJSON,
		},
		{
			name: "File",
			create: func(c *Client, e EntryCertificate) (EntryCertificate, error) {
				e.CertificateIdentifier = "c.p12"
				return c.Entries.Certificate.NewFile(e, []byte("certificate-bytes"))
			},
			wantData: map[string]any{"dataMode": "File", "fileName": "c.p12", "fileContent": base64.StdEncoding.EncodeToString([]byte("certificate-bytes")), "password": "pw", "expiresOn": "2099-01-01T05:00:00Z", "useDefaultCredentials": true},
			response: testCertificateFileJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testVaultEntriesPath, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				var body map[string]any
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				assert.Equal(t, "Document", body["type"])
				assert.Equal(t, "Certificate", body["subType"])
				assert.Equal(t, "F", body["path"])
				assert.Equal(t, tt.wantData, body["data"])

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(tt.response))
			})

			client := newTestClient(t, mux)

			expiration := time.Date(2099, 1, 1, 0, 0, 0, 0, time.FixedZone("EST", -5*3600))
			cert, err := tt.create(client, EntryCertificate{VaultId: testVaultID, Name: "cert", EntryFolderPath: "F", Password: new("pw"), Expiration: expiration, UseDefaultCredentials: true, CertificateIdentifier: "https://example.com/c.pfx"})
			require.NoError(t, err)
			assert.Equal(t, testEntryID, cert.Id)
			assert.Equal(t, new("pw"), cert.Password)
		})
	}
}

func TestCertificateUpdate_KeepsDataMode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			assert.Equal(t, "false", r.URL.Query().Get("includePasswords"), "a metadata update must not read the password")
			w.Write([]byte(testCertificateFileJSON))
		case http.MethodPut:
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, map[string]any{"fileName": "renamed.p12", "expiresOn": nil, "useDefaultCredentials": false}, body["data"], "no data mode, file content or password keeps the stored ones; a zero expiration clears it")
			assert.Equal(t, []any{}, body["tags"])
			w.Write([]byte(testCertificateFileJSON))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, mux)

	cert, err := client.Entries.Certificate.Update(EntryCertificate{Id: testEntryID, VaultId: testVaultID, Name: "cert", CertificateIdentifier: "renamed.p12"})
	require.NoError(t, err)
	assert.Equal(t, EntryCertificateDataModeFile, cert.GetDataMode())
	assert.Nil(t, cert.Password, "a kept password is unknown, not empty")
}

func TestCertificateUpdate_Password(t *testing.T) {
	tests := []struct {
		name         string
		password     *string
		wantPassword any
	}{
		{name: "Clear", password: new(""), wantPassword: ""},
		{name: "Set", password: new("new-pw"), wantPassword: "new-pw"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					w.Write([]byte(testCertificateUrlJSON))
				case http.MethodPut:
					var body map[string]any
					require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
					data, ok := body["data"].(map[string]any)
					require.True(t, ok)
					assert.Contains(t, data, "password")
					assert.Equal(t, tt.wantPassword, data["password"])
					w.Write([]byte(testCertificateUrlJSON))
				default:
					t.Errorf("unexpected method %s", r.Method)
				}
			})

			client := newTestClient(t, mux)

			cert, err := client.Entries.Certificate.Update(EntryCertificate{Id: testEntryID, VaultId: testVaultID, Name: "cert", CertificateIdentifier: "https://example.com/c.pfx", Password: tt.password})
			require.NoError(t, err)
			assert.Equal(t, tt.password, cert.Password)
		})
	}
}

func TestCertificateUpdate_FromGetSkipsRead(t *testing.T) {
	var gets int

	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			gets++
			w.Write([]byte(`{"id":"` + testEntryID + `","vaultId":"` + testVaultID + `","name":"cert","path":"","type":"Document","subType":"Certificate","description":"","tags":[],"data":{"dataMode":"Url","url":"https://example.com/c.pfx","expiresOn":"2099-01-01T05:00:00.1234567Z","useDefaultCredentials":false}}`))
		case http.MethodPut:
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, map[string]any{"url": "https://example.com/c.pfx", "expiresOn": "2099-01-01T05:00:00.1234567Z", "useDefaultCredentials": false}, body["data"], "an unchanged expiration keeps its fractional seconds")
			w.Write([]byte(testCertificateUrlJSON))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, mux)

	cert, err := client.Entries.Certificate.Get(testVaultID, testEntryID)
	require.NoError(t, err)
	assert.Equal(t, time.Date(2099, 1, 1, 5, 0, 0, 123456700, time.UTC), cert.Expiration)

	_, err = client.Entries.Certificate.Update(cert)
	require.NoError(t, err)
	assert.Equal(t, 1, gets, "an entry read with Get already has its data mode")
}

func TestCertificateDelete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(testEntryPath, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, mux)

	require.NoError(t, client.Entries.Certificate.Delete(testVaultID, testEntryID))
}
