package utils

import (
	"fmt"
	"os"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitalizeLogger()

func GetFileNamesFromDirectory(path string) []string {
	genLogger.Debug(fmt.Sprintf("Getting file names from directory: %s.", path))
	var names []string
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		for _, entry := range entries {
			// Check if the entry is not a directory and append its name to the list.
			if !entry.IsDir() {
				names = append(names, entry.Name())
			}
		}
	}
	genLogger.Debug(fmt.Sprintf("The directory '%s' contains the file names: %s.", path, names))
	return names
}
