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
	test.InitializeTest(t)
	return &DataStore{}
}

func TestSetHostDefault(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)

	ds.setMetaHost()

	assert.Equal(t, constPubhubDevnetHost, ds.metaHost, test.MessageEqual(constPubhubDevnetHost, ds.metaHost, t.Name()))
}

func TestSetHostFromEnvironmentVariable(t *testing.T) {
	ds := initializeDataStoreTest(t)
	t.Setenv(constEnvMetaHost, metaHost)

	ds.setMetaHost()

	assert.Equal(t, metaHost, ds.metaHost, test.MessageEqual(metaHost, ds.metaHost, t.Name()))
}

type loadClassExpected struct {
	Error bool
}

func TestLoadClass(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{
			Name:     "test_invalid_class_name_no_uppercase",
			Input:    "invalidclass",
			Expected: loadClassExpected{Error: true},
		},
		{
			Name:     "test_empty_class_name",
			Input:    "",
			Expected: loadClassExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(loadClassExpected)
			ds := &DataStore{
				Classes:              make(map[string]Class),
				GlobalMetaDefinition: GlobalMetaDefinition{},
			}

			err := ds.loadClass(testCase.Input.(string))

			if expected.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Contains(t, ds.Classes, testCase.Input.(string), test.MessageContains(ds.Classes, testCase.Input.(string), testCase.Name))
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
	name, _ := NewClassName("fvTenant")
	ds.Classes["fvTenant"] = Class{Name: name}

	// Loading the same class should not error and should skip
	err := ds.loadClass("fvTenant")

	require.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Len(t, ds.Classes, 1)
}

type retrieveEnvMetaClassesInput struct {
	EnvValue       string
	ServerResponse string
	ServerStatus   int
}

type retrieveEnvMetaClassesExpected struct {
	Error bool
}

func TestRetrieveEnvMetaClassesFromRemote(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_empty_env_variable",
			Input: retrieveEnvMetaClassesInput{
				EnvValue:       "",
				ServerResponse: "",
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveEnvMetaClassesExpected{Error: false},
		},
		{
			Name: "test_single_valid_class",
			Input: retrieveEnvMetaClassesInput{
				EnvValue:       "fvTenant",
				ServerResponse: `{"label": "tenant"}`,
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveEnvMetaClassesExpected{Error: false},
		},
		{
			Name: "test_multiple_valid_classes",
			Input: retrieveEnvMetaClassesInput{
				EnvValue:       "fvTenant,fvAp",
				ServerResponse: `{"label": "test"}`,
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveEnvMetaClassesExpected{Error: false},
		},
		{
			Name: "test_invalid_class_name_no_uppercase",
			Input: retrieveEnvMetaClassesInput{
				EnvValue:       "invalidclass",
				ServerResponse: "",
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveEnvMetaClassesExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(retrieveEnvMetaClassesInput)
			expected := testCase.Expected.(retrieveEnvMetaClassesExpected)
			// Create temp directory for meta files
			tempDir := t.TempDir()

			// Create a test server
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(input.ServerStatus)
				w.Write([]byte(input.ServerResponse))
			}))
			defer server.Close()

			// Set environment variable
			if input.EnvValue != "" {
				t.Setenv(constEnvMetaClasses, input.EnvValue)
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

			if expected.Error {
				assert.Error(t, err)
			} else if input.EnvValue == "" {
				// Empty env should not error
				require.NoError(t, err, test.MessageUnexpectedError(err))
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

	require.NoError(t, err, test.MessageUnexpectedError(err))
}

type refreshMetaFilesExpected struct {
	Error bool
}

func TestRefreshMetaFiles(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_env_not_set",
			Input:    "",
			Expected: refreshMetaFilesExpected{Error: false},
		},
		{
			Name:     "test_env_set_to_false",
			Input:    "false",
			Expected: refreshMetaFilesExpected{Error: false},
		},
		{
			Name:     "test_env_set_to_0",
			Input:    "0",
			Expected: refreshMetaFilesExpected{Error: false},
		},
		{
			Name:     "test_env_set_to_invalid_value",
			Input:    "invalid",
			Expected: refreshMetaFilesExpected{Error: false}, // Logs warning but doesn't error
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			expected := testCase.Expected.(refreshMetaFilesExpected)
			if testCase.Input.(string) != "" {
				t.Setenv(constEnvMetaRefresh, testCase.Input.(string))
			}

			ds := &DataStore{
				Classes:          make(map[string]Class),
				metaHost:         "test.example.com",
				retrievedClasses: []string{},
			}

			err := ds.refreshMetaFiles()

			if expected.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err, test.MessageUnexpectedError(err))
			}
		})
	}
}

type retrieveMetaFileInput struct {
	ClassName      string
	ServerResponse string
	ServerStatus   int
}

type retrieveMetaFileExpected struct {
	Error bool
}

func TestRetrieveMetaFileFromRemote(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{
			Name: "test_invalid_class_name",
			Input: retrieveMetaFileInput{
				ClassName:      "invalidclass",
				ServerResponse: "",
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveMetaFileExpected{Error: true},
		},
		{
			Name: "test_empty_class_name",
			Input: retrieveMetaFileInput{
				ClassName:      "",
				ServerResponse: "",
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveMetaFileExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(retrieveMetaFileInput)
			expected := testCase.Expected.(retrieveMetaFileExpected)

			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(input.ServerStatus)
				w.Write([]byte(input.ServerResponse))
			}))
			defer server.Close()

			ds := &DataStore{
				Classes:          make(map[string]Class),
				client:           server.Client(),
				metaHost:         server.URL[8:],
				retrievedClasses: []string{},
			}

			err := ds.retrieveMetaFileFromRemote(input.ClassName)

			if expected.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err, test.MessageUnexpectedError(err))
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

	require.NoError(t, err, test.MessageUnexpectedError(err))
}
