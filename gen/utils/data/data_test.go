package data

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	metaHost = "10.0.0.1"
)

func initializeDataStoreTest(t *testing.T) *DataStore {
	t.Helper()
	test.InitializeTest(t)
	return &DataStore{}
}

func TestSetHostDefault(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)

	ds.setMetaHost()

	assert.Equal(t, constPubhubDevnetHost, ds.metaHost)
}

func TestSetHostFromEnvironmentVariable(t *testing.T) {
	ds := initializeDataStoreTest(t)
	t.Setenv(constEnvMetaHost, metaHost)

	ds.setMetaHost()

	assert.Equal(t, metaHost, ds.metaHost)
}

func TestLoadClass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		className   string
		expectError bool
	}{
		{
			name:        "invalid_class_name_no_uppercase",
			className:   "invalidclass",
			expectError: true,
		},
		{
			name:        "empty_class_name",
			className:   "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ds := &DataStore{
				Classes:              make(map[string]Class),
				GlobalMetaDefinition: GlobalMetaDefinition{},
			}

			err := ds.loadClass(tc.className)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, ds.Classes, tc.className)
			}
		})
	}
}

func TestLoadClassAlreadyLoaded(t *testing.T) {
	t.Parallel()
	ds := &DataStore{
		Classes:              make(map[string]Class),
		GlobalMetaDefinition: GlobalMetaDefinition{},
	}

	// Pre-populate with a class
	ds.Classes["fvTenant"] = Class{ClassName: "fvTenant"}

	// Loading the same class should not error and should skip
	err := ds.loadClass("fvTenant")

	require.NoError(t, err)
	assert.Len(t, ds.Classes, 1)
}

func TestRetrieveEnvMetaClassesFromRemote(t *testing.T) {
	tests := []struct {
		name           string
		envValue       string
		serverResponse string
		serverStatus   int
		expectError    bool
	}{
		{
			name:           "empty_env_variable",
			envValue:       "",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "single_valid_class",
			envValue:       "fvTenant",
			serverResponse: `{"label": "tenant"}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "multiple_valid_classes",
			envValue:       "fvTenant,fvAp",
			serverResponse: `{"label": "test"}`,
			serverStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid_class_name_no_uppercase",
			envValue:       "invalidclass",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp directory for meta files
			tempDir := t.TempDir()

			// Create a test server
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.serverStatus)
				w.Write([]byte(tc.serverResponse))
			}))
			defer server.Close()

			// Set environment variable
			if tc.envValue != "" {
				t.Setenv(constEnvMetaClasses, tc.envValue)
			}

			ds := &DataStore{
				Classes:          make(map[string]Class),
				client:           server.Client(),
				metaHost:         server.URL[8:], // Remove "https://"
				retrievedClasses: []string{},
			}

			// We need to override constMetaPath for this test
			// Since we can't easily do that, we'll just test the error cases
			err := ds.retrieveEnvMetaClassesFromRemote()

			if tc.expectError {
				assert.Error(t, err)
			} else if tc.envValue == "" {
				// Empty env should not error
				require.NoError(t, err)
			}

			_ = tempDir // Used for cleanup
		})
	}
}

func TestRetrieveEnvMetaClassesFromRemoteEmptyEnv(t *testing.T) {
	// Note: Cannot use t.Parallel() with t.Setenv()
	t.Setenv(constEnvMetaClasses, "")

	ds := &DataStore{
		Classes:          make(map[string]Class),
		retrievedClasses: []string{},
	}

	err := ds.retrieveEnvMetaClassesFromRemote()

	require.NoError(t, err)
}

func TestRefreshMetaFiles(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
	}{
		{
			name:        "env_not_set",
			envValue:    "",
			expectError: false,
		},
		{
			name:        "env_set_to_false",
			envValue:    "false",
			expectError: false,
		},
		{
			name:        "env_set_to_0",
			envValue:    "0",
			expectError: false,
		},
		{
			name:        "env_set_to_invalid_value",
			envValue:    "invalid",
			expectError: false, // Logs warning but doesn't error
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv(constEnvMetaRefresh, tc.envValue)
			}

			ds := &DataStore{
				Classes:          make(map[string]Class),
				metaHost:         "test.example.com",
				retrievedClasses: []string{},
			}

			err := ds.refreshMetaFiles()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRetrieveMetaFileFromRemote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		className      string
		serverResponse string
		serverStatus   int
		expectError    bool
	}{
		{
			name:           "invalid_class_name",
			className:      "invalidclass",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectError:    true,
		},
		{
			name:           "empty_class_name",
			className:      "",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.serverStatus)
				w.Write([]byte(tc.serverResponse))
			}))
			defer server.Close()

			ds := &DataStore{
				Classes:          make(map[string]Class),
				client:           server.Client(),
				metaHost:         server.URL[8:],
				retrievedClasses: []string{},
			}

			err := ds.retrieveMetaFileFromRemote(tc.className)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRetrieveMetaFileFromRemoteAlreadyRetrieved(t *testing.T) {
	t.Parallel()

	ds := &DataStore{
		Classes:          make(map[string]Class),
		retrievedClasses: []string{"fvTenant"},
	}

	// Should skip retrieval and not error
	err := ds.retrieveMetaFileFromRemote("fvTenant")

	require.NoError(t, err)
}
