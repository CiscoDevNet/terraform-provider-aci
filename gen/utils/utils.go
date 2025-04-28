package utils

import "os"

func GetFileNamesFromDirectory(path string) []string {
	var names []string
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
	}
	return names
}
