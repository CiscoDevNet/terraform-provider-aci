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
	genLogger.Debugf("Getting file names from directory: %s.", path)
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
	genLogger.Debugf("The directory '%s' contains the file names: %s.", path, names)
	return names
}

// A simplified version of a function to create plural forms for words.
// This function is not comprehensive and is not intended to cover all pluralization rules in English.
// It is intended for basic use cases and does not handle irregular plurals like "mouse" -> "mice".
// If we need a more robust pluralization solution, we should consider using a external package or expanding the logic.
func Plural(word string) string {
	if strings.HasSuffix(word, "y") {
		word = fmt.Sprintf("%sies", strings.TrimSuffix(word, "y"))
	} else if !strings.HasSuffix(word, "s") {
		word = fmt.Sprintf("%ss", word)
	}
	return word
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

// GetValueFromMapWithOverride returns override when it is not the zero value of T.
// Otherwise it returns the value stored under key in m, type-asserted to T. Returns the
// zero value of T when m is nil, the key is missing, or the stored value is not a T.
//
// The "non-zero override wins" precedence fits string-shaped overrides (empty string means
// "unset") and any pointer/interface T (nil means "unset"). For bool or numeric T the zero
// value (false / 0) cannot be expressed as a meaningful override, so a different mechanism
// is required for those cases.
//
// Typical use: merging a per-field override from a definition file with a value read from
// a JSON-unmarshalled meta sub-map (map[string]any).
func GetValueFromMapWithOverride[T comparable](m map[string]any, key string, override T) T {
	var zero T
	if override != zero {
		return override
	}
	if m == nil {
		return zero
	}
	v, _ := m[key].(T)
	return v
}
