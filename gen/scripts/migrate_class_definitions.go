//go:build ignore
// +build ignore

// quick AI generated script for testing with migrated files
// GENERATE NEW FORMAT DEFINTIONS: go run gen/scripts/migrate_class_definitions.go
// REMOVE NEW FORMAT DEFINTIONS:go run gen/scripts/migrate_class_definitions.go clean
// Alternatively leverage git for removal

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Input struct matching the current top-level YAML keys in classes/ files.
type OldClassDefinition struct {
	AllowDelete     string   `yaml:"allow_delete"`
	ExcludeChildren []string `yaml:"exclude_children"`
	IncludeChildren []string `yaml:"include_children"`
	SubCategory     string   `yaml:"sub_category"`
	UiLocations     []string `yaml:"ui_locations"`
}

// Output struct matching the new ClassDefinition format.
type NewClassDefinition struct {
	AllowDelete     string             `yaml:"allow_delete,omitempty"`
	Documentation   DocumentationBlock `yaml:"documentation,omitempty"`
	ExcludeChildren []string           `yaml:"exclude_children,omitempty"`
	IncludeChildren []string           `yaml:"include_children,omitempty"`
}

type DocumentationBlock struct {
	SubCategory string   `yaml:"sub_category,omitempty"`
	UiLocations []string `yaml:"ui_locations,omitempty"`
}

func main() {
	classesDir := "gen/definitions/classes"
	outputDir := "gen/definitions"

	// If "clean" argument is passed, remove all .yaml files except global.yaml from the definitions folder.
	if len(os.Args) > 1 && os.Args[1] == "clean" {
		cleanFiles, err := filepath.Glob(filepath.Join(outputDir, "*.yaml"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading definitions dir: %v\n", err)
			os.Exit(1)
		}
		removed := 0
		for _, f := range cleanFiles {
			if filepath.Base(f) == "global.yaml" {
				continue
			}
			if err := os.Remove(f); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing %s: %v\n", f, err)
				continue
			}
			fmt.Printf("Removed: %s\n", f)
			removed++
		}
		fmt.Printf("\nDone: %d files removed\n", removed)
		return
	}

	files, err := filepath.Glob(filepath.Join(classesDir, "*.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading classes dir: %v\n", err)
		os.Exit(1)
	}

	migrated := 0
	skipped := 0

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			continue
		}

		var old OldClassDefinition
		if err := yaml.Unmarshal(data, &old); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file, err)
			continue
		}

		// Skip files that have no fields to migrate.
		if old.AllowDelete == "" && old.SubCategory == "" && len(old.UiLocations) == 0 && len(old.ExcludeChildren) == 0 && len(old.IncludeChildren) == 0 {
			skipped++
			continue
		}

		// Map allow_delete: false to "never" to match the expected string format.
		allowDelete := old.AllowDelete
		if allowDelete == "false" {
			allowDelete = "never"
		}

		newDef := NewClassDefinition{
			AllowDelete:     allowDelete,
			ExcludeChildren: old.ExcludeChildren,
			IncludeChildren: old.IncludeChildren,
			Documentation: DocumentationBlock{
				SubCategory: old.SubCategory,
				UiLocations: old.UiLocations,
			},
		}

		out, err := yaml.Marshal(newDef)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling %s: %v\n", file, err)
			continue
		}

		className := strings.TrimSuffix(filepath.Base(file), ".yaml")
		outputPath := filepath.Join(outputDir, className+".yaml")

		if err := os.WriteFile(outputPath, out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			continue
		}

		fmt.Printf("Migrated: %s -> %s\n", file, outputPath)
		migrated++
	}

	fmt.Printf("\nDone: %d migrated, %d skipped (no documentation fields)\n", migrated, skipped)
}
