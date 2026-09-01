# Generated output ownership

The generator is intended to produce the complete generated output set in a
single run. It must remove the previous generated output set before writing
the current set, while preserving handwritten and explicitly excluded files.

## Current inventory

The following output areas contain files produced by, or historically
associated with, `gen/generator.go`:

| Output area | Current inventory | Cleanup status |
| --- | ---: | --- |
| `internal/provider/` | 505 files | Managed output with explicit preserved paths |
| `internal/provider/models/` | 176 files | Managed by the active generator manifest |
| `internal/custom_types/` | 36 files | Managed output with explicit preserved paths |
| `docs/` | 255 files | Managed output with explicit preserved paths |
| `gen/testvars/` | 125 files | Managed output |
| `legacy-docs/` | 4 generated-marked files | Not managed by the active generator |

The generated marker is not, by itself, an ownership declaration. Some legacy
files with the marker are intentionally preserved, while some handwritten
files have no marker. Ownership is determined by the manifest and the
preserved-path policy.

## Preserved legacy paths

The following paths are outside the active generator manifest and must not be
removed during cleanup:

```text
internal/provider/provider_test.go
internal/provider/utils.go
internal/provider/test_constants.go
internal/provider/resource_aci_rest_managed.go
internal/provider/resource_aci_rest_managed_test.go
internal/provider/data_source_aci_rest_managed.go
internal/provider/data_source_aci_rest_managed_test.go
internal/provider/annotation_unsupported.go
internal/provider/data_source_aci_system.go
internal/provider/data_source_aci_system_test.go
internal/provider/function_compare_versions.go
internal/provider/function_compare_versions_test.go
internal/custom_types/ipAddress.go
internal/custom_types/arpLearning.go
docs/data-sources/system.md
examples/provider/provider.tf
examples/data-sources/aci_system/*
```

The list includes files that contain the legacy generated marker. They remain
preserved until their ownership is explicitly migrated to an active template.

## Cleanup contract

The active generator must use a manifest containing the files owned and
produced by the previous run. The manifest is the authoritative deletion
target. It allows files for removed classes or artifacts to be deleted without
inspecting every file in shared directories.

The manifest is stored at:

```text
gen/managed_files.json
```

It groups exact managed files by their containing directory and stores only the
filename beneath each directory key:

```json
{
  "directories": {
    "internal/provider/models": [
      "fvAp.go",
      "fvTenant.go"
    ]
  }
}
```

Directory keys and filenames are derived directly from rendered output paths
and written in deterministic sorted order. The manifest is generator
bookkeeping, not provider output, and is never included in its own deletion
list.

The manifest must not contain preserved legacy paths. A one-time migration or
recovery path may use the legacy preserve list and generated markers to
bootstrap or repair the manifest, but normal generation uses the manifest
directly.

`gen/managed_files.json` and the preserved-path policy are temporary migration
controls. They may be removed or simplified after all legacy artifacts have
been migrated and generated output is isolated into dedicated directories.

Render jobs are validated before managed files are deleted. Validation checks
that every referenced template is loaded and output paths are unique.
Generation then deletes all paths from the previous manifest before writing
the current output set. A managed file that is already absent is skipped;
other deletion errors stop generation.

`template.Generator` owns the canonical DataStore, loaded templates, render
jobs, managed-file cleanup, concurrent file rendering, and replacement
manifest. The outer generator only initializes the DataStore, constructs the
template generator, and calls `Generate()`.

`Generator.renderTemplate` is the shared renderer for every template and
output path. It renders one artifact into its own memory buffer, formats it
according to the output extension, validates its marker, creates its output
directory, and writes it directly to its final path. A bounded worker pool
calls this renderer concurrently because every artifact owns a unique path
and the normalized DataStore is read-only during rendering. The replacement
manifest is written only after every worker succeeds. If any worker fails,
every planned output path is removed and the previous manifest is retained.
Existing generated outputs may therefore be absent after a failed run, but
files written by an unsuccessful attempt cannot become unmanaged. Rerunning
the generator recreates the complete output set.

Generated Go is validated and formatted with `go/format`. Terraform
configuration is validated with the HCL parser and formatted with `hclwrite`.
Markdown has no canonical formatter and is therefore written exactly as
rendered by its template.

## Template ownership

The former templates are stored in `gen/templates/legacy_templates/` for
comparison and migration reference. The active templates will be stored in
`gen/templates/`.

Each active template declares its output location centrally, so the generator
can account for every expected file. Each active template must also emit a
format-valid generated marker as the first output line:

| Output format | Marker form |
| --- | --- |
| Go | `// Code generated by "gen/generator.go"; DO NOT EDIT.` |
| Terraform configuration | `# Configuration generated by "gen/generator.go"; DO NOT EDIT.` |
| Markdown | `<!-- Documentation generated by "gen/generator.go"; DO NOT EDIT. -->` |

The manifest remains authoritative for cleanup; the marker provides human
visibility and a safety check that a rendered file came from an active
template. Template-specific front matter must be tested to ensure the marker
does not invalidate the document format.

Every template executes with an internally constructed standard context
containing references to the current resolved class and the canonical
DataStore:

```go
type TemplateContext struct {
	Class     *data.Class
	DataStore *data.DataStore
}
```

This context does not duplicate or reshape datastore fields. Class-specific
information is read from `Class`; genuinely global or cross-class information
remains available from `DataStore`.
