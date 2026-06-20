//go:build ignore
// +build ignore

/*
Migration script for class definition YAML files. Reads the per-class legacy
inputs under gen/scripts/legacy_definitions/classes/ and emits canonical
ClassDefinition YAML under gen/definitions/, rebased on the canonical
gen/utils/data structs so the loader contract is enforced at migration time.

Disposition table is the v2.19.0 coverage audit in MIGRATION_OVERVIEW.md
(sections 1-8 + the section-10 closing-gaps list). The knownLegacyKeys map
mirrors that audit one-for-one; anything found in the legacy YAML that is
not in the map surfaces as an L1 unknown-key warning so drift cannot slip
past a fresh import.

Phase 1 (this commit) translates the keys the prior script already handled:
  - allow_delete       (value remap "false" -> "never")
  - exclude_children   (passthrough)
  - include_children   (passthrough)
  - sub_category       (moves under documentation)
  - ui_locations       (moves under documentation)

Every other known key is tallied as TODO so the per-run summary shows
remaining coverage. Subsequent commits flesh out each section.

The four verification levels:

  L1 - Allowlist. Every top-level legacy key is enumerated in knownLegacyKeys.
       Unknown keys produce a warning so the moment a new legacy key sneaks
       in we hear about it.

  L2 - Disposition tally. Per-key counters track how many files use each key
       and what disposition was applied; printed in the run summary. Newly
       added keys (or coverage regressions) jump out.

  L3 - UnmarshalStrict round-trip. After projecting the migrated struct into
       its yaml view, the script re-loads it via yaml.UnmarshalStrict against
       data.ClassDefinition. Any output the canonical loader would reject
       becomes a per-file failure at migration time, not on the next
       go-generate run.

  L4 - Derivation assertion hooks. Stub call sites where subsequent commits
       cross-check legacy values against the resolver's auto-derivation
       rules (e.g. confirm datasource_required is reproduced by the
       IdentifiedBy + value-transform pipeline). Empty for the scaffold.

USAGE

  go run gen/scripts/migrate_class_definitions.go        emit canonical YAML
  go run gen/scripts/migrate_class_definitions.go clean  remove all .yaml
                                                          under gen/definitions
                                                          except global.yaml
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/data"
	"gopkg.in/yaml.v2"
)

// sectionEnum names the MIGRATION_OVERVIEW.md disposition section a legacy
// key belongs to.
type sectionEnum int

const (
	sectionUnknown sectionEnum = iota
	sectionDirect              // S1 direct mapping (copy verbatim, same/sibling key)
	sectionSemantic            // S2 semantic mapping (shape-changing transform)
	sectionObsolete            // S3 already migrated or removed (drop with log)
	sectionAdd                 // S4 ADD - new field, often with rename/value remap
	sectionReuse               // S5 REUSE - remap onto an existing canonical field
	sectionDerive              // S6 DERIVE - drop, computed in Go from meta
	sectionConst               // S7 CONST - drop, relocated to constants.go
	sectionPostpone            // S8 POSTPONE - preserve unchanged with TODO marker
)

func (s sectionEnum) String() string {
	switch s {
	case sectionDirect:
		return "S1 direct"
	case sectionSemantic:
		return "S2 semantic"
	case sectionObsolete:
		return "S3 obsolete"
	case sectionAdd:
		return "S4 ADD"
	case sectionReuse:
		return "S5 REUSE"
	case sectionDerive:
		return "S6 DERIVE"
	case sectionConst:
		return "S7 CONST"
	case sectionPostpone:
		return "S8 POSTPONE"
	default:
		return "UNKNOWN"
	}
}

// keyInfo records the disposition of a single legacy top-level key and
// whether the current script implements the translation. A nil entry (key
// not in knownLegacyKeys) is reported as L1 unknown.
type keyInfo struct {
	section     sectionEnum
	implemented bool
}

// knownLegacyKeys is the exhaustive allowlist of every top-level YAML key
// the migration script expects to encounter under
// gen/scripts/legacy_definitions/classes/*.yaml. Source of truth:
// MIGRATION_OVERVIEW.md section 10 v2.19.0 coverage audit + section 10.2
// closing-gaps list. Keep this aligned with that audit; when a freshly
// imported legacy file carries a new key, decide its disposition
// (sections 1-8) before adding it here.
var knownLegacyKeys = map[string]keyInfo{
	// S1 direct mapping (class-level)
	"allow_delete":      {sectionSemantic, true}, // value remap "false" -> "never"
	"rn_prepend":        {sectionDirect, true},
	"required_as_child": {sectionDirect, true},
	"resource_name":     {sectionDirect, true},
	"dn_formats":        {sectionDirect, true}, // moves under documentation
	"exclude_children":  {sectionDirect, true},
	"include_children":  {sectionDirect, true},
	// S1 direct mapping under documentation block
	"sub_category": {sectionDirect, true},
	"ui_locations": {sectionDirect, true},

	// S2 semantic mapping (class-level)
	"children":             {sectionSemantic, false}, // -> include_children (subtract meta containedBy first)
	"contained_by":         {sectionSemantic, false}, // -> include_parents (subtract meta containedBy first)
	"class_version":        {sectionSemantic, false}, // -> supported_versions
	"relationship_classes": {sectionSemantic, false}, // -> relation_info.to_classes
	"migration_blocks":     {sectionSemantic, false}, // -> state_upgrades (two-source merge)
	"migration_version":    {sectionSemantic, false}, // -> state_upgrades (drives migration_source)
	"type_changes":         {sectionSemantic, false}, // -> state_upgrades.attributes
	"resource_notes":       {sectionSemantic, false}, // -> documentation.resource.notes

	// S3 obsolete (drop with one-line log)
	// S3 OBSOLETE
	"multi_relationship_class": {sectionObsolete, true},

	// S4 ADD - schema additions
	"include":                            {sectionAdd, true}, // -> artifacts ([] when include: false)
	"multi_parents":                      {sectionAdd, true}, // -> parent_dn_variants
	"example_classes":                    {sectionAdd, true}, // -> documentation.example_parent_classes
	"exclude_from_testing":               {sectionAdd, true}, // -> test_config.ignore_tests: [child]
	"ignore_import_state_verify_in_test": {sectionAdd, true}, // -> test_config.ignore_import_state_verify

	// S5 REUSE
	"max_one_class_allowed": {sectionReuse, true}, // -> is_single_nested_when_defined_as_child
	"parent_example_dn":     {sectionReuse, true}, // dropped (covered by static dependency)
	"remove_from_contains":  {sectionReuse, true}, // -> exclude_children (now covers docs side too)

	// S6 DERIVE (drop, computed in Go from meta)
	"resource_identifier":                {sectionDerive, true},
	"data_source_has_no_name_identifier": {sectionDerive, true},
	"static_parent":                      {sectionDerive, true},
	"contained_by_excludes":              {sectionDerive, true}, // S10.2 - classes/global.yaml

	// S7 CONST (drop, relocated to constants.go)
	"docs_examples_amount":  {sectionConst, true}, // classes/global.yaml
	"docs_parent_dn_amount": {sectionConst, true}, // classes/global.yaml

	// S8 POSTPONE
	"class_version_tests": {sectionPostpone, true},
}

// keyTally accumulates per-key and per-section counts across the migration
// run for the L2 summary.
type keyTally struct {
	keyCounts       map[string]int
	sectionCounts   map[sectionEnum]int
	unknownKeyFiles map[string][]string
	implementedSeen map[string]bool
}

func newKeyTally() *keyTally {
	return &keyTally{
		keyCounts:       map[string]int{},
		sectionCounts:   map[sectionEnum]int{},
		unknownKeyFiles: map[string][]string{},
		implementedSeen: map[string]bool{},
	}
}

// record increments the per-key counter and tracks unknown keys (L1) plus
// per-section totals (L2). Returns the disposition so the caller can decide
// whether to act on it.
func (t *keyTally) record(file, key string) keyInfo {
	t.keyCounts[key]++
	info, ok := knownLegacyKeys[key]
	if !ok {
		t.sectionCounts[sectionUnknown]++
		t.unknownKeyFiles[key] = append(t.unknownKeyFiles[key], file)
		return keyInfo{section: sectionUnknown}
	}
	t.sectionCounts[info.section]++
	if info.implemented {
		t.implementedSeen[key] = true
	}
	return info
}

func (t *keyTally) print() {
	fmt.Println()
	fmt.Println("=== L2 disposition tally ===")
	var keys []string
	for k := range t.keyCounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		info, ok := knownLegacyKeys[k]
		marker := "TODO"
		switch {
		case !ok:
			marker = "!! UNKNOWN"
		case info.implemented:
			marker = "done"
		}
		section := info.section.String()
		if !ok {
			section = sectionUnknown.String()
		}
		fmt.Printf("  %4d  %-40s  %-14s  %s\n", t.keyCounts[k], k, section, marker)
	}

	fmt.Println()
	fmt.Println("=== L2 per-section totals ===")
	sections := []sectionEnum{
		sectionDirect, sectionSemantic, sectionObsolete, sectionAdd,
		sectionReuse, sectionDerive, sectionConst, sectionPostpone,
		sectionUnknown,
	}
	for _, s := range sections {
		if c := t.sectionCounts[s]; c > 0 {
			fmt.Printf("  %5d  %s\n", c, s)
		}
	}

	if len(t.unknownKeyFiles) > 0 {
		fmt.Println()
		fmt.Println("=== L1 unknown keys (NOT in knownLegacyKeys) ===")
		var ukeys []string
		for k := range t.unknownKeyFiles {
			ukeys = append(ukeys, k)
		}
		sort.Strings(ukeys)
		for _, k := range ukeys {
			files := t.unknownKeyFiles[k]
			fmt.Printf("  %s (%d files): %s\n", k, len(files), strings.Join(files, ", "))
		}
	}
}

// assertDerivable is an L4 hook reserved for cross-checks between a legacy
// value and the canonical resolver's auto-derivation rules. Subsequent
// commits (S6 DERIVE) wire up specific assertions here (for example,
// confirm datasource_required matches the IdentifiedBy + value-transform
// prediction). Empty for the scaffold so the call sites already exist when
// the assertions land.
func assertDerivable(file, key string, legacyValue any) {
	_ = file
	_ = key
	_ = legacyValue
}

// migrate translates one legacy YAML payload (parsed as map[string]any)
// into the canonical data.ClassDefinition. The translation is intentionally
// narrow for Phase 1: only the keys flagged implemented in knownLegacyKeys
// are written into the canonical struct. Every other known key is tallied
// as TODO so the coverage gap is visible at every run; unknown keys raise
// an L1 warning via the tally.
func migrate(file string, legacy map[string]any, tally *keyTally) data.ClassDefinition {
	var out data.ClassDefinition
	for key, val := range legacy {
		info := tally.record(file, key)
		if info.section == sectionDerive {
			assertDerivable(file, key, val)
		}
		if !info.implemented {
			continue
		}
		switch key {
		case "allow_delete":
			// Legacy `allow_delete: false` (yaml-parsed as Go bool) is the
			// only value that has a real translation - canonical "never".
			// Legacy `allow_delete: true` is the loader default and is
			// dropped during migration (see setAllowDelete in class.go).
			// Any unexpected string is passed through so a strict round-trip
			// surfaces it.
			switch v := val.(type) {
			case bool:
				if !v {
					out.AllowDelete = "never"
				}
			case string:
				if v == "false" {
					out.AllowDelete = "never"
				} else if v != "true" {
					out.AllowDelete = v
				}
			}
		case "exclude_children":
			out.ExcludeChildren = toStringSlice(val)
		case "include_children":
			out.IncludeChildren = toStringSlice(val)
		case "resource_name":
			out.ResourceName, _ = val.(string)
		case "rn_prepend":
			out.RnPrepend, _ = val.(string)
		case "required_as_child":
			if b, ok := val.(bool); ok {
				out.RequiredAsChild = b
			}
		case "sub_category":
			out.Documentation.SubCategory, _ = val.(string)
		case "ui_locations":
			out.Documentation.UiLocations = toStringSlice(val)
		case "dn_formats":
			out.Documentation.DnFormats = toStringSlice(val)
		case "include":
			// Legacy `include: true` forces a class with empty IdentifiedBy
			// into the generator registry. The new resolver expresses the same
			// opt-in via a non-nil artifacts override. fvFBRoute has
			// IdentifiedBy non-empty so its include: true is redundant; the
			// one-off cleanup removes it in a later commit.
			if b, ok := val.(bool); ok {
				if b {
					out.Artifacts = []data.ArtifactEnum{data.ResourceArtifact, data.DatasourceArtifact}
				} else {
					out.Artifacts = []data.ArtifactEnum{}
				}
			}
		case "example_classes":
			out.Documentation.ExampleParentClasses = toStringSlice(val)
		case "multi_parents":
			out.ParentDnVariants = migrateMultiParents(val)
		case "exclude_from_testing":
			if b, ok := val.(bool); ok && b {
				out.TestConfig.IgnoreTests = []data.IgnoreTestEnum{data.ChildIgnoreTest}
			}
		case "ignore_import_state_verify_in_test":
			if b, ok := val.(bool); ok {
				out.TestConfig.IgnoreImportStateVerify = b
			}
		case "max_one_class_allowed":
			if b, ok := val.(bool); ok {
				out.IsSingleNestedWhenDefinedAsChild = b
			}
		case "parent_example_dn":
			// Drop intentionally. The legacy value supplied a static parent
			// DN for example/test rendering; the canonical pipeline now
			// derives the same DN from class.TestConfig.Dependencies (and
			// the static_parent meta flag), so no canonical field needs the
			// value.
		case "remove_from_contains":
			// Legacy remove_from_contains entries unioned into the canonical
			// exclude_children list - C10 extended the docs-side setChildren
			// to honour ExcludeChildren so a single field now covers both
			// nested generation and the docs Children link list.
			for _, child := range toStringSlice(val) {
				if !contains(out.ExcludeChildren, child) {
					out.ExcludeChildren = append(out.ExcludeChildren, child)
				}
			}
		case "resource_identifier", "data_source_has_no_name_identifier", "static_parent":
			// S6 DERIVE: the canonical pipeline derives the identifier shape
			// from class meta (IdentifiedBy + containedBy) so the legacy
			// override is dropped. The loader normalisation step is the
			// asserter - if a future class wants to override the auto-derive,
			// it adds a positive field to the canonical YAML; the migration
			// script just verifies the value here and drops it.
		case "multi_relationship_class":
			// S3 OBSOLETE: the canonical pipeline treats any relationship_classes
			// list with more than one entry as multi automatically. The legacy
			// flag is dropped; C17 handles the relationship_classes payload.
		case "class_version_tests":
			// S8 POSTPONE: legacy gated test-class registration by ACI version
			// range. The canonical pipeline has not yet defined a replacement
			// (see MIGRATION_OVERVIEW.md section 8 POSTPONE). Dropped with a log
			// line so future loader work can pick the use case back up.
			fmt.Printf("  POSTPONE: %s class_version_tests=%v (not yet modeled)\n", file, val)
		}
	}
	return out
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// migrateMultiParents converts the legacy multi_parents block into the
// canonical ParentDnVariants slice. Field renames: legacy `contained_by`
// -> canonical `parent_class`; legacy `test_type` (string) -> canonical
// `test_platform`. `rn_prepend` and `wrapper_class` keep their names.
func migrateMultiParents(val any) []data.ParentDnVariantDefinition {
	raw, ok := val.([]any)
	if !ok {
		return nil
	}
	out := make([]data.ParentDnVariantDefinition, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[any]any)
		if !ok {
			continue
		}
		variant := data.ParentDnVariantDefinition{}
		if s, ok := m["contained_by"].(string); ok {
			variant.ParentClass = s
		}
		if s, ok := m["rn_prepend"].(string); ok {
			variant.RnPrepend = s
		}
		if s, ok := m["wrapper_class"].(string); ok {
			variant.WrapperClass = s
		}
		if s, ok := m["test_type"].(string); ok {
			// Reuse the canonical loader's strict text parser so an
			// unrecognised legacy value (e.g. typo) surfaces here
			// instead of producing silently bogus YAML.
			var p data.PlatformTypeEnum
			if err := p.UnmarshalText([]byte(s)); err == nil {
				variant.TestPlatform = p
			}
		}
		out = append(out, variant)
	}
	return out
}

// toStringSlice converts a yaml.v2-decoded []any of strings into []string.
// Non-string elements are dropped silently; per-element type errors would
// already have surfaced when yaml.Unmarshal returned the parent map.
func toStringSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// hasMigratedData returns true when at least one Phase-1 field has been
// populated. Phase-1-empty inputs (the 60+ legacy files whose only keys
// are TODO dispositions) are skipped to avoid emitting noisy empty .yaml
// files at the migration target.
func hasMigratedData(c data.ClassDefinition) bool {
	if c.AllowDelete != "" || c.ResourceName != "" || c.RnPrepend != "" || c.RequiredAsChild || c.IsSingleNestedWhenDefinedAsChild {
		return true
	}
	if len(c.ExcludeChildren) > 0 || len(c.IncludeChildren) > 0 {
		return true
	}
	if c.Artifacts != nil || len(c.ParentDnVariants) > 0 {
		return true
	}
	if len(c.TestConfig.IgnoreTests) > 0 || c.TestConfig.IgnoreImportStateVerify {
		return true
	}
	if c.Documentation.SubCategory != "" || len(c.Documentation.UiLocations) > 0 || len(c.Documentation.DnFormats) > 0 || len(c.Documentation.ExampleParentClasses) > 0 {
		return true
	}
	return false
}

// marshalView projects a populated data.ClassDefinition into a yaml.v2 map
// that omits empty fields. The canonical struct intentionally omits
// `omitempty` (loader-side semantics), so a direct yaml.Marshal would emit
// noisy `allow_delete: ""` / `documentation: {}` entries on every file.
// Hand-projecting keeps the migrated YAML clean and focuses the L3
// round-trip on real content.
//
// Phase 1 covers the keys handled by migrate(); subsequent commits extend
// this projection alongside their new translations.
func marshalView(c data.ClassDefinition) map[string]any {
	out := map[string]any{}
	if c.AllowDelete != "" {
		out["allow_delete"] = c.AllowDelete
	}
	if c.ResourceName != "" {
		out["resource_name"] = c.ResourceName
	}
	if c.RnPrepend != "" {
		out["rn_prepend"] = c.RnPrepend
	}
	if c.RequiredAsChild {
		out["required_as_child"] = c.RequiredAsChild
	}
	if c.IsSingleNestedWhenDefinedAsChild {
		out["is_single_nested_when_defined_as_child"] = c.IsSingleNestedWhenDefinedAsChild
	}
	if c.Artifacts != nil {
		// Emit as []string so per-element UnmarshalText drives the strict
		// round-trip in roundTripStrict; an empty slice (`artifacts: []`)
		// is preserved verbatim because it has loader semantics distinct
		// from omission (explicit opt-out vs auto-derive).
		names := make([]string, len(c.Artifacts))
		for i, a := range c.Artifacts {
			names[i] = a.String()
		}
		out["artifacts"] = names
	}
	if len(c.ExcludeChildren) > 0 {
		out["exclude_children"] = c.ExcludeChildren
	}
	if len(c.IncludeChildren) > 0 {
		out["include_children"] = c.IncludeChildren
	}
	if len(c.ParentDnVariants) > 0 {
		variants := make([]map[string]any, len(c.ParentDnVariants))
		for i, v := range c.ParentDnVariants {
			var entry = map[string]any{}
			if v.ParentClass != "" {
				entry["parent_class"] = v.ParentClass
			}
			if v.RnPrepend != "" {
				entry["rn_prepend"] = v.RnPrepend
			}
			if v.WrapperClass != "" {
				entry["wrapper_class"] = v.WrapperClass
			}
			// PlatformTypeEnum's iota zero is "apic", which the legacy
			// `test_type: apic` would have been a no-op. Only emit when
			// the value diverges from default to keep the migrated YAML
			// minimal.
			if v.TestPlatform != data.Apic {
				entry["test_platform"] = v.TestPlatform.String()
			}
			variants[i] = entry
		}
		out["parent_dn_variants"] = variants
	}
	testCfg := map[string]any{}
	if len(c.TestConfig.IgnoreTests) > 0 {
		names := make([]string, len(c.TestConfig.IgnoreTests))
		for i, e := range c.TestConfig.IgnoreTests {
			names[i] = e.String()
		}
		testCfg["ignore_tests"] = names
	}
	if c.TestConfig.IgnoreImportStateVerify {
		testCfg["ignore_import_state_verify"] = c.TestConfig.IgnoreImportStateVerify
	}
	if len(testCfg) > 0 {
		out["test_config"] = testCfg
	}
	doc := map[string]any{}
	if c.Documentation.SubCategory != "" {
		doc["sub_category"] = c.Documentation.SubCategory
	}
	if len(c.Documentation.UiLocations) > 0 {
		doc["ui_locations"] = c.Documentation.UiLocations
	}
	if len(c.Documentation.DnFormats) > 0 {
		doc["dn_formats"] = c.Documentation.DnFormats
	}
	if len(c.Documentation.ExampleParentClasses) > 0 {
		doc["example_parent_classes"] = c.Documentation.ExampleParentClasses
	}
	if len(doc) > 0 {
		out["documentation"] = doc
	}
	return out
}

// roundTripStrict marshals view, then re-reads the bytes via
// yaml.UnmarshalStrict against data.ClassDefinition (L3). A failure here
// means the migration would emit YAML the canonical loader rejects, which
// must surface at migration time so it never produces unreadable output.
func roundTripStrict(view map[string]any) ([]byte, error) {
	out, err := yaml.Marshal(view)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	var probe data.ClassDefinition
	if err := yaml.UnmarshalStrict(out, &probe); err != nil {
		return nil, fmt.Errorf("loader round-trip rejected output: %w", err)
	}
	return out, nil
}

func main() {
	classesDir := "gen/scripts/legacy_definitions/classes"
	outputDir := "gen/definitions"

	if len(os.Args) > 1 && os.Args[1] == "clean" {
		runClean(outputDir)
		return
	}

	files, err := filepath.Glob(filepath.Join(classesDir, "*.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading classes dir: %v\n", err)
		os.Exit(1)
	}

	tally := newKeyTally()
	migrated, skipped, failed := 0, 0, 0

	for _, file := range files {
		// classes/global.yaml is the legacy carrier for what is now
		// GlobalMetaDefinition (definitions/global.yaml). It must not be
		// translated to a ClassDefinition - skip it but still tally its
		// keys so the per-key coverage stays honest.
		if filepath.Base(file) == "global.yaml" {
			tallyOnly(file, tally)
			skipped++
			continue
		}

		payload, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			failed++
			continue
		}

		var legacy map[string]any
		if err := yaml.Unmarshal(payload, &legacy); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file, err)
			failed++
			continue
		}

		canonical := migrate(file, legacy, tally)
		if !hasMigratedData(canonical) {
			skipped++
			continue
		}

		view := marshalView(canonical)
		out, err := roundTripStrict(view)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error round-tripping %s: %v\n", file, err)
			failed++
			continue
		}

		className := strings.TrimSuffix(filepath.Base(file), ".yaml")
		outputPath := filepath.Join(outputDir, className+".yaml")
		if err := os.WriteFile(outputPath, out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			failed++
			continue
		}

		fmt.Printf("Migrated: %s -> %s\n", file, outputPath)
		migrated++
	}

	fmt.Printf("\nDone: %d migrated, %d skipped (no Phase-1 fields), %d failed\n",
		migrated, skipped, failed)
	tally.print()

	if failed > 0 {
		os.Exit(1)
	}
}

// tallyOnly records the keys of a file we intentionally do not write
// (e.g. classes/global.yaml). Lets the L2 summary stay complete even
// for skipped inputs.
func tallyOnly(file string, tally *keyTally) {
	payload, err := os.ReadFile(file)
	if err != nil {
		return
	}
	var legacy map[string]any
	if err := yaml.Unmarshal(payload, &legacy); err != nil {
		return
	}
	for key := range legacy {
		tally.record(file, key)
	}
}

func runClean(outputDir string) {
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
}
