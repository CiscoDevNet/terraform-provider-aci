package manifest

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteAndLoadManifest(t *testing.T) {
	useTemporaryRepository(t)
	managedFiles := newTestManifest(t, []string{
		"internal/provider/models/fvTenant.go",
		"internal/provider/models/fvAp.go",
	})

	if err := managedFiles.Write(); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	contents, err := os.ReadFile(FilePath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	expectedContents := "{\n  \"directories\": {\n    \"internal/provider/models\": [\n      \"fvAp.go\",\n      \"fvTenant.go\"\n    ]\n  }\n}\n"
	if string(contents) != expectedContents {
		t.Fatalf("manifest contents %q, expected %q", contents, expectedContents)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	expectedDirectories := map[string][]string{
		"internal/provider/models": {"fvAp.go", "fvTenant.go"},
	}
	if !reflect.DeepEqual(loaded.Directories, expectedDirectories) {
		t.Fatalf("loaded directories %#v, expected %#v", loaded.Directories, expectedDirectories)
	}
}

func TestLoadMissingManifest(t *testing.T) {
	useTemporaryRepository(t)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load missing manifest: %v", err)
	}
	if loaded == nil || len(loaded.Directories) != 0 {
		t.Fatalf("expected empty manifest, got %#v", loaded)
	}
}

func TestWriteEmptyManifestUsesObject(t *testing.T) {
	useTemporaryRepository(t)
	managedFiles := newTestManifest(t, nil)

	if err := managedFiles.Write(); err != nil {
		t.Fatalf("write empty manifest: %v", err)
	}
	contents, err := os.ReadFile(FilePath)
	if err != nil {
		t.Fatalf("read empty manifest: %v", err)
	}
	if string(contents) != "{\n  \"directories\": {}\n}\n" {
		t.Fatalf("unexpected empty manifest contents: %q", contents)
	}
}

func TestDeleteManagedFilesSkipsAlreadyAbsentFiles(t *testing.T) {
	useTemporaryRepository(t)
	modelDirectory := filepath.Join("internal", "provider", "models")
	if err := os.MkdirAll(modelDirectory, 0o755); err != nil {
		t.Fatalf("create model directory: %v", err)
	}
	existingPath := "internal/provider/models/fvTenant.go"
	if err := os.WriteFile(filepath.FromSlash(existingPath), []byte("generated"), 0o600); err != nil {
		t.Fatalf("write managed file: %v", err)
	}

	managedFiles := newTestManifest(t, []string{
		existingPath,
		"internal/provider/models/alreadyDeleted.go",
	})
	if err := managedFiles.DeleteManagedFiles(); err != nil {
		t.Fatalf("delete managed files: %v", err)
	}
	if _, err := os.Stat(filepath.FromSlash(existingPath)); !os.IsNotExist(err) {
		t.Fatalf("expected managed file to be deleted, got error %v", err)
	}
}

func TestLoadManifestRejectsInvalidJSON(t *testing.T) {
	useTemporaryRepository(t)
	if err := os.WriteFile(FilePath, []byte(`{"directories": {`), 0o600); err != nil {
		t.Fatalf("write manifest fixture: %v", err)
	}
	if _, err := Load(); err == nil {
		t.Fatal("expected invalid manifest error")
	}
}

func TestManifestGroupsFilesByContainingDirectory(t *testing.T) {
	useTemporaryRepository(t)
	managedFiles := newTestManifest(t, []string{
		"internal/provider/models/fvTenant.go",
		"docs/resources/nested/tenant.md",
	})

	if err := managedFiles.Write(); err != nil {
		t.Fatalf("write multi-directory manifest: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("load multi-directory manifest: %v", err)
	}
	expectedDirectories := map[string][]string{
		"docs/resources/nested":    {"tenant.md"},
		"internal/provider/models": {"fvTenant.go"},
	}
	if !reflect.DeepEqual(loaded.Directories, expectedDirectories) {
		t.Fatalf("loaded directories %#v, expected %#v", loaded.Directories, expectedDirectories)
	}
}

func newTestManifest(t *testing.T, files []string) *Manifest {
	t.Helper()
	return New(files)
}

func useTemporaryRepository(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(filepath.Dir(FilePath), 0o755); err != nil {
		t.Fatalf("create manifest directory: %v", err)
	}
}
