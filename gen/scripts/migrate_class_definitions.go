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
	"children":             {sectionSemantic, true}, // -> include_children (subtract meta containedBy first)
	"contained_by":         {sectionSemantic, true}, // -> include_parents (subtract meta containedBy first)
	"class_version":        {sectionSemantic, true}, // -> supported_versions
	"relationship_classes": {sectionSemantic, true}, // -> relation_info.to_classes
	"migration_blocks":     {sectionSemantic, true}, // -> state_upgrades (two-source merge)
	"migration_version":    {sectionSemantic, true}, // -> state_upgrades (drives migration_source)
	"type_changes":         {sectionSemantic, true}, // -> state_upgrades.attributes
	"resource_notes":       {sectionSemantic, true}, // -> documentation.resource.notes

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

	// Property-level top-level keys (loaded from
	// gen/scripts/legacy_definitions/properties/*.yaml). Tallied under the
	// "prop:" prefix so they group separately from class-level keys with
	// matching names (e.g. "documentation" exists on both sides). Same
	// disposition codes apply; the per-key migration logic lives in
	// migrateProperties.

	// S1 direct
	"prop:documentation": {sectionDirect, true},

	// S2 semantic (this commit)
	"prop:overwrites":                {sectionSemantic, true},
	"prop:read_only_properties":      {sectionSemantic, true},
	"prop:resource_required":         {sectionSemantic, true},
	"prop:ignores":                   {sectionSemantic, true},
	"prop:type_overwrites":           {sectionSemantic, true},
	"prop:default_values":            {sectionSemantic, true},
	"prop:remove_valid_values":       {sectionSemantic, true},
	"prop:add_valid_values":          {sectionSemantic, true},
	"prop:ignore_properties_in_test": {sectionSemantic, true},

	// S2 semantic (later commits)
	"prop:test_values": {sectionSemantic, true}, // C20 (bucket merge)
	"prop:parents":     {sectionSemantic, true}, // this commit (C21 - decision tree + polymorphic)
	"prop:targets":     {sectionSemantic, true}, // this commit (C21 - decision tree + polymorphic)

	// S3 obsolete (later commit)
	"prop:exclude_targets":             {sectionObsolete, false}, // C22 (verify polymorphic-same-type)
	"prop:resource_name_doc_overwrite": {sectionReuse, false},    // C22 (drop - already in global)

	// S6 DERIVE (verify-and-drop; meta auto-derives ip_address from validateAsIPv4OrIPv6)
	"prop:static_custom_type": {sectionDerive, true},

	// S8 POSTPONE - drop with per-file warning, no slot in canonical struct yet
	"prop:ignore_custom_type_docs":     {sectionPostpone, true}, // section 8.3
	"prop:example_value_overwrite":     {sectionPostpone, true}, // section 8.4
	"prop:custom_test_dependency_name": {sectionPostpone, true}, // section 8.5
	"prop:datasource_required":         {sectionPostpone, true}, // section 8.2 (topSystem top-level only)
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

// recordProperty is the property-level companion to record. The "prop:" prefix
// keeps property and class top-level keys in distinct rows of the L2 tally
// even when they share a name (e.g. "documentation" or "overwrites").
func (t *keyTally) recordProperty(file, key string) keyInfo {
	return t.record(file, "prop:"+key)
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
		case "resource_notes":
			// -> documentation.resource.notes (resource-only, NOT the shared
			// documentation.notes). The legacy key always documented behaviour
			// specific to the resource side; sharing it onto the datasource
			// docs would render notes that don't apply.
			notes := toStringSlice(val)
			if len(notes) > 0 {
				out.Documentation.Resource.Notes = append(out.Documentation.Resource.Notes, notes...)
			}
		case "children":
			// -> include_children (union with anything already present).
			// The loader merges meta `containedBy`/`rnMap` itself and dedupes
			// via sortAndConvertToClassNames, so we do NOT subtract meta here.
			for _, child := range toStringSlice(val) {
				if !contains(out.IncludeChildren, child) {
					out.IncludeChildren = append(out.IncludeChildren, child)
				}
			}
		case "contained_by":
			// -> include_parents (union with anything already present).
			// Loader merges meta `containedBy` itself and dedupes.
			for _, parent := range toStringSlice(val) {
				if !contains(out.IncludeParents, parent) {
					out.IncludeParents = append(out.IncludeParents, parent)
				}
			}
		case "class_version":
			if s, ok := val.(string); ok {
				out.SupportedVersions = s
			}
		case "relationship_classes":
			// -> relation_info.to_classes. The list is taken verbatim; the
			// loader's NewClassName handles both bare ("fvAEPg") and meta-style
			// ("fv:AEPg") forms. multi_relationship_class is dropped separately
			// (C15) because the canonical pipeline derives the multi-class
			// shape from len(to_classes) > 1 automatically.
			out.RelationInfo.ToClasses = toStringSlice(val)
		case "migration_blocks", "migration_version", "type_changes":
			// Handled in migrateStateUpgrades after the switch so all three
			// keys can share the same versions map (the version index that
			// migration_version pins also has to be referenced by
			// migration_blocks attributes and type_changes entries).
		}
	}
	migrateStateUpgrades(file, legacy, &out)
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

// snakeToCamel converts "application_profile_dn" -> "applicationProfileDn".
// The legacy `migration_blocks` value is the new attribute_name (snake_case);
// the new state_upgrades schema keys by meta camelCase property name. Both
// the auto-derived attribute_name (when no override) and the manual
// attribute_name overrides round-trip to the meta camelCase via lowercase-first
// segment + Title-case the rest. Synthetic names that are already camelCase
// (parent_dn -> parentDn, tDn -> tDn, tnAaaDomainName -> tnAaaDomainName)
// round-trip through this rule cleanly.
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	var b strings.Builder
	b.WriteString(parts[0])
	for _, p := range parts[1:] {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

// parseLegacyAttributeType translates the legacy type_changes string
// ("SetNestedAttribute", "StringAttribute", ...) into the canonical
// LegacyAttributeTypeEnum value. The legacy value is the Go identifier
// from the prior SDKv2 framework; the canonical YAML uses the snake_case
// string ("set_nested_attribute"). Returns UndefinedLegacyAttributeType
// when the input does not match a known type so the caller can emit a
// warning rather than crashing on a typo.
func parseLegacyAttributeType(s string) data.LegacyAttributeTypeEnum {
	switch s {
	case "StringAttribute":
		return data.StringAttribute
	case "BoolAttribute":
		return data.BoolAttribute
	case "Int64Attribute":
		return data.Int64Attribute
	case "Float64Attribute":
		return data.Float64Attribute
	case "ListAttribute":
		return data.ListAttribute
	case "SetAttribute":
		return data.SetAttribute
	case "MapAttribute":
		return data.MapAttribute
	case "SingleNestedAttribute":
		return data.SingleNestedAttribute
	case "ListNestedAttribute":
		return data.ListNestedAttribute
	case "SetNestedAttribute":
		return data.SetNestedAttribute
	case "MapNestedAttribute":
		return data.MapNestedAttribute
	}
	return data.UndefinedLegacyAttributeType
}

// migrateStateUpgrades reads the three SDKv2-migration keys
// (migration_version, migration_blocks, type_changes) from the legacy YAML
// and assembles a single state_upgrades:[] tree on `out`. Also sets
// out.MigrationSource = FromSDKv2 when any of the three keys is present
// (validates the cross-field rule in class.go validateStateUpgrades).
//
// The two-source merge with legacy_definitions/schema-git-commit-e21fb3e5.json
// is NOT done here today - legacy_type and legacy_restriction default to
// the zero value (interpreted as "inherit from current property"), which
// covers the common case of a pure name rename. Authors needing a real
// type or restriction diff write the explicit override into the migrated
// YAML by hand; the script's only obligation is to capture the legacy
// name pairs faithfully.
//
// Per MIGRATION_OVERVIEW.md section 2.1, the script:
//   - Keys top-level migration_blocks entries (className == file) into
//     attributes[<camelCase(newName)>]
//   - Keys nested migration_blocks entries (className != file) into
//     children[<className>].attributes[<camelCase(newName-suffix)>]
//   - Splits dotted new names (rel.inner) so the inner suffix becomes
//     the attribute name and the outer prefix is dropped (the className
//     already carries the block identity)
//   - Maps type_changes[i] to either attributes[<attribute>] (scalar
//     legacy types) or children[<attribute>] (nested legacy types) on
//     the entry matching version, creating a fresh entry when no
//     migration_version landed at that version
func migrateStateUpgrades(file string, legacy map[string]any, out *data.ClassDefinition) {
	_, hasBlocks := legacy["migration_blocks"]
	versionVal, hasVersion := legacy["migration_version"]
	_, hasTypes := legacy["type_changes"]
	if !hasBlocks && !hasVersion && !hasTypes {
		return
	}
	selfClass := strings.TrimSuffix(filepath.Base(file), ".yaml")

	// versions indexes state_upgrades by PriorSchemaVersion so all three
	// keys can find/create the same entry.
	versions := map[int]*data.StateUpgradeDefinition{}
	ensure := func(v int) *data.StateUpgradeDefinition {
		if entry, ok := versions[v]; ok {
			return entry
		}
		entry := &data.StateUpgradeDefinition{
			PriorSchemaVersion: v,
			Attributes:         map[string]data.AttributeUpgradeDefinition{},
			Children:           map[string]data.AttributeUpgradeDefinition{},
		}
		versions[v] = entry
		return entry
	}

	primaryVersion := 0
	if hasVersion {
		switch v := versionVal.(type) {
		case int:
			primaryVersion = v
		case int64:
			primaryVersion = int(v)
		}
	}
	if hasBlocks || hasVersion {
		ensure(primaryVersion)
	}

	if hasBlocks {
		blocks, _ := legacy["migration_blocks"].(map[any]any)
		for classKey, attrMap := range blocks {
			className, _ := classKey.(string)
			attrs, _ := attrMap.(map[any]any)
			entry := ensure(primaryVersion)
			isSelf := className == selfClass
			for oldKey, newVal := range attrs {
				oldName, _ := oldKey.(string)
				newName, _ := newVal.(string)
				// Dotted form: outer prefix is the relation/block attribute
				// (already implicit in className), suffix is the inner attr.
				inner := newName
				if idx := strings.Index(newName, "."); idx >= 0 {
					inner = newName[idx+1:]
				}
				camelKey := snakeToCamel(inner)
				upgrade := data.AttributeUpgradeDefinition{
					LegacyAttribute: oldName,
				}
				if isSelf {
					entry.Attributes[camelKey] = upgrade
				} else {
					child := entry.Children[className]
					if child.Attributes == nil {
						child.Attributes = map[string]data.AttributeUpgradeDefinition{}
					}
					child.Attributes[camelKey] = upgrade
					entry.Children[className] = child
				}
			}
		}
	}

	if hasTypes {
		typesList, _ := legacy["type_changes"].([]any)
		for _, item := range typesList {
			entryMap, _ := item.(map[any]any)
			if entryMap == nil {
				continue
			}
			attrName, _ := entryMap["attribute"].(string)
			oldType, _ := entryMap["old_type"].(string)
			versionInt := 0
			switch v := entryMap["version"].(type) {
			case int:
				versionInt = v
			case int64:
				versionInt = int(v)
			}
			parsed := parseLegacyAttributeType(oldType)
			if parsed == data.UndefinedLegacyAttributeType {
				fmt.Printf("  WARN: %s type_changes attribute=%s old_type=%q is not a recognised legacy type; entry kept with zero legacy_type\n", file, attrName, oldType)
			}
			entry := ensure(versionInt)
			// Nested legacy types route to .children; scalar legacy types
			// route to .attributes. The split mirrors the state_upgrades
			// loader: children carries block-shaped entries, attributes
			// carries scalar entries.
			isNested := strings.Contains(oldType, "Nested")
			upgrade := data.AttributeUpgradeDefinition{
				LegacyType: parsed,
			}
			if isNested {
				entry.Children[attrName] = upgrade
			} else {
				entry.Attributes[attrName] = upgrade
			}
		}
	}

	// Order entries by ascending PriorSchemaVersion so the emitted YAML
	// is reproducible across runs (Go map iteration is randomised).
	keys := make([]int, 0, len(versions))
	for k := range versions {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		out.StateUpgrades = append(out.StateUpgrades, *versions[k])
	}

	// MigrationSource is the cross-field that drives the docs migration
	// warning and the validator that requires at least one state_upgrades
	// entry. The legacy YAML implies SDKv2 lineage whenever any of the
	// three keys is present, so the migration script sets the enum on
	// every translated file.
	out.MigrationSource = data.FromSDKv2
}

// projectStateUpgrade emits a StateUpgradeDefinition as a YAML-friendly map
// with enums hand-encoded via String() (yaml.v2 has no MarshalText support)
// and zero-value fields omitted. Mirrors the canonical loader's strict
// expectation: the YAML keys are meta camelCase, and unset fields surface
// as the loader zero value (interpreted as "inherit from current property").
func projectStateUpgrade(su data.StateUpgradeDefinition) map[string]any {
	entry := map[string]any{
		"prior_schema_version": su.PriorSchemaVersion,
	}
	if len(su.Attributes) > 0 {
		entry["attributes"] = projectAttributeMap(su.Attributes)
	}
	if len(su.Children) > 0 {
		entry["children"] = projectAttributeMap(su.Children)
	}
	return entry
}

// projectAttributeMap walks AttributeUpgradeDefinition values recursively
// (Attributes/Children at any depth) and hand-projects enums to strings.
func projectAttributeMap(m map[string]data.AttributeUpgradeDefinition) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		out[k] = projectAttributeUpgrade(v)
	}
	return out
}

func projectAttributeUpgrade(a data.AttributeUpgradeDefinition) map[string]any {
	out := map[string]any{}
	if a.LegacyAttribute != "" {
		out["legacy_attribute"] = a.LegacyAttribute
	}
	if a.LegacyType != data.UndefinedLegacyAttributeType {
		out["legacy_type"] = a.LegacyType.String()
	}
	if a.LegacyRestriction != data.UndefinedRestriction {
		out["legacy_restriction"] = a.LegacyRestriction.String()
	}
	// LegacyStatus's zero value (Functioning) is the documented default;
	// emit only when non-zero so the migrated YAML stays minimal.
	if a.LegacyStatus != data.Functioning {
		out["legacy_status"] = a.LegacyStatus.String()
	}
	if len(a.Attributes) > 0 {
		out["attributes"] = projectAttributeMap(a.Attributes)
	}
	if len(a.Children) > 0 {
		out["children"] = projectAttributeMap(a.Children)
	}
	return out
}

// parseValueType maps a legacy `type_overwrites` value (single token) to the
// canonical ValueTypeEnum. Today every entry across the 3 affected files
// uses "string" - the migration logs and bypasses any other value rather
// than crashing the run.
func parseValueType(file, propName, s string) data.ValueTypeEnum {
	switch s {
	case "string":
		return data.String
	case "set":
		return data.Set
	case "object":
		return data.Object
	case "ip_address":
		return data.IpAddress
	case "semantic_equality":
		return data.SemanticEquality
	default:
		fmt.Printf("WARN: %s: unknown type_overwrites value %q for %s - left as UndefinedValueType\n", file, s, propName)
		return data.UndefinedValueType
	}
}

// upsertProperty fetches (or creates) the PropertyDefinition for the given
// meta property name. Centralising the lookup keeps the per-key handlers in
// migrateProperties terse and avoids subtle map-value-vs-pointer bugs.
func upsertProperty(out *data.ClassDefinition, metaName string) data.PropertyDefinition {
	if out.Properties == nil {
		out.Properties = map[string]data.PropertyDefinition{}
	}
	return out.Properties[metaName]
}

// setRestriction writes the restriction onto the property; warns when the
// caller would overwrite an already-set restriction with a different value
// (e.g. a property listed in both `resource_required` and `ignores`, which
// is a legacy-file authoring error worth surfacing).
func setRestriction(file, metaName string, prop *data.PropertyDefinition, r data.RestrictionEnum, source string) {
	if prop.Restriction != data.UndefinedRestriction && prop.Restriction != r {
		fmt.Printf("WARN: %s: restriction conflict on %s: already %s, %s wants %s\n",
			file, metaName, prop.Restriction, source, r)
		return
	}
	prop.Restriction = r
}

// asStringSlice converts a yaml.v2 []any (or []string) to a []string and
// silently drops non-string entries with a warning. The shape `[]any` is
// what yaml.v2 produces for inline lists.
func asStringSlice(v any) []string {
	switch xs := v.(type) {
	case []string:
		return xs
	case []any:
		out := make([]string, 0, len(xs))
		for _, e := range xs {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// migrateProperties translates one legacy property YAML payload (parsed as
// map[string]any) into per-property PropertyDefinition entries on the
// outgoing ClassDefinition. The translation is narrow: only keys flagged
// `implemented:true` in the "prop:" entries of knownLegacyKeys are written
// into the canonical struct. Every other known key is tallied so the L2
// coverage gap stays visible; unknown keys raise an L1 warning via the
// tally just like the class loader does.
func migrateProperties(file string, legacy map[string]any, tally *keyTally, out *data.ClassDefinition) {
	// The test_values bucket sub-keys are snake_case new attribute names
	// (post-overwrites), and the canonical struct keys PropertyDefinition
	// entries by meta camelCase. Build the snake_new -> meta_camel inverter
	// up front so every bucket entry can resolve its target property in O(1)
	// regardless of map iteration order. Properties without an `overwrites`
	// entry fall through to snakeToCamel, which auto-derives the meta name.
	invertOverwrites := buildOverwritesInverter(legacy)

	for key, val := range legacy {
		info := tally.recordProperty(file, key)
		if !info.implemented {
			continue
		}
		switch key {
		case "documentation":
			// Map<metaName, description-template-string> -> per-property
			// PropertyDefinition.Documentation.Description.
			docs, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: documentation is not a map: %T\n", file, val)
				continue
			}
			for metaKey, descVal := range docs {
				metaName, _ := metaKey.(string)
				desc, _ := descVal.(string)
				if metaName == "" {
					continue
				}
				prop := upsertProperty(out, metaName)
				prop.Documentation.Description = desc
				out.Properties[metaName] = prop
			}

		case "overwrites":
			// Map<snake_case_old_attr, snake_case_new_attr>. The new
			// schema keys per-property overrides by meta camelCase name,
			// so transform the legacy snake_case attribute name to the
			// camelCase meta name via snakeToCamel (same rule used for
			// state_upgrades). The value (new attribute name) lands in
			// PropertyDefinition.AttributeName verbatim.
			ows, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: overwrites is not a map: %T\n", file, val)
				continue
			}
			for snakeKey, newNameVal := range ows {
				oldSnake, _ := snakeKey.(string)
				newName, _ := newNameVal.(string)
				if oldSnake == "" || newName == "" {
					continue
				}
				metaName := snakeToCamel(oldSnake)
				prop := upsertProperty(out, metaName)
				prop.AttributeName = newName
				out.Properties[metaName] = prop
			}

		case "read_only_properties":
			for _, metaName := range asStringSlice(val) {
				prop := upsertProperty(out, metaName)
				setRestriction(file, metaName, &prop, data.ReadOnly, "read_only_properties")
				out.Properties[metaName] = prop
			}

		case "resource_required":
			for _, metaName := range asStringSlice(val) {
				prop := upsertProperty(out, metaName)
				setRestriction(file, metaName, &prop, data.Required, "resource_required")
				out.Properties[metaName] = prop
			}

		case "ignores":
			// Legacy `ignores` is class-scoped (cross-class lives in
			// properties/global.yaml as "ignores" and migrates to
			// global.exclude_properties). Per-class entries map to
			// PropertyDefinition.Restriction: exclude.
			for _, metaName := range asStringSlice(val) {
				prop := upsertProperty(out, metaName)
				setRestriction(file, metaName, &prop, data.Exclude, "ignores")
				out.Properties[metaName] = prop
			}

		case "type_overwrites":
			tos, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: type_overwrites is not a map: %T\n", file, val)
				continue
			}
			for metaKey, typeVal := range tos {
				metaName, _ := metaKey.(string)
				typeStr, _ := typeVal.(string)
				if metaName == "" {
					continue
				}
				vt := parseValueType(file, metaName, typeStr)
				if vt == data.UndefinedValueType {
					continue
				}
				prop := upsertProperty(out, metaName)
				prop.ValueType = vt
				out.Properties[metaName] = prop
			}

		case "default_values":
			// Legacy: `metaPropName: value` (single value).
			// Canonical: `metaPropName: { value: versionRange }`. Empty
			// version string means "applies to all versions" - migrate
			// every entry with empty version range so the meta-derived
			// version-scoped semantics carry over unchanged.
			dvs, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: default_values is not a map: %T\n", file, val)
				continue
			}
			for metaKey, defVal := range dvs {
				metaName, _ := metaKey.(string)
				if metaName == "" {
					continue
				}
				defStr := fmt.Sprintf("%v", defVal)
				prop := upsertProperty(out, metaName)
				if prop.DefaultValues == nil {
					prop.DefaultValues = map[string]string{}
				}
				prop.DefaultValues[defStr] = ""
				out.Properties[metaName] = prop
			}

		case "remove_valid_values":
			rvs, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: remove_valid_values is not a map: %T\n", file, val)
				continue
			}
			for metaKey, listVal := range rvs {
				metaName, _ := metaKey.(string)
				if metaName == "" {
					continue
				}
				prop := upsertProperty(out, metaName)
				prop.RemoveValidValues = append(prop.RemoveValidValues, asStringSlice(listVal)...)
				out.Properties[metaName] = prop
			}

		case "add_valid_values":
			avs, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: add_valid_values is not a map: %T\n", file, val)
				continue
			}
			for metaKey, listVal := range avs {
				metaName, _ := metaKey.(string)
				if metaName == "" {
					continue
				}
				prop := upsertProperty(out, metaName)
				prop.AddValidValues = append(prop.AddValidValues, asStringSlice(listVal)...)
				out.Properties[metaName] = prop
			}

		case "ignore_properties_in_test":
			// Legacy shape is map<metaName, "no"> (the value is the
			// historical "yes/no" string; "no" is the only observed
			// value). New shape is a boolean test_config.ignore_in_test
			// per property; flip true regardless of legacy value because
			// the key's presence is what disables the test entry.
			ips, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: ignore_properties_in_test is not a map: %T\n", file, val)
				continue
			}
			for metaKey := range ips {
				metaName, _ := metaKey.(string)
				if metaName == "" {
					continue
				}
				prop := upsertProperty(out, metaName)
				prop.TestConfig.IgnoreInTest = true
				out.Properties[metaName] = prop
			}

		case "static_custom_type":
			// Legacy shape: map<metaName, custom-type-token>.
			//   "ip_address"        -> auto-derived from meta validateAsIPv4OrIPv6
			//                          (drop silently; covers 14 of 15 entries).
			//   "vmm_arp_learning"  -> POSTPONE: log a warning, drop from output.
			//                          Hand-restore once section 8.1 lands.
			sct, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: static_custom_type is not a map: %T\n", file, val)
				continue
			}
			for metaKey, typeVal := range sct {
				metaName, _ := metaKey.(string)
				typeStr, _ := typeVal.(string)
				if typeStr == "ip_address" {
					continue
				}
				fmt.Printf("POSTPONE: %s: static_custom_type[%s]=%s dropped (no slot in canonical struct yet; restore after section 8 lands)\n",
					file, metaName, typeStr)
			}

		case "ignore_custom_type_docs",
			"example_value_overwrite",
			"custom_test_dependency_name",
			"datasource_required":
			// S8 POSTPONE: no slot in canonical struct yet; drop with a
			// per-file warning so the loss is visible during migration.
			fmt.Printf("POSTPONE: %s: %s dropped (no canonical slot; restore after section 8 lands)\n", file, key)

		case "test_values":
			migrateTestValues(file, val, invertOverwrites, out)

		case "parents":
			migrateParents(file, val, out)

		case "targets":
			migrateTargets(file, val, out)
		}
	}
}

// buildOverwritesInverter returns a map<snake_case_new_attr_name, meta_camel_case>
// derived from the legacy `overwrites` block. Used by migrateTestValues to
// resolve a bucket entry's snake_case attribute key back to the canonical
// PropertyDefinition meta key. Returns an empty map when overwrites is
// absent or malformed; callers fall back to snakeToCamel for unmapped keys.
func buildOverwritesInverter(legacy map[string]any) map[string]string {
	out := map[string]string{}
	raw, ok := legacy["overwrites"]
	if !ok {
		return out
	}
	ows, ok := raw.(map[any]any)
	if !ok {
		return out
	}
	for snakeKey, newNameVal := range ows {
		oldSnake, _ := snakeKey.(string)
		newName, _ := newNameVal.(string)
		if oldSnake == "" || newName == "" {
			continue
		}
		out[newName] = snakeToCamel(oldSnake)
	}
	return out
}

// legacyValueToHCL renders a yaml.v2-decoded test value as a single string
// suitable for ConfigValue. Scalars stringify via fmt.Sprintf; lists render
// as bracketed comma-separated literals so the renderer can copy them into
// HCL verbatim. The single observed v2.19.0 list usage is commSsh's
// kex_algorithms: [curve25519-sha256], which becomes ["curve25519-sha256"].
// TODO once the test renderer lands: refine list element quoting to honour
// the property's ValueType (e.g. drop quotes around numeric/boolean tokens).
func legacyValueToHCL(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case []any:
		parts := make([]string, 0, len(x))
		for _, e := range x {
			parts = append(parts, fmt.Sprintf("%q", legacyValueToHCL(e)))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	default:
		return fmt.Sprintf("%v", x)
	}
}

// migrateTestValues translates one legacy `test_values:` block into per
// property TestConfigDefinition entries on the outgoing ClassDefinition.
//
// Bucket mapping (per MIGRATION_OVERVIEW.md section 2 test_values row):
//
//	all      -> PropertyDefinition.TestConfig.Create
//	default  -> PropertyDefinition.TestConfig.Default
//	legacy   -> PropertyDefinition.TestConfig.Legacy
//	update   -> PropertyDefinition.TestConfig.Update   (no v2.19.0 data)
//	force_new -> PropertyDefinition.TestConfig.ForceNew (no v2.19.0 data)
//
// Nested keys (per the same row):
//
//	custom_type            -> obsolete (section 3); silently dropped.
//	                          The 10 IpAddress files fold the IPv6 author
//	                          choice into the standard buckets; the 1
//	                          vmm_arp_learning ride-along is section 8.1.
//	datasource_required    -> DERIVE (section 6); IdentifiedBy + Default
//	                          bucket drive the datasource lookup config.
//	datasource_non_existing -> DERIVE (section 6); type-aware non-matching
//	                          transform of resource_required.
//	child_*, ignore_in_*   -> class-level test_config.children / dependency
//	                          overrides; v2.19.0 carries no data here.
//
// Bucket entries are keyed by snake_case new attribute name; resolveMetaName
// consults invertOverwrites first, then falls back to snakeToCamel. The
// resolved meta camelCase name becomes the PropertyDefinition map key.
func migrateTestValues(file string, val any, invertOverwrites map[string]string, out *data.ClassDefinition) {
	tv, ok := val.(map[any]any)
	if !ok {
		fmt.Printf("WARN: %s: test_values is not a map: %T\n", file, val)
		return
	}
	// Sort bucket names so the per-file run order is deterministic and the
	// per-property Create/Default/Legacy slices land in the same order on
	// every invocation.
	bucketNames := make([]string, 0, len(tv))
	for k := range tv {
		if s, ok := k.(string); ok {
			bucketNames = append(bucketNames, s)
		}
	}
	sort.Strings(bucketNames)

	for _, bucket := range bucketNames {
		entry := tv[bucket]
		switch bucket {
		case "all", "default", "legacy", "update", "force_new":
			entries, ok := entry.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: test_values.%s is not a map: %T\n", file, bucket, entry)
				continue
			}
			// Sort the per-bucket property keys for stable output.
			keys := make([]string, 0, len(entries))
			for k := range entries {
				if s, ok := k.(string); ok {
					keys = append(keys, s)
				}
			}
			sort.Strings(keys)
			for _, snakeKey := range keys {
				metaName := resolveMetaName(snakeKey, invertOverwrites)
				if metaName == "" {
					continue
				}
				tve := data.TestValueEntryDefinition{
					ConfigValue: legacyValueToHCL(entries[snakeKey]),
				}
				prop := upsertProperty(out, metaName)
				switch bucket {
				case "all":
					prop.TestConfig.Create = append(prop.TestConfig.Create, tve)
				case "default":
					prop.TestConfig.Default = append(prop.TestConfig.Default, tve)
				case "legacy":
					prop.TestConfig.Legacy = append(prop.TestConfig.Legacy, tve)
				case "update":
					prop.TestConfig.Update = append(prop.TestConfig.Update, tve)
				case "force_new":
					prop.TestConfig.ForceNew = append(prop.TestConfig.ForceNew, tve)
				}
				out.Properties[metaName] = prop
			}

		case "custom_type":
			// Obsolete per section 3: the legacy `custom_type` bucket was
			// only used to author IPv6 examples for the 10 IpAddress
			// files; the standard create/default buckets cover the same
			// data after migration. The vmm_arp_learning ride-along is
			// section 8.1.
			continue

		case "datasource_required", "datasource_non_existing":
			// DERIVE per section 6: IdentifiedBy + Default bucket drive
			// the datasource lookup config and the non-existing branch;
			// no migration of these nested keys today.
			continue

		case "resource_required", "test_values_for_parent":
			// DERIVE per section 6 (resource_required is the variant used
			// when the property is required-vs-optional; test_values_for_
			// parent is the values used when this class is rendered as a
			// parent in another class's test). Both are reproducible from
			// the standard Create/Default buckets plus the property's
			// Restriction; drop them silently and let the loader/renderer
			// rebuild them.
			continue

		default:
			// child_*, ignore_in_*, and any future nested key. Today
			// v2.19.0 carries no data here; log so the gap is visible if
			// the corpus ever ships an entry.
			if strings.HasPrefix(bucket, "child_") || strings.HasPrefix(bucket, "ignore_in_") {
				continue
			}
			fmt.Printf("WARN: %s: test_values.%s ignored (unknown bucket)\n", file, bucket)
		}
	}
}

// resolveMetaName returns the meta camelCase property name for a bucket
// entry's snake_case attribute name. invertOverwrites takes precedence;
// any unmapped key falls back to snakeToCamel for auto-derived properties.
func resolveMetaName(snakeKey string, invertOverwrites map[string]string) string {
	if v, ok := invertOverwrites[snakeKey]; ok {
		return v
	}
	return snakeToCamel(snakeKey)
}

// classifyReference picks a ReferenceTypeEnum for a parent_dn / target_dn
// value: anything that starts with `aci_` or `data.aci_` is a Terraform
// reference expression; everything else (uni/..., topology/..., empty) is
// a static DN literal. The single observed `data.aci_*` usage today is in
// migrated scenarios that pre-date the renderer; v2.19.0 only emits the
// `aci_*` resource form.
func classifyReference(ref string) data.ReferenceTypeEnum {
	switch {
	case strings.HasPrefix(ref, "data.aci_"):
		return data.DataSourceReference
	case strings.HasPrefix(ref, "aci_"):
		return data.ResourceReference
	default:
		return data.StaticReference
	}
}

// stringifyConfigOverrides converts a legacy `properties:` map (yaml.v2's
// map[any]any of scalar values) to the canonical ConfigOverrides shape
// map[string]string. Lists collapse to bracketed literals via
// legacyValueToHCL so the renderer can copy them into HCL verbatim.
func stringifyConfigOverrides(file, who string, raw any) map[string]string {
	m, ok := raw.(map[any]any)
	if !ok {
		fmt.Printf("WARN: %s: %s properties is not a map: %T\n", file, who, raw)
		return nil
	}
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		key, _ := k.(string)
		if key == "" {
			continue
		}
		out[key] = legacyValueToHCL(v)
	}
	return out
}

// migrateParents translates one legacy `parents:` block into class-level
// TestDependencyDefinition entries with Role=Parent, plus polymorphic-
// same-type detection (section 8.6).
//
// Sub-key mapping (per MIGRATION_OVERVIEW.md section 2.2 parents row):
//
//	class_name             -> ClassName.
//	parent_dn              -> Reference + ReferenceType (aci_*/data.aci_*
//	                          vs static DN).
//	properties             -> ConfigOverrides.
//	parent_dependency      -> recursive single-level Dependencies[] entry.
//	parent_dependency_name -> the dep's Reference (rare, 1 file today).
//	target_classes         -> READ for polymorphic-same-type detection;
//	                          NOT emitted (the filter now lives on
//	                          Relation.ToClasses + the polymorphic
//	                          auto-detector).
//	class_in_parent        -> dropped (covered by the recursive dependency
//	                          shape; warn if data ever differs).
//
// Auto-resolution filtering (dropping entries already produced by
// meta containedBy + resolveParentDependencies) is intentionally NOT
// performed today — the migration script has no meta JSON loader yet.
// All legacy entries are emitted verbatim with `replace_auto_resolved`
// left at its default (false), trusting the loader's additive merge to
// deduplicate. A follow-on commit can prune redundant entries once meta
// loading lands; see MIGRATION_OVERVIEW section 2.2.
func migrateParents(file string, val any, out *data.ClassDefinition) {
	if val == nil {
		return
	}
	list, ok := val.([]any)
	if !ok {
		fmt.Printf("WARN: %s: parents is not a list: %T\n", file, val)
		return
	}
	// Track <entry.class_name, entry.target_classes> for polymorphic-same-
	// type detection at the end. Only fvRsSecInherited matches today.
	type parentSample struct {
		ClassName     string
		TargetClasses []string
	}
	samples := make([]parentSample, 0, len(list))

	for _, raw := range list {
		entry, ok := raw.(map[any]any)
		if !ok {
			fmt.Printf("WARN: %s: parents entry is not a map: %T\n", file, raw)
			continue
		}
		className, _ := entry["class_name"].(string)
		if className == "" {
			// Legacy quirk: 4 v2.19.0 files specify only parent_dependency
			// in their parents entry (no class_name), implicitly relying on
			// the legacy generator's "first parent from containedBy" rule.
			// The new pipeline's resolveParentDependencies covers the same
			// behaviour from meta containedBy, so drop the entry silently;
			// log just to keep the corpus check visible.
			fmt.Printf("DROP: %s: parents entry without class_name (auto-resolution covers via meta containedBy)\n", file)
			continue
		}
		dep := data.TestDependencyDefinition{
			ClassName: className,
			Role:      data.Parent,
		}
		if pd, ok := entry["parent_dn"].(string); ok && pd != "" {
			dep.Reference = pd
			dep.ReferenceType = classifyReference(pd)
		}
		if props, ok := entry["properties"]; ok {
			dep.ConfigOverrides = stringifyConfigOverrides(file, "parents."+className, props)
		}
		if pdep, ok := entry["parent_dependency"].(string); ok && pdep != "" {
			inner := data.TestDependencyDefinition{ClassName: pdep}
			if pdepName, ok := entry["parent_dependency_name"].(string); ok && pdepName != "" {
				inner.Reference = pdepName
				inner.ReferenceType = classifyReference(pdepName)
			}
			dep.Dependencies = []data.TestDependencyDefinition{inner}
		} else if pdepName, ok := entry["parent_dependency_name"].(string); ok && pdepName != "" {
			fmt.Printf("WARN: %s: parents.%s has parent_dependency_name without parent_dependency; dropped\n", file, className)
		}
		if cip, ok := entry["class_in_parent"]; ok {
			fmt.Printf("WARN: %s: parents.%s.class_in_parent=%v dropped (covered by recursive dependency shape)\n", file, className, cip)
		}
		out.TestConfig.Dependencies = append(out.TestConfig.Dependencies, dep)

		// Collect target_classes for polymorphic detection. The slot is
		// always a []string when present; missing slot is recorded as nil.
		sample := parentSample{ClassName: className}
		if tcRaw, ok := entry["target_classes"]; ok {
			sample.TargetClasses = asStringSlice(tcRaw)
		}
		samples = append(samples, sample)
	}

	// Polymorphic-same-type detection (section 8.6): when every parents
	// entry's target_classes list is exactly [class_name] (1:1 match), the
	// resolved Relation.ToClasses must equal the distinct union of those
	// class_names so the polymorphic auto-detector at render time can
	// drive multi-scenario rendering without any further YAML hint. Only
	// fvRsSecInherited triggers today, expanding ToClasses from
	// [fvAEPg, fvESg] to [fvAEPg, fvESg, l3extInstP].
	if len(samples) < 2 {
		return
	}
	allOneToOne := true
	for _, s := range samples {
		if len(s.TargetClasses) != 1 || s.TargetClasses[0] != s.ClassName {
			allOneToOne = false
			break
		}
	}
	if !allOneToOne {
		return
	}
	seen := map[string]bool{}
	for _, c := range out.RelationInfo.ToClasses {
		seen[c] = true
	}
	for _, s := range samples {
		if !seen[s.ClassName] {
			out.RelationInfo.ToClasses = append(out.RelationInfo.ToClasses, s.ClassName)
			seen[s.ClassName] = true
		}
	}
}

// migrateTargets translates one legacy `targets:` block into class-level
// TestDependencyDefinition entries with Role=Target.
//
// Sub-key mapping (per MIGRATION_OVERVIEW.md section 2.2 targets row):
//
//	class_name              -> ClassName.
//	target_dn / target_dn_ref -> Reference + ReferenceType. target_dn_ref
//	                             wins when present (the 1 v2.19.0 file
//	                             using it is the LookupTestValue dead-code
//	                             entry, but the ref form is the canonical
//	                             one). `static: true` forces StaticReference.
//	properties              -> ConfigOverrides.
//	parent_dependency       -> recursive single-level Dependencies[] entry.
//	relation_resource_name  -> dropped (canonical struct derives the HCL
//	                             resource name from class_name).
//	overwrite_parent_dn_key -> dropped (canonical schema infers the
//	                             parent_dn attribute name on the target).
//	target_dn_overwrite_docs -> dropped (S3 obsolete, doc-only).
//	shared_classes          -> dropped (no canonical slot; 3 files).
//	parent_dependency_dn_ref -> dropped (canonical struct derives; 2 files).
//
// Same auto-resolution caveat as migrateParents: emits all legacy entries
// verbatim; pruning the entries that single-target auto-resolution would
// reproduce is a follow-on commit gated on meta JSON loading.
func migrateTargets(file string, val any, out *data.ClassDefinition) {
	if val == nil {
		return
	}
	list, ok := val.([]any)
	if !ok {
		fmt.Printf("WARN: %s: targets is not a list: %T\n", file, val)
		return
	}
	for _, raw := range list {
		entry, ok := raw.(map[any]any)
		if !ok {
			fmt.Printf("WARN: %s: targets entry is not a map: %T\n", file, raw)
			continue
		}
		className, _ := entry["class_name"].(string)
		if className == "" {
			fmt.Printf("WARN: %s: targets entry missing class_name\n", file)
			continue
		}
		dep := data.TestDependencyDefinition{
			ClassName: className,
			Role:      data.Target,
		}
		// target_dn_ref wins when present; otherwise fall back to
		// target_dn. The `static: true` flag forces StaticReference even
		// for values that look like resource references.
		if ref, ok := entry["target_dn_ref"].(string); ok && ref != "" {
			dep.Reference = ref
			dep.ReferenceType = classifyReference(ref)
		} else if td, ok := entry["target_dn"].(string); ok && td != "" {
			dep.Reference = td
			dep.ReferenceType = classifyReference(td)
		}
		if isStatic, ok := entry["static"].(bool); ok && isStatic {
			dep.ReferenceType = data.StaticReference
		}
		if props, ok := entry["properties"]; ok {
			dep.ConfigOverrides = stringifyConfigOverrides(file, "targets."+className, props)
		}
		if pdep, ok := entry["parent_dependency"].(string); ok && pdep != "" {
			dep.Dependencies = []data.TestDependencyDefinition{{ClassName: pdep}}
		}
		// The remaining sub-keys (relation_resource_name,
		// overwrite_parent_dn_key, target_dn_overwrite_docs,
		// shared_classes, parent_dependency_dn_ref) intentionally have no
		// canonical slot — see the function doc for the rationale.
		out.TestConfig.Dependencies = append(out.TestConfig.Dependencies, dep)
	}
}

// projectTestDependency emits a TestDependencyDefinition as a YAML-friendly
// map, omitting zero-value fields and recursing through Dependencies.
// Enums are hand-emitted via String() (yaml.v2-no-MarshalText pattern).
func projectTestDependency(d data.TestDependencyDefinition) map[string]any {
	m := map[string]any{}
	if d.ClassName != "" {
		m["class_name"] = d.ClassName
	}
	if d.Reference != "" {
		m["reference"] = d.Reference
	}
	if d.ReferenceType != data.ResourceReference {
		m["reference_type"] = d.ReferenceType.String()
	}
	if d.Role != data.UndefinedRole {
		m["role"] = d.Role.String()
	}
	if len(d.ConfigOverrides) > 0 {
		m["config_overrides"] = d.ConfigOverrides
	}
	if len(d.Dependencies) > 0 {
		nested := make([]map[string]any, len(d.Dependencies))
		for i, dd := range d.Dependencies {
			nested[i] = projectTestDependency(dd)
		}
		m["dependencies"] = nested
	}
	return m
}

func projectTestDependencies(deps []data.TestDependencyDefinition) []map[string]any {
	out := make([]map[string]any, len(deps))
	for i, d := range deps {
		out[i] = projectTestDependency(d)
	}
	return out
}

// projectProperty emits a PropertyDefinition as a YAML-friendly map. Hand
// emits enums via String() (same yaml.v2-no-MarshalText pattern as the
// Artifacts and StateUpgrades projections) and omits zero-value fields.
func projectProperty(p data.PropertyDefinition) map[string]any {
	out := map[string]any{}
	if p.AttributeName != "" {
		out["attribute_name"] = p.AttributeName
	}
	if p.Restriction != data.UndefinedRestriction {
		out["restriction"] = p.Restriction.String()
	}
	if p.ValueType != data.UndefinedValueType {
		out["value_type"] = p.ValueType.String()
	}
	if len(p.DefaultValues) > 0 {
		out["default_values"] = p.DefaultValues
	}
	if len(p.AddValidValues) > 0 {
		out["add_valid_values"] = p.AddValidValues
	}
	if len(p.RemoveValidValues) > 0 {
		out["remove_valid_values"] = p.RemoveValidValues
	}
	if p.Documentation.Description != "" {
		out["documentation"] = map[string]any{
			"description": p.Documentation.Description,
		}
	}
	tc := projectTestConfig(p.TestConfig)
	if len(tc) > 0 {
		out["test_config"] = tc
	}
	return out
}

// projectTestConfig builds the YAML test_config block for a property,
// omitting empty buckets. Each TestValueEntryDefinition becomes a single
// `config_value: <string>` map (other entry fields are zero today; the
// renderer fleshes them out when the higher-fidelity buckets are
// implemented).
func projectTestConfig(tc data.TestConfigDefinition) map[string]any {
	out := map[string]any{}
	if len(tc.Create) > 0 {
		out["create"] = projectTestEntries(tc.Create)
	}
	if len(tc.Default) > 0 {
		out["default"] = projectTestEntries(tc.Default)
	}
	if len(tc.Update) > 0 {
		out["update"] = projectTestEntries(tc.Update)
	}
	if len(tc.ForceNew) > 0 {
		out["force_new"] = projectTestEntries(tc.ForceNew)
	}
	if len(tc.Legacy) > 0 {
		out["legacy"] = projectTestEntries(tc.Legacy)
	}
	if tc.IgnoreInTest {
		out["ignore_in_test"] = tc.IgnoreInTest
	}
	return out
}

func projectTestEntries(entries []data.TestValueEntryDefinition) []map[string]any {
	out := make([]map[string]any, len(entries))
	for i, e := range entries {
		m := map[string]any{}
		if e.ConfigValue != "" {
			m["config_value"] = e.ConfigValue
		}
		if e.ConfigInclude != nil {
			m["config_include"] = *e.ConfigInclude
		}
		if e.AssertValue != "" {
			m["assert_value"] = e.AssertValue
		}
		// StringValue is the iota zero / default; only emit when overridden.
		if e.ValueType != data.StringValue {
			m["value_type"] = e.ValueType.String()
		}
		out[i] = m
	}
	return out
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
	if c.SupportedVersions != "" {
		return true
	}
	if c.MigrationSource != data.UndefinedMigrationSource || len(c.StateUpgrades) > 0 {
		return true
	}
	if len(c.RelationInfo.ToClasses) > 0 || c.RelationInfo.FromClass != "" || c.RelationInfo.Disabled || c.RelationInfo.Type != data.UndefinedRelationshipType {
		return true
	}
	if len(c.ExcludeChildren) > 0 || len(c.IncludeChildren) > 0 || len(c.IncludeParents) > 0 {
		return true
	}
	if c.Artifacts != nil || len(c.ParentDnVariants) > 0 {
		return true
	}
	if len(c.TestConfig.IgnoreTests) > 0 || c.TestConfig.IgnoreImportStateVerify || len(c.TestConfig.Dependencies) > 0 || c.TestConfig.ReplaceAutoResolved {
		return true
	}
	if c.Documentation.SubCategory != "" || len(c.Documentation.UiLocations) > 0 || len(c.Documentation.DnFormats) > 0 || len(c.Documentation.ExampleParentClasses) > 0 {
		return true
	}
	if len(c.Documentation.Resource.Notes) > 0 {
		return true
	}
	if len(c.Properties) > 0 {
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
	if len(c.IncludeParents) > 0 {
		out["include_parents"] = c.IncludeParents
	}
	if c.SupportedVersions != "" {
		out["supported_versions"] = c.SupportedVersions
	}
	rel := map[string]any{}
	if c.RelationInfo.Disabled {
		rel["disabled"] = c.RelationInfo.Disabled
	}
	if c.RelationInfo.Type != data.UndefinedRelationshipType {
		rel["type"] = c.RelationInfo.Type.String()
	}
	if c.RelationInfo.FromClass != "" {
		rel["from_class"] = c.RelationInfo.FromClass
	}
	if len(c.RelationInfo.ToClasses) > 0 {
		rel["to_classes"] = c.RelationInfo.ToClasses
	}
	if len(rel) > 0 {
		out["relation_info"] = rel
	}
	if c.MigrationSource != data.UndefinedMigrationSource {
		out["migration_source"] = c.MigrationSource.String()
	}
	if len(c.StateUpgrades) > 0 {
		entries := make([]map[string]any, len(c.StateUpgrades))
		for i, su := range c.StateUpgrades {
			entries[i] = projectStateUpgrade(su)
		}
		out["state_upgrades"] = entries
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
	if c.TestConfig.ReplaceAutoResolved {
		testCfg["replace_auto_resolved"] = c.TestConfig.ReplaceAutoResolved
	}
	if len(c.TestConfig.Dependencies) > 0 {
		testCfg["dependencies"] = projectTestDependencies(c.TestConfig.Dependencies)
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
	if len(c.Documentation.Resource.Notes) > 0 {
		doc["resource"] = map[string]any{
			"notes": c.Documentation.Resource.Notes,
		}
	}
	if len(doc) > 0 {
		out["documentation"] = doc
	}
	if len(c.Properties) > 0 {
		props := map[string]any{}
		for name, p := range c.Properties {
			props[name] = projectProperty(p)
		}
		out["properties"] = props
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
	propertiesDir := "gen/scripts/legacy_definitions/properties"
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

	// Pre-load all property YAML files, keyed by basename (e.g. "fvBD.yaml").
	// Entries are removed as the class loop consumes them so the remaining
	// map at the end of the class loop is exactly the orphan property files
	// (classes with no legacy class YAML - 7 comm* files today).
	propFiles, err := filepath.Glob(filepath.Join(propertiesDir, "*.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading properties dir: %v\n", err)
		os.Exit(1)
	}
	propMap := map[string]string{} // basename -> full path
	tally := newKeyTally()
	for _, pf := range propFiles {
		base := filepath.Base(pf)
		// resource_name_overwrite.yaml is the entire-file S3 obsolete drop
		// (section 3): its 7 relation-to-resource-name entries are now expressed
		// as per-class `resource_name` on each relation class. Log a single
		// skip line rather than tallying each entry as an UNKNOWN key.
		// The file itself is deleted in the C22 cleanup commit.
		if base == "resource_name_overwrite.yaml" {
			fmt.Printf("Skipped (S3 obsolete file): %s\n", pf)
			continue
		}
		// properties/global.yaml is the legacy carrier for what is now
		// GlobalMetaDefinition (definitions/global.yaml) - same disposition
		// as classes/global.yaml. Skip from the per-class merge so it does
		// not overwrite the existing GlobalMetaDefinition file.
		if base == "global.yaml" {
			fmt.Printf("Skipped (properties/global.yaml - covered by GlobalMetaDefinition): %s\n", pf)
			continue
		}
		propMap[base] = pf
	}

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

		// Merge the matching property YAML (if any) into the same
		// ClassDefinition before deciding emit-or-skip and before
		// projecting to view.
		propBase := filepath.Base(file)
		if propPath, ok := propMap[propBase]; ok {
			if err := loadAndMigrateProperties(propPath, tally, &canonical); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", propPath, err)
				failed++
				continue
			}
			delete(propMap, propBase)
		}

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

	// Orphan property files: a property YAML whose class YAML never existed
	// (the 7 comm* files today). Each emits a fresh ClassDefinition seeded
	// only from the property loader. Use a sorted iteration order so the
	// per-run "Migrated:" output stays stable.
	orphanBases := make([]string, 0, len(propMap))
	for base := range propMap {
		orphanBases = append(orphanBases, base)
	}
	sort.Strings(orphanBases)
	for _, base := range orphanBases {
		propPath := propMap[base]
		canonical := data.ClassDefinition{}
		if err := loadAndMigrateProperties(propPath, tally, &canonical); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading orphan %s: %v\n", propPath, err)
			failed++
			continue
		}
		if !hasMigratedData(canonical) {
			skipped++
			continue
		}
		view := marshalView(canonical)
		out, err := roundTripStrict(view)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error round-tripping orphan %s: %v\n", propPath, err)
			failed++
			continue
		}
		className := strings.TrimSuffix(base, ".yaml")
		outputPath := filepath.Join(outputDir, className+".yaml")
		if err := os.WriteFile(outputPath, out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			failed++
			continue
		}
		fmt.Printf("Migrated (orphan property): %s -> %s\n", propPath, outputPath)
		migrated++
	}

	fmt.Printf("\nDone: %d migrated, %d skipped (no Phase-1 fields), %d failed\n",
		migrated, skipped, failed)
	tally.print()

	if failed > 0 {
		os.Exit(1)
	}
}

// loadAndMigrateProperties parses one legacy property YAML and merges the
// translated values onto the supplied ClassDefinition. Returns an error
// only on read or parse failure so callers can attribute it to the file.
func loadAndMigrateProperties(file string, tally *keyTally, out *data.ClassDefinition) error {
	payload, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var legacy map[string]any
	if err := yaml.Unmarshal(payload, &legacy); err != nil {
		return err
	}
	migrateProperties(file, legacy, tally, out)
	return nil
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
