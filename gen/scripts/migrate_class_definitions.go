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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
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
	"prop:exclude_targets":             {sectionObsolete, true}, // this commit (C22 - drop, polymorphic auto-detector covers)
	"prop:resource_name_doc_overwrite": {sectionReuse, true},    // this commit (C22 - lived only on properties/global.yaml, already skipped)

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
			// into the generator registry. The canonical pipeline auto-derives
			// the same opt-in from a non-empty IdentifiedBy, so when the meta
			// already carries one the legacy override is redundant - dropping
			// it lets the auto-derive own Artifacts going forward. Otherwise
			// (true with empty IdentifiedBy, or an explicit false) we record
			// the Artifacts list explicitly. fvFBRoute hits the drop branch;
			// fvCrtrn, vmmUplinkPCont, vzAny and fvSiteAssociated keep the
			// explicit list because their IdentifiedBy is empty.
			if b, ok := val.(bool); ok {
				selfClass := strings.TrimSuffix(filepath.Base(file), ".yaml")
				identifiedBy := metaReg.ClassIdentifiedBy[selfClass]
				if b && len(identifiedBy) > 0 {
					fmt.Printf("  DROP: %s include=true is redundant (meta identifiedBy=%v drives auto-derive)\n", file, identifiedBy)
				} else if b {
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

// looksLikeIPOrCIDR returns true when s parses as an IPv4 / IPv6 address or
// a CIDR prefix (e.g. "10.0.0.1", "fe80::1", "2001:db8::/32"). Used to gate
// custom_type bucket migration so only IP-shaped values land in Update; the
// rare non-IP entries (parent_dn-shaped scope, custom-enum ride-alongs)
// would otherwise pollute Update with literals that override dependency
// wiring or duplicate Create harmlessly.
func looksLikeIPOrCIDR(s string) bool {
	if s == "" {
		return false
	}
	if ip, _, err := net.ParseCIDR(s); err == nil && ip != nil {
		return true
	}
	return net.ParseIP(s) != nil
}

// propagateRequiredCreateToUpdate copies TestConfig.Create into Update for
// any Required property that has an explicit Create but no Update. Required
// attributes don't change values across the template's Update step, so the
// asserted Update value must equal Create. Without this copy the loader's
// per-bucket auto-derive would synthesize a different Update value (e.g.
// attr_2 / validValues[1]) that fails apply-time validation for format-
// constrained types (IP, MAC) and silently produces inconsistent assertions
// for required enums. Run as a post-pass after all bucket migrations so
// every property's final Restriction and Create are visible.
func propagateRequiredCreateToUpdate(file string, out *data.ClassDefinition) {
	for propName, prop := range out.Properties {
		if prop.Restriction != data.Required {
			continue
		}
		if len(prop.TestConfig.Create) == 0 || len(prop.TestConfig.Update) > 0 {
			continue
		}
		prop.TestConfig.Update = append(prop.TestConfig.Update, prop.TestConfig.Create...)
		out.Properties[propName] = prop
		fmt.Printf("COPY: %s: %s.test_config.update <- create (required attr, value unchanged across Update step)\n", file, propName)
	}
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// liftTestDefaultsToDefaultValues rewrites the legacy `test_values.default`
// semantic onto the canonical `default_values:` map where they coincide.
// The legacy bucket meant "omit X from HCL config and assert APIC returns
// Y", which the loader expresses through Documentation.DefaultValues:
// generateDefault uses each DefaultValue.Value as the AssertValue and sets
// ConfigInclude=false. The explicit `test_config.default` path instead
// emits ConfigInclude=true (puts Y in the HCL config), which is the
// opposite of the legacy intent.
//
// Per entry, the decision is:
//   - ConfigValue is list-shaped (starts with '['): keep. DefaultValues is
//     string-only and cannot encode list defaults; the explicit
//     test_config.default is the only path for those.
//   - DefaultValues already contains the ConfigValue as a key: drop the
//     test entry as redundant. The auto-derive will assert the same value.
//   - DefaultValues exists but does not contain the ConfigValue: keep both.
//     The legacy author chose a test default value distinct from the
//     declared schema default; preserve that override verbatim.
//   - DefaultValues is empty:
//     - Meta carries a default for this property AND it differs from the
//       ConfigValue: keep. The legacy author intentionally asserted a
//       value that diverges from the APIC server default; lifting would
//       rewrite the "Default:" docs line to a value APIC does not use.
//     - Meta has no default OR meta default equals ConfigValue: lift the
//       ConfigValue into DefaultValues[cv]="" and drop the test entry.
//
// Run as a post-pass after migrateTestValues + the default_values handler
// in loadAndMigrateProperties so both inputs are visible.
func liftTestDefaultsToDefaultValues(file, className string, out *data.ClassDefinition) {
	for propName, prop := range out.Properties {
		if len(prop.TestConfig.Default) == 0 {
			continue
		}
		kept := prop.TestConfig.Default[:0]
		for _, entry := range prop.TestConfig.Default {
			cv := entry.ConfigValue
			if strings.HasPrefix(cv, "[") {
				kept = append(kept, entry)
				continue
			}
			if len(prop.DefaultValues) > 0 {
				if _, exists := prop.DefaultValues[cv]; exists {
					fmt.Printf("DROP: %s: %s.test_config.default[%q] already in default_values\n", file, propName, cv)
					continue
				}
				kept = append(kept, entry)
				continue
			}
			if metaDefault, ok := metaReg.metaDefaultFor(className, propName); ok && metaDefault != cv {
				kept = append(kept, entry)
				continue
			}
			if prop.DefaultValues == nil {
				prop.DefaultValues = map[string]string{}
			}
			prop.DefaultValues[cv] = ""
			fmt.Printf("LIFT: %s: %s.test_config.default[%q] -> default_values\n", file, propName, cv)
		}
		if len(kept) == 0 {
			prop.TestConfig.Default = nil
		} else {
			prop.TestConfig.Default = kept
		}
		out.Properties[propName] = prop
	}
}

// snakeToCamel converts "application_profile_dn" -> "applicationProfileDn".
// Last-resort fallback only: prefer metaReg.resolveMetaName when a class
// context is available so renamed properties resolve to their meta
// camelCase rather than a snake-case round-trip that drops the rename.
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

// reverseUnderscoreOrCamel returns the meta camelCase property name whose
// Underscore form equals snake, by linear scan over metaProps. Falls back
// to snakeToCamel(snake) when no meta property matches, which preserves
// the auto-derived camelCase shape for keys that name properties not
// listed in the supplied meta (e.g. classes whose meta JSON wasn't
// loaded). Distinct from metaRegistry.resolveMetaName: this helper takes
// the meta property list directly and runs no overwrite-inverter steps,
// so it can be called during metaRegistry initialisation before the
// package-level metaReg is assigned. Used by buildOverwritesInverter to
// fix acronym round-trips like "mcast_arp_drop" -> "mcastARPDrop" that
// snakeToCamel alone collapses to the lowercase "mcastArpDrop".
func reverseUnderscoreOrCamel(snake string, metaProps []string) string {
	for _, meta := range metaProps {
		if utils.Underscore(meta) == snake {
			return meta
		}
	}
	return snakeToCamel(snake)
}

// metaRegistry holds cross-file lookup tables used by state_upgrades
// migration to resolve a v1 snake-case attribute name back to the meta
// camelCase property name that keys ClassDefinition.StateUpgrades.
//
// Resolution order is documented on resolveMetaName; the registry pre-loads
// all three sources up front in main() so migrateStateUpgrades stays a pure
// data transform.
type metaRegistry struct {
	// ClassMetaProperties: className -> meta camelCase property names from
	// gen/meta/<className>.json. Used by the reverse-Underscore lookup
	// (default attribute_name = utils.Underscore(metaName)) and by the
	// tn-property friendly-name fallback for named-relation classes.
	ClassMetaProperties map[string][]string
	// ClassIdentifiedBy: className -> meta identifiedBy slice from
	// gen/meta/<className>.json. Drives the legacy `include: true`
	// redundancy check: a non-empty IdentifiedBy already satisfies the
	// canonical artifact auto-derive, so the legacy override is dropped
	// rather than emitted as a verbatim Artifacts list (which would mask
	// the auto-derive going forward).
	ClassIdentifiedBy map[string][]string
	// GlobalAttributeInverter: inverted from
	// gen/definitions/global.yaml.attribute_name_overrides
	// (the canonical store for cross-class renames the legacy generator
	// produced via per-property attribute_name in properties/global.yaml).
	// Maps attribute_name -> meta camelCase.
	GlobalAttributeInverter map[string]string
	// ClassOverwriteInverters: className -> {new_snake -> meta_camel} from
	// each class's per-property overwrites block. Pre-built so child-class
	// inverters are available cross-class (a fvAEPg state_upgrade entry
	// references fvRsBd's overwrites).
	ClassOverwriteInverters map[string]map[string]string
	// ClassMetaDefaults: className -> {metaName -> default-as-string}.
	// Populated from each property's `default` JSON field; only present
	// when the meta carries an explicit default. Used by
	// liftTestDefaultsToDefaultValues to decide whether a legacy
	// test_values.default value diverges from the APIC server default
	// (lift-skip) or matches/has no meta entry (lift-eligible).
	ClassMetaDefaults map[string]map[string]string
	// ClassContainedBy: className -> sorted slice of meta `containedBy`
	// keys (with the `fv:` -> `fv` separator stripped). Used by
	// canParentEntryAutoResolve to decide whether a legacy parents-block
	// entry exactly matches what resolveParentDependencies would auto-
	// generate from meta and can therefore be dropped during migration.
	ClassContainedBy map[string][]string
	// ClassResourceName: className -> resource_name from
	// legacy_definitions/classes/<class>.yaml. Pre-loaded so the parents
	// prune can compare a legacy parent_dn / parent_dependency_name
	// reference against the canonical `aci_<resource_name>.test.id` form
	// the loader would emit.
	ClassResourceName map[string]string
	// ClassRelationToMo: className -> sanitised target class name from
	// meta `relationInfo.toMo` (e.g. `fv:BD` -> `fvBD`). Populated only
	// when toMo is a single class string. Used by canTargetEntryAutoResolve
	// to decide whether a legacy targets-block entry exactly matches the
	// single-target shape resolveTargetDependencies would synthesise.
	ClassRelationToMo map[string]string
	// ClassRnFormat: className -> meta `rnFormat` (e.g. `BD-{name}`).
	// Used by isCanonicalAutoTargetRef to recognise the legacy static DN
	// form `uni/tn-test_name/<rnPrefix><resourceName>_name_<N>` as
	// equivalent to the canonical resource reference the loader would
	// emit after the prune.
	ClassRnFormat map[string]string
	// NoMetaFile: className -> resource_name for classes without meta
	// JSON, loaded from gen/definitions/global.yaml's no_meta_file. Lets
	// the targets prune resolve canonical references for vzBrCP / vzCPIf
	// / monEPGPol / ... whose source-side relation lives in a meta file
	// but whose target class is registered only via no_meta_file.
	NoMetaFile map[string]string
}

// newMetaRegistry pre-loads the three lookup tables. Errors during a single
// file load (malformed JSON, missing top-level key) are logged and the
// entry is skipped so the script remains usable when a meta file is being
// edited.
func newMetaRegistry(metaDir, classesDir, globalDefPath, propertiesDir string) (*metaRegistry, error) {
	reg := &metaRegistry{
		ClassMetaProperties:     map[string][]string{},
		ClassIdentifiedBy:       map[string][]string{},
		GlobalAttributeInverter: map[string]string{},
		ClassOverwriteInverters: map[string]map[string]string{},
		ClassMetaDefaults:       map[string]map[string]string{},
		ClassContainedBy:        map[string][]string{},
		ClassResourceName:       map[string]string{},
		ClassRelationToMo:       map[string]string{},
		ClassRnFormat:           map[string]string{},
		NoMetaFile:              map[string]string{},
	}

	// Meta property names.
	metaFiles, err := filepath.Glob(filepath.Join(metaDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob meta dir %s: %w", metaDir, err)
	}
	for _, mf := range metaFiles {
		raw, err := os.ReadFile(mf)
		if err != nil {
			fmt.Printf("  WARN: failed to read meta file %s: %v\n", mf, err)
			continue
		}
		var top map[string]any
		if err := json.Unmarshal(raw, &top); err != nil {
			fmt.Printf("  WARN: failed to parse meta file %s: %v\n", mf, err)
			continue
		}
		for metaKey, body := range top {
			bodyMap, ok := body.(map[string]any)
			if !ok {
				continue
			}
			propsRaw, ok := bodyMap["properties"].(map[string]any)
			if !ok {
				continue
			}
			className := strings.Replace(metaKey, ":", "", 1)
			names := make([]string, 0, len(propsRaw))
			for name := range propsRaw {
				names = append(names, name)
			}
			sort.Strings(names)
			reg.ClassMetaProperties[className] = names
			defaults := map[string]string{}
			for name, p := range propsRaw {
				propMap, ok := p.(map[string]any)
				if !ok {
					continue
				}
				if raw, ok := propMap["default"]; ok && raw != nil {
					defaults[name] = fmt.Sprintf("%v", raw)
				}
			}
			if len(defaults) > 0 {
				reg.ClassMetaDefaults[className] = defaults
			}
			if idRaw, ok := bodyMap["identifiedBy"].([]any); ok {
				ids := make([]string, 0, len(idRaw))
				for _, v := range idRaw {
					if s, ok := v.(string); ok {
						ids = append(ids, s)
					}
				}
				reg.ClassIdentifiedBy[className] = ids
			}
			if cbRaw, ok := bodyMap["containedBy"].(map[string]any); ok {
				cbs := make([]string, 0, len(cbRaw))
				for pk := range cbRaw {
					cbs = append(cbs, strings.Replace(pk, ":", "", 1))
				}
				sort.Strings(cbs)
				reg.ClassContainedBy[className] = cbs
			}
			if rnFmt, ok := bodyMap["rnFormat"].(string); ok && rnFmt != "" {
				reg.ClassRnFormat[className] = rnFmt
			}
			if riRaw, ok := bodyMap["relationInfo"].(map[string]any); ok {
				if toMo, ok := riRaw["toMo"].(string); ok && toMo != "" {
					reg.ClassRelationToMo[className] = strings.Replace(toMo, ":", "", 1)
				}
			}
		}
	}

	// resource_name from legacy_definitions/classes/*.yaml so the parents
	// prune can resolve the canonical `aci_<resource_name>.test.id` form
	// for arbitrary parent classes.
	classFiles, err := filepath.Glob(filepath.Join(classesDir, "*.yaml"))
	if err == nil {
		for _, cf := range classFiles {
			base := filepath.Base(cf)
			if base == "global.yaml" {
				continue
			}
			raw, err := os.ReadFile(cf)
			if err != nil {
				continue
			}
			var doc map[string]any
			if err := yaml.Unmarshal(raw, &doc); err != nil {
				continue
			}
			if rn, ok := doc["resource_name"].(string); ok && rn != "" {
				reg.ClassResourceName[strings.TrimSuffix(base, ".yaml")] = rn
			}
		}
	}

	// Global attribute_name_overrides.
	globalRaw, err := os.ReadFile(globalDefPath)
	if err == nil {
		var globalDoc map[string]any
		if err := yaml.Unmarshal(globalRaw, &globalDoc); err == nil {
			if overrides, ok := globalDoc["attribute_name_overrides"].(map[any]any); ok {
				for metaKeyAny, attrAny := range overrides {
					metaKey, _ := metaKeyAny.(string)
					attr, _ := attrAny.(string)
					if metaKey != "" && attr != "" {
						reg.GlobalAttributeInverter[attr] = metaKey
					}
				}
			}
			if nmf, ok := globalDoc["no_meta_file"].(map[any]any); ok {
				for cnAny, rnAny := range nmf {
					cn, _ := cnAny.(string)
					rn, _ := rnAny.(string)
					if cn != "" && rn != "" {
						reg.NoMetaFile[cn] = rn
					}
				}
			}
		}
	} else {
		fmt.Printf("  WARN: failed to read global definitions %s: %v\n", globalDefPath, err)
	}

	// Per-class overwrites inverters from properties YAMLs.
	propFiles, err := filepath.Glob(filepath.Join(propertiesDir, "*.yaml"))
	if err == nil {
		for _, pf := range propFiles {
			base := filepath.Base(pf)
			if base == "global.yaml" {
				continue
			}
			raw, err := os.ReadFile(pf)
			if err != nil {
				continue
			}
			var doc map[string]any
			if err := yaml.Unmarshal(raw, &doc); err != nil {
				continue
			}
			inv := buildOverwritesInverter(doc, reg.ClassMetaProperties[strings.TrimSuffix(base, ".yaml")])
			if len(inv) > 0 {
				className := strings.TrimSuffix(base, ".yaml")
				reg.ClassOverwriteInverters[className] = inv
			}
		}
	}

	return reg, nil
}

// resolveMetaName maps a v1 snake-case attribute name to the meta camelCase
// property name used as the state_upgrades map key.
//
// The chain is documented in MIGRATION_OVERVIEW.md section 2.1 and applied
// in this order:
//
//  0. Direct match against the class's meta property names. Callers that
//     pass an already-camelCase key (e.g. default_values keys, which are
//     authored camelCase in v2.19.0) land here without trips through the
//     fallback warning chain.
//  1. Synthetic "parent_dn" -> "parentDn".
//  2. Per-class overwrites inverter - explicit per-property attribute_name
//     renames live here (e.g. fvAEPg.pc_enf_pref -> intra_epg_isolation
//     inverts to intra_epg_isolation -> pcEnfPref).
//  3. Global attribute_name_overrides inverter - cross-class renames
//     authored once globally (e.g. tDn -> target_dn inverts to target_dn
//     -> tDn).
//  4. Reverse-Underscore against the class's meta property names. Most
//     unrenamed properties land here (mode, ipLearning, ...).
//  5. tn-property friendly-name fallback: named-relation classes
//     (fvRsBd, fvRsAEPgMonPol, ...) expose exactly one tn<TargetCap>Name
//     meta property and the legacy generator named the v1 attribute after
//     the target's resource_name + "_name". When the class has exactly one
//     tn<X>Name meta property, return it.
//  6. Last-resort snakeToCamel + WARN log so editorial divergence surfaces.
//
// `ok` is false only when step 6 ran; callers may still use the returned
// value (snakeToCamel form) but should treat the warning as a coverage
// gap.
func (r *metaRegistry) resolveMetaName(className, attrSnake string) (string, bool) {
	if r.classHasMetaProperty(className, attrSnake) {
		return attrSnake, true
	}
	if attrSnake == "parent_dn" {
		return "parentDn", true
	}
	if classInv := r.ClassOverwriteInverters[className]; classInv != nil {
		if meta, ok := classInv[attrSnake]; ok {
			return meta, true
		}
	}
	if meta, ok := r.GlobalAttributeInverter[attrSnake]; ok {
		// Verify the global override actually applies to this class. The
		// global table is keyed by attribute_name across all classes; the
		// inverter can match an attribute name that the class does not
		// expose (e.g. global "ctxName: vrf_name" inverts to "vrf_name ->
		// ctxName" but fvRsCtx has no ctxName meta property).
		if r.classHasMetaProperty(className, meta) {
			return meta, true
		}
	}
	for _, meta := range r.ClassMetaProperties[className] {
		if utils.Underscore(meta) == attrSnake {
			return meta, true
		}
	}
	// Relation-class fallback: the legacy generator named the user-facing
	// inner attribute `<relation_resource_name>_name` even when the meta
	// has no friendly-name property (path-based relations expose tDn). Prefer
	// the single tn<X>Name property when present; fall back to tDn for
	// path-based relations that only carry tDn.
	var tnProps []string
	hasTDn := false
	for _, meta := range r.ClassMetaProperties[className] {
		if strings.HasPrefix(meta, "tn") && strings.HasSuffix(meta, "Name") {
			tnProps = append(tnProps, meta)
		}
		if meta == "tDn" {
			hasTDn = true
		}
	}
	if len(tnProps) == 1 {
		return tnProps[0], true
	}
	if len(tnProps) == 0 && hasTDn && strings.HasSuffix(attrSnake, "_name") {
		return "tDn", true
	}
	return snakeToCamel(attrSnake), false
}

// classHasMetaProperty returns true iff className's meta JSON lists
// metaName under .properties. Linear scan is fine here - meta property
// lists are short (typically <= 30 entries) and resolveMetaName is called
// O(state_upgrade entries) times per run.
func (r *metaRegistry) classHasMetaProperty(className, metaName string) bool {
	for _, p := range r.ClassMetaProperties[className] {
		if p == metaName {
			return true
		}
	}
	return false
}

// metaDefaultFor returns the meta JSON `default` value for the given class
// property as its raw string representation. `ok` is false when the meta
// does not declare a default for this property (most properties).
func (r *metaRegistry) metaDefaultFor(className, metaName string) (string, bool) {
	if defaults, ok := r.ClassMetaDefaults[className]; ok {
		v, found := defaults[metaName]
		return v, found
	}
	return "", false
}

// metaReg is the package-level registry initialised once in main(). Migration
// helpers consult it through resolveMetaName; tests can swap it for a fake.
var metaReg *metaRegistry

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
				upgrade := data.AttributeUpgradeDefinition{
					LegacyAttribute: oldName,
				}
				dotIdx := strings.Index(newName, ".")
				if isSelf {
					// Self-class: newName is the new flat attribute name on
					// the parent class itself. Resolve back to the meta
					// PropertyName for use as the state_upgrades map key.
					camelKey, ok := metaReg.resolveMetaName(selfClass, newName)
					if !ok {
						fmt.Printf("  WARN: %s migration_blocks self attr %q has no meta resolution; falling back to %q\n", file, newName, camelKey)
					}
					entry.Attributes[camelKey] = upgrade
					continue
				}
				// Non-self (child class) entry. Two shapes:
				//   - Dotted "<block>.<inner>": inner property rename inside
				//     the named child class. Key by meta property of the
				//     child class; the block prefix is purely informational.
				//   - Bare "<block>": block-level rename (the legacy flat
				//     attribute became the renamed child block as a whole).
				//     Set legacy_attribute on the child directly, no inner.
				child := entry.Children[className]
				if dotIdx < 0 {
					child.LegacyAttribute = oldName
					entry.Children[className] = child
					continue
				}
				inner := newName[dotIdx+1:]
				camelKey, ok := metaReg.resolveMetaName(className, inner)
				if !ok {
					fmt.Printf("  WARN: %s migration_blocks %s.%s has no meta resolution; falling back to %q\n", file, className, inner, camelKey)
				}
				if child.Attributes == nil {
					child.Attributes = map[string]data.AttributeUpgradeDefinition{}
				}
				child.Attributes[camelKey] = upgrade
				entry.Children[className] = child
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
//
// className is the meta JSON class name (e.g. "fvBD") used by metaReg.
// resolveMetaName to recover the meta camelCase property key from the
// snake-case keys that legacy property YAMLs carry. Without it the
// overwrites / resource_required / default_values / test_values handlers
// would fall back to snakeToCamel and mint phantom property entries for
// acronym-bearing names (mcastARPDrop -> phantom mcastArpDrop), renamed
// properties (addr -> phantom gateway_address), and parent_dn synthetic.
func migrateProperties(file, className string, legacy map[string]any, tally *keyTally, out *data.ClassDefinition) {
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
			// so resolve the legacy OLD snake attribute name (which is the
			// utils.Underscore form of the meta camelCase) back to the
			// meta camelCase via metaReg.resolveMetaName. The value (new
			// attribute name) lands in PropertyDefinition.AttributeName
			// verbatim.
			//
			// Filter: drop the synthetic `relation_from_<x>_to_<y>` ->
			// `relation_to_<y>` aliases. v2.19.0 used these to rename the
			// parent-facing nested attribute that the legacy generator
			// emitted under the full "relation_from_..._to_..." snake. The
			// canonical pipeline computes the same name as Class.
			// ResourceNameNested (`relation_to_<toClass>`, pluralised when
			// IdentifiedBy is non-empty) from the relation's single target
			// class, so the legacy rename target equals the auto-derived
			// nested name. The synthesised property key (relationFromXToY)
			// is not in any class's meta, so the loader's setProperties
			// pass skips it - the entry is pure noise in the canonical
			// YAML and confuses readers who expect every property key to
			// be a real meta property.
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
				if strings.HasPrefix(oldSnake, "relation_from_") && strings.HasPrefix(newName, "relation_to_") {
					fmt.Printf("DROP: %s: overwrites[%s]=%s redundant (canonical Class.ResourceNameNested auto-derives the same name)\n", file, oldSnake, newName)
					continue
				}
				metaName, _ := metaReg.resolveMetaName(className, oldSnake)
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
			// Values are NEW snake attribute names (post-overwrite),
			// e.g. `gateway_address` for fvnsAddrInst.addr. Resolve via
			// metaReg.resolveMetaName so the per-class overwrites inverter
			// step recovers the meta camelCase (`addr`); without this the
			// snake key was written into PropertyDefinition keys verbatim
			// and minted a phantom restriction-only entry alongside the
			// real property.
			for _, snakeName := range asStringSlice(val) {
				metaName, _ := metaReg.resolveMetaName(className, snakeName)
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
			//
			// Keys are usually camelCase meta property names verbatim
			// (`enableRogueExceptMac`) but a handful of legacy files use
			// the synthetic `parent_dn` snake key. Route every key through
			// metaReg.resolveMetaName so step 0 (already-camel match)
			// preserves the common case and step 1 rewrites parent_dn ->
			// parentDn instead of writing the literal snake form.
			//
			// The two `@aci_gen_default_value_overwrite_to_*!` sentinel
			// values (string/list variants) were docs-only hints in the
			// legacy resource.md.tmpl: they told the template to render
			// `Default: ""` (or `Default: []`) instead of the literal
			// sentinel. The canonical pipeline expresses the same intent
			// directly via an empty-key default (DefaultValues{"": ""}),
			// which the loader translates into a `DefaultValue{Value: ""}`
			// docs entry. Translate the sentinel here so the canonical
			// YAML never carries the legacy marker forward.
			dvs, ok := val.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: default_values is not a map: %T\n", file, val)
				continue
			}
			for metaKey, defVal := range dvs {
				rawKey, _ := metaKey.(string)
				if rawKey == "" {
					continue
				}
				metaName, _ := metaReg.resolveMetaName(className, rawKey)
				defStr := fmt.Sprintf("%v", defVal)
				switch defStr {
				case "@aci_gen_default_value_overwrite_to_empty_string!",
					"@aci_gen_default_value_overwrite_to_empty_list!":
					defStr = ""
				}
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
			migrateTestValues(file, className, val, out)

		case "parents":
			migrateParents(file, className, val, out)

		case "targets":
			migrateTargets(file, className, val, out)

		case "exclude_targets":
			// S3 obsolete per MIGRATION_OVERVIEW section 3. The legacy
			// generator's getExcludeTargets subtracted entries from a
			// SetModelTestDependencies-driven auto-union of child-class
			// `targets:` lists. The new pipeline removes both legs of
			// that mechanism: (a) collectChildDrivenDependencies is
			// value-driven so no auto-union exists to subtract from, and
			// (b) the polymorphic-same-type auto-detector (section 8.6)
			// derives parent-class -> target-class matching from
			// Parents intersect ToClasses plus the rendering site. Drop
			// the key with a single log line per file; full editorial
			// verification (catching divergence from the polymorphic
			// rule) is gated on meta JSON loading and lands with the
			// auto-resolution prune commit.
			fmt.Printf("DROP: %s: exclude_targets=%v (polymorphic-same-type auto-detector now handles same-class filtering)\n", file, val)
		}
	}
}

// buildOverwritesInverter returns a map<snake_case_new_attr_name, meta_camel_case>
// derived from the legacy `overwrites` block. Used by metaRegistry init to
// seed ClassOverwriteInverters; runtime handlers resolve attribute keys
// via metaRegistry.resolveMetaName which consults the same table.
// Returns an empty map when overwrites is absent or malformed.
//
// metaProps is the class's meta property list, used by
// reverseUnderscoreOrCamel to recover the meta camelCase from the legacy
// OLD snake key (e.g. "mcast_arp_drop" -> "mcastARPDrop"). When metaProps
// is nil (e.g. class not in any meta JSON) the helper degrades to
// snakeToCamel, matching the pre-fix behaviour.
func buildOverwritesInverter(legacy map[string]any, metaProps []string) map[string]string {
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
		out[newName] = reverseUnderscoreOrCamel(oldSnake, metaProps)
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
// Bucket entries are keyed by snake_case NEW attribute name (post-
// overwrites). metaReg.resolveMetaName handles the resolution chain (step
// 2 per-class overwrite inverter; step 4 reverse-Underscore against the
// class meta; step 5 tn-property friendly-name fallback for relation
// classes) so renamed and acronym-bearing properties land on their meta
// camelCase key rather than minting phantom snakeToCamel entries.
func migrateTestValues(file, className string, val any, out *data.ClassDefinition) {
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
		case "legacy":
			// Legacy alias tests verify the schema-upgrade mapping; any
			// valid value works. The loader's setLegacyTestValues clones
			// Create into Legacy when the legacy alias has a non-divergent
			// type. Drop the explicit migration entirely and let auto-derive
			// handle every case; the loader emits a warning when divergent
			// types prevent cloning so authors can hand-author overrides.
			entries, ok := entry.(map[any]any)
			if !ok {
				continue
			}
			for k := range entries {
				if s, ok := k.(string); ok {
					fmt.Printf("DROP: %s: test_values.legacy[%s] dropped (loader auto-derives Legacy from Create)\n", file, s)
				}
			}
		case "all", "default", "update", "force_new":
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
				metaName, _ := metaReg.resolveMetaName(className, snakeKey)
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
				case "update":
					prop.TestConfig.Update = append(prop.TestConfig.Update, tve)
				case "force_new":
					prop.TestConfig.ForceNew = append(prop.TestConfig.ForceNew, tve)
				}
				out.Properties[metaName] = prop
			}

		case "custom_type":
			// Reuse the legacy custom_type values (mostly IPv6 examples for
			// the IpAddress custom-typed properties) as test_config.update
			// overrides. The Update step then asserts the custom type's
			// semantic equality across format variations (e.g. fe80::0001
			// == fe80::1) and exercises the type plumbing in both Create
			// and Update halves of the same test. Only entries whose value
			// parses as an IP / CIDR are migrated; sibling ride-along values
			// (e.g. fvTrackMember.scope literal DN, vmmDomP.arp_learning
			// enum) are dropped because overlaying them on Update would
			// either override dependency-wired references or duplicate the
			// Create value harmlessly.
			entries, ok := entry.(map[any]any)
			if !ok {
				fmt.Printf("WARN: %s: test_values.custom_type is not a map: %T\n", file, entry)
				continue
			}
			keys := make([]string, 0, len(entries))
			for k := range entries {
				if s, ok := k.(string); ok {
					keys = append(keys, s)
				}
			}
			sort.Strings(keys)
			for _, snakeKey := range keys {
				value := legacyValueToHCL(entries[snakeKey])
				if !looksLikeIPOrCIDR(value) {
					fmt.Printf("DROP: %s: test_values.custom_type[%s]=%q is not IP/CIDR (ride-along, skipped)\n", file, snakeKey, value)
					continue
				}
				metaName, _ := metaReg.resolveMetaName(className, snakeKey)
				if metaName == "" {
					continue
				}
				tve := data.TestValueEntryDefinition{ConfigValue: value}
				prop := upsertProperty(out, metaName)
				prop.TestConfig.Update = append(prop.TestConfig.Update, tve)
				out.Properties[metaName] = prop
				fmt.Printf("REUSE: %s: test_values.custom_type[%s]=%s -> test_config.update (exercise custom IpAddress type)\n", file, snakeKey, value)
			}

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

// canParentEntryAutoResolve returns true when a legacy `parents:` block
// entry for selfClass exactly matches what the loader's
// resolveParentDependencies + buildDependency synthesise from meta
// containedBy. Such an entry can be dropped during migration; the loader
// reproduces the same TestDependency tree (and adds the second-instance
// test_2.id slot used for ForceNew tests).
//
// All of the following must hold:
//
//   - meta containedBy for selfClass has exactly one entry equal to
//     entry.class_name. A single-entry containedBy is unambiguous; multi-
//     parent classes keep their entries verbatim.
//   - entry has no `properties` (would become non-empty ConfigOverrides),
//     no `target_classes` (polymorphic-detector input), and no
//     `class_in_parent`.
//   - entry's `parent_dn`, if present, equals the canonical
//     `aci_<resourceName>.test.id` form. The bare `<resourceName>.test.id`
//     form (legacy quirk) is also accepted: the loader will replace it
//     with the canonical reference, so the prune doubles as a fix-up for
//     those legacy entries.
//   - entry's `parent_dependency`, if present, equals the meta
//     containedBy[0] of entry.class_name (the chain matches the meta
//     walk), and `parent_dependency_name` (if present) matches the same
//     canonical / bare reference shape.
func canParentEntryAutoResolve(selfClass string, entry map[any]any) bool {
	if metaReg == nil {
		return false
	}
	className, _ := entry["class_name"].(string)
	if className == "" {
		return false
	}
	cb := metaReg.ClassContainedBy[selfClass]
	if len(cb) != 1 || cb[0] != className {
		return false
	}
	if props, ok := entry["properties"]; ok && props != nil {
		if m, ok := props.(map[any]any); ok && len(m) > 0 {
			return false
		}
	}
	if _, ok := entry["target_classes"]; ok {
		return false
	}
	if _, ok := entry["class_in_parent"]; ok {
		return false
	}
	if pd, ok := entry["parent_dn"].(string); ok && pd != "" && !isCanonicalAutoRef(className, pd) {
		return false
	}
	if pdep, ok := entry["parent_dependency"].(string); ok && pdep != "" {
		parentCb := metaReg.ClassContainedBy[className]
		if len(parentCb) != 1 || parentCb[0] != pdep {
			return false
		}
		if pdepName, ok := entry["parent_dependency_name"].(string); ok && pdepName != "" && !isCanonicalAutoRef(pdep, pdepName) {
			return false
		}
	} else if pdepName, ok := entry["parent_dependency_name"].(string); ok && pdepName != "" {
		// parent_dependency_name without parent_dependency: migrateParents already
		// drops it with a WARN; the entry itself stays so the test loses nothing.
		return false
	}
	return true
}

// isCanonicalAutoRef returns true when ref is one of the two forms the
// loader's resolveParentDependencies would emit (or accept as equivalent)
// for className: the canonical `aci_<resourceName>.test.id` and the bare
// `<resourceName>.test.id` legacy variant (which the loader will rewrite).
// Returns false when the resource_name is unknown so the parent entry stays
// verbatim.
func isCanonicalAutoRef(className, ref string) bool {
	res := metaReg.ClassResourceName[className]
	if res == "" {
		return false
	}
	return ref == "aci_"+res+".test.id" || ref == res+".test.id"
}

// targetResourceName returns the resource_name for className from any of
// the two registry sources the loader's getResourceNameForClass consults:
// DataStore.Classes (legacy classes/<class>.yaml) and global no_meta_file.
func targetResourceName(className string) string {
	if res := metaReg.ClassResourceName[className]; res != "" {
		return res
	}
	return metaReg.NoMetaFile[className]
}

// isCanonicalAutoTargetRef returns true when ref matches what
// resolveTargetDependencies would emit for className, OR matches the
// legacy static-DN form the prune is willing to replace with the
// canonical reference.
//
// Accepted forms:
//   - `aci_<resourceName>.test.id`     (canonical instance 1)
//   - `aci_<resourceName>.test_2.id`   (canonical instance 2)
//   - `<resourceName>.test.id`         (bare legacy variant)
//   - `uni/tn-test_name/<rnPrefix><resourceName>_name_<N>` for N in 0..2,
//     where rnPrefix is the meta rnFormat with `{name}` stripped (e.g.
//     `BD-` for fvBD). This is the static-reference form that drives the
//     "named relation target uses a static reference" generator
//     diagnostic; we treat it as auto-resolvable because the loader's
//     post-prune output replaces it with `aci_<resourceName>.test.id`.
//   - `uni/tn-test_name/<anything><resourceName>_name_<N>` for N in 0..2
//     when no meta rnFormat is registered (NoMetaFile targets such as
//     vzBrCP / monEPGPol / fvEpRetPol). The legacy fixtures uniformly
//     used the `<resourceName>_name_<N>` token as the rightmost DN
//     segment; the `static: true` flag (already excluded above) covers
//     the rare case where an author pinned a non-conventional instance.
//
// Returns false when the resource_name is unknown (target class not in
// DataStore or NoMetaFile), so the entry stays verbatim.
func isCanonicalAutoTargetRef(className, ref string) bool {
	res := targetResourceName(className)
	if res == "" {
		return false
	}
	if ref == "aci_"+res+".test.id" || ref == "aci_"+res+".test_2.id" || ref == res+".test.id" {
		return true
	}
	const tenantBase = "uni/tn-test_name/"
	if !strings.HasPrefix(ref, tenantBase) {
		return false
	}
	if rnFmt := metaReg.ClassRnFormat[className]; rnFmt != "" {
		idx := strings.Index(rnFmt, "{name}")
		if idx < 0 {
			return false
		}
		expectedPrefix := tenantBase + rnFmt[:idx] + res + "_name_"
		if !strings.HasPrefix(ref, expectedPrefix) {
			return false
		}
		switch ref[len(expectedPrefix):] {
		case "0", "1", "2":
			return true
		}
		return false
	}
	// NoMetaFile fallback: require the rightmost DN segment to be the
	// canonical `<resourceName>_name_<N>` token (the legacy generator
	// uniformly used this form for relation targets).
	last := ref[strings.LastIndex(ref, "/")+1:]
	dash := strings.LastIndex(last, "-")
	if dash < 0 {
		return false
	}
	name := last[dash+1:]
	expected := res + "_name_"
	if !strings.HasPrefix(name, expected) {
		return false
	}
	switch name[len(expected):] {
	case "0", "1", "2":
		return true
	}
	return false
}

// canTargetEntryAutoResolve returns true when a legacy `targets:` block
// entry for selfClass exactly matches what the loader's
// resolveTargetDependencies + buildDependency synthesise from meta
// `relationInfo.toMo` + the target's `containedBy`. Such an entry can be
// dropped during migration; the loader reproduces the same TestDependency
// tree (canonical `aci_<resourceName>.test.id` / `aci_<resourceName>.test_2.id`
// pair, parent chain walked via meta).
//
// All of the following must hold:
//
//   - singleTargetClass is true. The loader's single-target auto-resolve
//     only fires for ToClasses length 1; multi-class target groups
//     (polymorphic, e.g. fvRsSecInherited's fvAEPg/fvESg/l3extInstP) keep
//     every entry verbatim so the polymorphic auto-detector still receives
//     all hints. Same-class repeats (2 vzBrCP entries on fvRsCons) ARE
//     prunable because the loader auto-emits exactly two instances.
//   - selfClass has meta relationInfo.toMo equal to entry.class_name. The
//     migration trusts the meta as the source of truth for the target
//     class identity.
//   - target class has a resource_name (DataStore or NoMetaFile). Without
//     one, getResourceNameForClass returns empty and the loader emits the
//     no-resource diagnostic instead.
//   - entry has no `target_dn_ref` and no `static: true`. target_dn_ref
//     wins over target_dn in migrateTargets, so its presence indicates an
//     authored reference shape we should preserve. `static: true` forces
//     the migrated dep to StaticReference, defeating the auto-resolution
//     this prune relies on.
//   - entry's `parent_dependency`, if present, matches the target's meta
//     containedBy[0]. For NoMetaFile targets (no meta containedBy), any
//     `parent_dependency` is accepted because the loader rebuilds the
//     chain from the target class's own Parents and the legacy parent
//     was redundant with the source class's own parent chain.
//   - entry's `target_dn`, if present, satisfies isCanonicalAutoTargetRef.
//   - entry's `properties`, if present, is either empty or contains only
//     a single `name` override. Single-key name overrides match the
//     legacy cosmetic pattern (`<resourceName>_name_<N>`) that the loader
//     replaces with its `name_1` / `name_2` default; semantic overrides
//     (certificate_chain, gateway_address, alloc_mode, polymorphic name
//     disambiguators, ...) carry additional keys or non-name keys and
//     pin the entry as verbatim.
//
// The legacy keys overwrite_parent_dn_key, shared_classes,
// target_dn_overwrite_docs, parent_dependency_dn_ref, target_classes,
// class_in_parent are NOT pinning signals: migrateTargets either ignores
// them entirely or drops them with a per-file log, so their presence does
// not affect the migrated output we're comparing against.
func canTargetEntryAutoResolve(selfClass string, entry map[any]any, singleTargetClass bool) bool {
	if metaReg == nil || !singleTargetClass {
		return false
	}
	className, _ := entry["class_name"].(string)
	if className == "" {
		return false
	}
	if toMo := metaReg.ClassRelationToMo[selfClass]; toMo == "" || toMo != className {
		return false
	}
	if targetResourceName(className) == "" {
		return false
	}
	if ref, ok := entry["target_dn_ref"].(string); ok && ref != "" {
		return false
	}
	if isStatic, ok := entry["static"].(bool); ok && isStatic {
		return false
	}
	if pdep, ok := entry["parent_dependency"].(string); ok && pdep != "" {
		cb := metaReg.ClassContainedBy[className]
		switch {
		case len(cb) == 0:
			// NoMetaFile target: chain rebuilds from the target's own
			// Parents; the legacy parent_dependency was redundant with
			// the source's parent chain.
		case len(cb) == 1 && cb[0] == pdep:
			// chain matches meta containedBy walk
		default:
			return false
		}
	}
	if td, ok := entry["target_dn"].(string); ok && td != "" && !isCanonicalAutoTargetRef(className, td) {
		return false
	}
	if props, ok := entry["properties"]; ok && props != nil {
		m, ok := props.(map[any]any)
		if !ok {
			return false
		}
		if len(m) > 1 {
			return false
		}
		if len(m) == 1 {
			if _, hasName := m["name"]; !hasName {
				return false
			}
		}
	}
	return true
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
// Auto-resolution filtering: entries whose shape exactly matches what the
// loader's resolveParentDependencies + buildDependency would synthesise
// from meta `containedBy` are dropped here (see canParentEntryAutoResolve).
// Reduces noise in gen/definitions/*.yaml and lets the loader auto-emit
// the bonus second-instance test_2.id slot used for ForceNew tests.
// Entries with ConfigOverrides, non-canonical references, polymorphic
// inputs (target_classes), or off-meta-chain parents are kept verbatim.
func migrateParents(file, selfClass string, val any, out *data.ClassDefinition) {
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
		if canParentEntryAutoResolve(selfClass, entry) {
			fmt.Printf("DROP: %s: parents.%s redundant with meta containedBy auto-resolution\n", file, className)
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
// Auto-resolution filtering: entries whose shape exactly matches what the
// loader's resolveTargetDependencies + buildDependency would synthesise
// from meta `relationInfo.toMo` are dropped here (see
// canTargetEntryAutoResolve). Only fires when the file holds a single
// target entry, mirroring the loader's single-target auto-resolve path;
// multi-target groups keep every entry verbatim so the polymorphic
// renderer still receives all hints. Single-key `name` overrides match
// the legacy cosmetic naming pattern (`<resourceName>_name_<N>`) and are
// accepted as auto-resolvable; the loader replaces them with its
// `name_1` / `name_2` default. Entries with semantic overrides
// (certificate_chain, gateway_address, alloc_mode, polymorphic name
// disambiguators) stay verbatim.
func migrateTargets(file, selfClass string, val any, out *data.ClassDefinition) {
	if val == nil {
		return
	}
	list, ok := val.([]any)
	if !ok {
		fmt.Printf("WARN: %s: targets is not a list: %T\n", file, val)
		return
	}
	// Single-target-class detection: count distinct class_name values across
	// all entries. The prune only fires when the file describes one target
	// class (loader's single-target auto-resolve emits exactly 2 instances)
	// AND has at most 2 entries — preventing data loss for the rare polymorphic
	// or 3+ same-class shapes.
	distinctClasses := map[string]bool{}
	for _, raw := range list {
		entry, ok := raw.(map[any]any)
		if !ok {
			continue
		}
		if cn, _ := entry["class_name"].(string); cn != "" {
			distinctClasses[cn] = true
		}
	}
	singleTargetClass := len(distinctClasses) == 1 && len(list) <= 2
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
		if canTargetEntryAutoResolve(selfClass, entry, singleTargetClass) {
			fmt.Printf("DROP: %s: targets.%s redundant with meta relationInfo auto-resolution\n", file, className)
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
		// Always emit config_value (even when empty) so an explicit
		// legacy `bucket.attr: ""` keeps its empty-string semantics.
		// Dropping it would silently turn an entry that asserts the
		// empty-string roundtrip into an empty `{}` map.
		m["config_value"] = e.ConfigValue
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
	metaDir := "gen/meta"
	globalDefPath := "gen/definitions/global.yaml"

	if len(os.Args) > 1 && os.Args[1] == "clean" {
		runClean(outputDir)
		return
	}

	reg, err := newMetaRegistry(metaDir, classesDir, globalDefPath, propertiesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building meta registry: %v\n", err)
		os.Exit(1)
	}
	metaReg = reg

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
		className := strings.TrimSuffix(filepath.Base(file), ".yaml")
		propBase := filepath.Base(file)
		if propPath, ok := propMap[propBase]; ok {
			if err := loadAndMigrateProperties(propPath, className, tally, &canonical); err != nil {
				fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", propPath, err)
				failed++
				continue
			}
			delete(propMap, propBase)
		}

		// Post-merge cleanup: state_upgrades (from the class YAML) and
		// test_values (from the property YAML) live on different sources
		// but interact - run dedup once both are loaded.
		liftTestDefaultsToDefaultValues(file, className, &canonical)
		propagateRequiredCreateToUpdate(file, &canonical)

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
		orphanClass := strings.TrimSuffix(base, ".yaml")
		canonical := data.ClassDefinition{}
		if err := loadAndMigrateProperties(propPath, orphanClass, tally, &canonical); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading orphan %s: %v\n", propPath, err)
			failed++
			continue
		}
		liftTestDefaultsToDefaultValues(propPath, orphanClass, &canonical)
		propagateRequiredCreateToUpdate(propPath, &canonical)
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
func loadAndMigrateProperties(file, className string, tally *keyTally, out *data.ClassDefinition) error {
	payload, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var legacy map[string]any
	if err := yaml.Unmarshal(payload, &legacy); err != nil {
		return err
	}
	migrateProperties(file, className, legacy, tally, out)
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
