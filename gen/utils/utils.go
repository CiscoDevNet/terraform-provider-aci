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

// A simplified version of a function to create plural forms for words.
// This function is not comprehensive and is not intended to cover all pluralization rules in English.
// It is intended for basic use cases and does not handle irregular plurals like "mouse" -> "mice".
// If we need a more robust pluralization solution, we should consider using a external package or expanding the logic.
func Plural(notPlural string) (string, error) {
	if strings.HasSuffix(notPlural, "y") {
		return fmt.Sprintf("%s%s", strings.TrimSuffix(notPlural, "y"), "ies"), nil
	} else if !strings.HasSuffix(notPlural, "s") {
		return fmt.Sprintf("%ss", notPlural), nil
	}
	genLogger.Error(fmt.Sprintf("The word '%s' has no plural rule defined.", notPlural))
	return "", fmt.Errorf("the word '%s' has no plural rule defined", notPlural)
}

// Reused from https://github.com/buxizhizhoum/inflection/blob/master/inflection.go#L8 to avoid importing the whole package
func Underscore(s string) string {
	for _, reStr := range []string{`([A-Z]+)([A-Z][a-z])`, `([a-z\d])([A-Z])`} {
		re := regexp.MustCompile(reStr)
		s = re.ReplaceAllString(s, "${1}_${2}")
	}
	s = strings.ReplaceAll(strings.ToLower(s), " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}
