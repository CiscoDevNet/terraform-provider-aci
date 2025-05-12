package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

// Initialize a logger instance for the generator.
var genLogger = logger.InitializeLogger()

func GetFileNamesFromDirectory(path string, removeExtension bool) []string {
	genLogger.Debug(fmt.Sprintf("Getting file names from directory: %s.", path))
	var names []string
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		for _, entry := range entries {
			// Check if the entry is not a directory and append its name to the list.
			if !entry.IsDir() {
				name := entry.Name()
				// When removeExtension is true, remove the file extension from the name. (ex file1.json -> file1)
				// When removeExtension is false, keep the file extension in the name. (ex file1.json -> file1.json)
				if removeExtension {
					name = strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
				}
				names = append(names, name)
			}
		}
	}
	genLogger.Debug(fmt.Sprintf("The directory '%s' contains the file names: %s.", path, names))
	return names
}

// Reused from https://github.com/buxizhizhoum/inflection/blob/master/inflection.go#L8 to avoid importing the whole package
func Underscore(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1}_${2}")
	}
	return strings.ReplaceAll(strings.ToLower(s), " ", "_")
}
