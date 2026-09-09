package data

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

const (
	metaHost = "10.0.0.1"
)

func initializeDataStoreTest(t *testing.T) *DataStore {
	test.InitializeTest(t)
	return &DataStore{
		Classes:          make(map[string]Class),
		ClassDefinitions: make(map[string]ClassDefinition),
	}
}

func TestFilterClassChildren(t *testing.T) {
	t.Parallel()

	ds := &DataStore{Classes: map[string]Class{
		"fvTenant": {
			Name: testClassName("fvTenant"),
			Children: []*ClassName{
				testClassName("fvAp"),
				testClassName("unsupportedChild"),
			},
		},
		"fvAp": {Name: testClassName("fvAp")},
	}}

	ds.filterClassChildren()

	assert.Equal(t, []string{"fvAp"}, classNamesToStrings(ds.Classes["fvTenant"].Children))
}

func TestGetResourceName(t *testing.T) {
	t.Parallel()

	resourceClassName := testClassName("fvTenant")
	ds := &DataStore{
		Classes: map[string]Class{
			"fvTenant": {
				Name:         resourceClassName,
				ResourceName: "tenant",
				Artifacts:    []ArtifactEnum{ResourceArtifact},
			},
		},
		ClassDefinitions: map[string]ClassDefinition{
			"fvTenant":   {ResourceName: "ignored_definition_name"},
			"vzBrCP":     {ResourceName: "contract"},
			"fabricNode": {ResourceName: "fabric_node"},
			"fvABDPol":   {ResourceName: "bridge_domain"},
		},
	}

	assert.Equal(t, "tenant", ds.getResourceName("fvTenant"))
	assert.Equal(t, "contract", ds.getResourceName("vzBrCP"))
	assert.Equal(t, "fabric_node", ds.getResourceName("fabricNode"))
	assert.Equal(t, "bridge_domain", ds.getResourceName("fvABDPol"))
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

func TestRefreshUnsupportedAnnotationClassesNotRequested(t *testing.T) {
	ds := &DataStore{}
	if err := ds.refreshUnsupportedAnnotationClasses(); err != nil {
		t.Fatalf("skip unsupported annotation refresh: %v", err)
	}
	assert.Nil(t, ds.UnsupportedAnnotationClasses)
}

func TestRefreshUnsupportedAnnotationClasses(t *testing.T) {
	t.Setenv(constEnvAnnotationUnsupported, "true")
	client := &http.Client{Transport: test.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path != "/doc/jsonmeta/aci-meta.json" {
			t.Errorf("unexpected metadata path %q", request.URL.Path)
		}
		body := `{
			"classes": {
				"aaaConfig": {"isConfigurable": true, "properties": {}},
				"fvTenant": {"isConfigurable": true, "properties": {"annotation": {}}},
				"nullAnnotation": {"isConfigurable": true, "properties": {"annotation": null}},
				"readOnly": {"isConfigurable": false, "properties": {}}
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
			Request:    request,
		}, nil
	})}

	ds := &DataStore{
		client:   client,
		metaHost: "metadata.example.com",
	}
	if err := ds.refreshUnsupportedAnnotationClasses(); err != nil {
		t.Fatalf("refresh unsupported annotation classes: %v", err)
	}
	assert.Equal(t, []string{"aaaConfig", "nullAnnotation"}, ds.UnsupportedAnnotationClasses)
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
				assert.NoError(t, err, test.MessageUnexpectedError(err))
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

	assert.NoError(t, err, test.MessageUnexpectedError(err))
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
		{
			Name: "test_remote_error",
			Input: retrieveEnvMetaClassesInput{
				EnvValue:       "fvTenant",
				ServerResponse: `not found`,
				ServerStatus:   http.StatusNotFound,
			},
			Expected: retrieveEnvMetaClassesExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(retrieveEnvMetaClassesInput)
			expected := testCase.Expected.(retrieveEnvMetaClassesExpected)
			tempDir := t.TempDir()
			t.Chdir(tempDir)
			if err := os.MkdirAll(constMetaPath, 0o755); err != nil {
				t.Fatalf("create metadata directory: %v", err)
			}

			t.Setenv(constEnvMetaClasses, input.EnvValue)

			ds := &DataStore{
				Classes: make(map[string]Class),
				client: &http.Client{Transport: test.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: input.ServerStatus,
						Status:     http.StatusText(input.ServerStatus),
						Body:       io.NopCloser(strings.NewReader(input.ServerResponse)),
						Header:     make(http.Header),
						Request:    request,
					}, nil
				})},
				metaHost:         "metadata.example.com",
				retrievedClasses: make(map[string]bool),
			}

			err := ds.retrieveEnvMetaClassesFromRemote()

			if expected.Error {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			for _, className := range strings.Split(input.EnvValue, ",") {
				if className == "" {
					continue
				}
				contents, err := os.ReadFile(filepath.Join(constMetaPath, className+".json"))
				assert.NoError(t, err)
				assert.Equal(t, input.ServerResponse, string(contents))
			}
		})
	}
}

func TestRetrieveEnvMetaClassesFromRemoteEmptyEnv(t *testing.T) {
	// Note: Cannot use t.Parallel() with t.Setenv()
	t.Setenv(constEnvMetaClasses, "")

	ds := &DataStore{
		Classes:          make(map[string]Class),
		retrievedClasses: make(map[string]bool),
	}

	err := ds.retrieveEnvMetaClassesFromRemote()

	assert.NoError(t, err, test.MessageUnexpectedError(err))
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
				retrievedClasses: make(map[string]bool),
			}

			err := ds.refreshMetaFiles()

			if expected.Error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
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
	testCases := []test.TestCase{
		{
			Name: "test_success",
			Input: retrieveMetaFileInput{
				ClassName:      "fvTenant",
				ServerResponse: `{"label": "tenant"}`,
				ServerStatus:   http.StatusOK,
			},
			Expected: retrieveMetaFileExpected{Error: false},
		},
		{
			Name: "test_remote_error",
			Input: retrieveMetaFileInput{
				ClassName:      "fvTenant",
				ServerResponse: `not found`,
				ServerStatus:   http.StatusNotFound,
			},
			Expected: retrieveMetaFileExpected{Error: true},
		},
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
			input := testCase.Input.(retrieveMetaFileInput)
			expected := testCase.Expected.(retrieveMetaFileExpected)
			t.Chdir(t.TempDir())
			if err := os.MkdirAll(constMetaPath, 0o755); err != nil {
				t.Fatalf("create metadata directory: %v", err)
			}

			ds := &DataStore{
				Classes: make(map[string]Class),
				client: &http.Client{Transport: test.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: input.ServerStatus,
						Status:     http.StatusText(input.ServerStatus),
						Body:       io.NopCloser(strings.NewReader(input.ServerResponse)),
						Header:     make(http.Header),
						Request:    request,
					}, nil
				})},
				metaHost:         "metadata.example.com",
				retrievedClasses: make(map[string]bool),
			}

			err := ds.retrieveMetaFileFromRemote(input.ClassName)

			if expected.Error {
				assert.Error(t, err)
				if input.ClassName == "fvTenant" {
					_, statErr := os.Stat(filepath.Join(constMetaPath, "fvTenant.json"))
					assert.ErrorIs(t, statErr, os.ErrNotExist)
				}
				return
			}
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			contents, readErr := os.ReadFile(filepath.Join(constMetaPath, input.ClassName+".json"))
			assert.NoError(t, readErr)
			assert.Equal(t, input.ServerResponse, string(contents))
		})
	}
}

func TestRetrieveMetaFileFromRemoteAlreadyRetrieved(t *testing.T) {
	t.Parallel()

	ds := &DataStore{
		Classes:          make(map[string]Class),
		retrievedClasses: map[string]bool{"fvTenant": true},
	}

	// Should skip retrieval and not error
	err := ds.retrieveMetaFileFromRemote("fvTenant")

	assert.NoError(t, err, test.MessageUnexpectedError(err))
}

// TestRetrieveMetaFileFromRemote_HTTPTransportError verifies the
// fmt.Errorf("retrieve meta file for class '%s': %w", ...) error path when
// the HTTP client cannot reach the server.
func TestRetrieveMetaFileFromRemote_HTTPTransportError(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ds := &DataStore{
		Classes: make(map[string]Class),
		client: &http.Client{Transport: test.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("transport failure")
		})},
		metaHost:         "metadata.example.com",
		retrievedClasses: make(map[string]bool),
	}

	err := ds.retrieveMetaFileFromRemote("fvTenant")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retrieve meta file for class 'fvTenant'")
}
