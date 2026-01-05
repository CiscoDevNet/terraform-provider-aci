package utils

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	constTestDirectoryForGetFileNamesFromDirectory = "./test/data/test_get_file_names_from_directory"
	constTestFile1WithoutExtension                 = "file1"
	constTestFile1WithExtension                    = "file1.json"
	constTestFile2WithoutExtension                 = "file2"
	constTestFile2WithExtension                    = "file2.json"
	constTestFile3                                 = "file3"
	constTestDir1                                  = "dir1"
)

func TestGetFileNamesFromDirectoryWithExtension(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, false)

	require.NotEmpty(t, filenames, "file names should not be empty")
	assert.Len(t, filenames, 3)
	assert.Contains(t, filenames, constTestFile1WithExtension)
	assert.Contains(t, filenames, constTestFile2WithExtension)
	assert.Contains(t, filenames, constTestFile3)
	assert.NotContains(t, filenames, constTestDir1)
}

func TestGetFileNamesFromDirectoryWithoutExtension(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, true)

	require.NotEmpty(t, filenames, "file names should not be empty")
	assert.Len(t, filenames, 3)
	assert.Contains(t, filenames, constTestFile1WithoutExtension)
	assert.Contains(t, filenames, constTestFile2WithoutExtension)
	assert.Contains(t, filenames, constTestFile3)
	assert.NotContains(t, filenames, constTestDir1)
}

func TestUnderscore(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "lowercase", input: "tenant", expected: "tenant"},
		{name: "capitalized", input: "Tenant", expected: "tenant"},
		{name: "with_number", input: "Tenant1", expected: "tenant1"},
		{name: "camel_case", input: "ApplicationEndpointGroup", expected: "application_endpoint_group"},
		{name: "space_separated", input: "Application Endpoint Group", expected: "application_endpoint_group"},
		{name: "hyphen_separated", input: "Application-Endpoint-Group", expected: "application_endpoint_group"},
		{name: "mixed_separators", input: "Application Endpoint-Group", expected: "application_endpoint_group"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := Underscore(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPlural(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "policy_to_policies", input: "monitor_policy", expected: "monitor_policies"},
		{name: "add_s", input: "annotation", expected: "annotations"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := Plural(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPluralError(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	_, err := Plural("contracts")
	require.Error(t, err)
	assert.ErrorContains(t, err, "no plural rule defined")
}

func TestUnderscoreEdgeCases(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty_string", input: "", expected: ""},
		{name: "single_lowercase", input: "a", expected: "a"},
		{name: "single_uppercase", input: "A", expected: "a"},
		{name: "numbers_only", input: "123", expected: "123"},
		{name: "leading_number", input: "1Tenant", expected: "1_tenant"},
		{name: "underscore_input", input: "already_snake", expected: "already_snake"},
		{name: "multiple_underscores", input: "a__b", expected: "a__b"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := Underscore(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetFileNamesFromDirectoryNonExistent(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory("./non_existent_directory", false)

	assert.Empty(t, filenames)
}
