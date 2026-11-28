package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/logger"
)

const FilePath = "gen/managed_files.json"

var genLogger = logger.InitializeLogger()

type Manifest struct {
	Directories map[string][]string `json:"directories"`
}

func New(files []string) *Manifest {
	return &Manifest{Directories: groupManagedFiles(files)}
}

func Load() (*Manifest, error) {
	contents, err := os.ReadFile(FilePath)
	if errors.Is(err, fs.ErrNotExist) {
		genLogger.Debugf("Managed-files manifest does not exist yet: %s.", FilePath)
		return &Manifest{Directories: map[string][]string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read managed-files manifest %q: %w", FilePath, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return nil, fmt.Errorf("decode managed-files manifest %q: %w", FilePath, err)
	}

	fileCount := 0
	for _, files := range manifest.Directories {
		fileCount += len(files)
	}
	genLogger.Debugf("Successfully loaded %d managed files from: %s.", fileCount, FilePath)
	return &manifest, nil
}

func (m *Manifest) DeleteManagedFiles() error {
	var deletionErrors []error
	for outputDirectory, files := range m.Directories {
		for _, fileName := range files {
			managedPath := pathpkg.Join(outputDirectory, fileName)
			err := os.Remove(filepath.FromSlash(managedPath))
			if errors.Is(err, fs.ErrNotExist) {
				genLogger.Debugf("Skipping already absent managed file: %s.", managedPath)
				continue
			}
			if err != nil {
				deletionErrors = append(deletionErrors, fmt.Errorf("delete managed file %q: %w", managedPath, err))
				continue
			}
			genLogger.Debugf("Successfully deleted managed file: %s.", managedPath)
		}
	}

	return errors.Join(deletionErrors...)
}

func (m *Manifest) Write() error {
	contents, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encode managed-files manifest: %w", err)
	}
	contents = append(contents, '\n')
	if err := os.WriteFile(FilePath, contents, 0o644); err != nil {
		return fmt.Errorf("write managed-files manifest %q: %w", FilePath, err)
	}

	fileCount := 0
	for _, files := range m.Directories {
		fileCount += len(files)
	}
	genLogger.Debugf("Successfully wrote %d managed files to: %s.", fileCount, FilePath)
	return nil
}

func groupManagedFiles(files []string) map[string][]string {
	managedFiles := append([]string{}, files...)
	sort.Strings(managedFiles)
	directories := make(map[string][]string)

	for _, managedPath := range managedFiles {
		outputDirectory := pathpkg.Dir(managedPath)
		directories[outputDirectory] = append(directories[outputDirectory], pathpkg.Base(managedPath))
	}

	return directories
}
