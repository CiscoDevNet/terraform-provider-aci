package utils

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

const (
	testDirectoryForGetFileNamesFromDirectory = "./test_data/test_get_file_names_from_directory"
	testFile1                                 = "file1.json"
	testFile2                                 = "file2.json"
	testDir1                                  = "dir1"
)

func TestGetFileNamesFromDirectory(t *testing.T) {
	test.InitializeTest(t, 0)

	filenames := GetFileNamesFromDirectory(testDirectoryForGetFileNamesFromDirectory)
	assert.NotEmpty(t, filenames, "Expected to get file names from directory, but got empty list")
	assert.Equal(t, len(filenames), 2, fmt.Sprintf("Expected to get 2 file names from directory, but got %d", len(filenames)))
	assert.Contains(t, filenames, testFile1, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", testFile1))
	assert.Contains(t, filenames, testFile2, fmt.Sprintf("Expected to find file name '%s' in the list, but it was not found", testFile2))
	assert.NotContains(t, filenames, testDir1, fmt.Sprintf("Expected to not find directory name '%s' in the list, but it was found", testDir1))

}
