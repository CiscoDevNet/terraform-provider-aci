package utils

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
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
	test.InitializeTest(t)
	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, false)
	assert.NotEmpty(t, filenames, "Expected to get file names from directory, but got empty list")
	assert.Equal(t, len(filenames), 3, fmt.Sprintf("Expected to get 2 file names from directory, but got %d", len(filenames)))
	assert.Contains(t, filenames, constTestFile1WithExtension, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile1WithExtension))
	assert.Contains(t, filenames, constTestFile2WithExtension, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile2WithExtension))
	assert.Contains(t, filenames, constTestFile3, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile3))
	assert.NotContains(t, filenames, constTestDir1, fmt.Sprintf("Expected to not find directory name '%s' in the list, but it was found", constTestDir1))
}

func TestGetFileNamesFromDirectoryWithoutExtension(t *testing.T) {
	test.InitializeTest(t)
	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, true)
	assert.NotEmpty(t, filenames, "Expected to get file names from directory, but got empty list")
	assert.Equal(t, len(filenames), 3, fmt.Sprintf("Expected to get 2 file names from directory, but got %d", len(filenames)))
	assert.Contains(t, filenames, constTestFile1WithoutExtension, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile1WithoutExtension))
	assert.Contains(t, filenames, constTestFile2WithoutExtension, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile2WithoutExtension))
	assert.Contains(t, filenames, constTestFile3, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", constTestFile3))
	assert.NotContains(t, filenames, constTestDir1, fmt.Sprintf("Expected to not find directory name '%s' in the list, but it was found", constTestDir1))
}

func TestUnderscore(t *testing.T) {
	test.InitializeTest(t)

	tests := []map[string]string{
		{"input": "tenant", "expected": "tenant"},
		{"input": "Tenant", "expected": "tenant"},
		{"input": "Tenant1", "expected": "tenant1"},
		{"input": "ApplicationEndpointGroup", "expected": "application_endpoint_group"},
		{"input": "Application Endpoint Group", "expected": "application_endpoint_group"},
	}

	for _, test := range tests {
		genLogger.Info(fmt.Sprintf("Executing: %s' with input '%s' and expected output '%s'", t.Name(), test["input"], test["expected"]))
		result := Underscore(test["input"])
		assert.Equal(t, test["expected"], result, fmt.Sprintf("Expected '%s', but got '%s'", test["expected"], result))
	}

}
