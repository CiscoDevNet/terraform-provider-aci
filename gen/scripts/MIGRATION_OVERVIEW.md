# Definition Migration Overview

Status of every top-level key found across [gen/definitions/classes/](../definitions/classes/), [gen/definitions/properties/](../definitions/properties/), and the two old `global.yaml` files, mapped against the new combined `ClassDefinition` / `PropertyDefinition` / `GlobalMetaDefinition` format defined in [gen/utils/data/definitions.go](../utils/data/definitions.go).

The two legacy directories and `gen/definitions/schema-git-commit-e21fb3e5.json` intentionally retain the paths and contents from `develop`. They are immutable migration inputs: canonical corrections belong in `migrate_class_definitions.go`, while flat per-class `gen/definitions/*.yaml` files are generated outputs (`global.yaml` is the manual exception). Because the migration discovers both legacy directories dynamically, files changed or added on `develop` are included automatically; an unrecognized top-level key fails the migration until its disposition is defined.

`Files` = number of YAML files containing the key. `Target` = field path in the new format (`class.*` = `ClassDefinition`, `property.*` = `PropertyDefinition` under `properties.<metaPropName>`, `global.*` = `GlobalMetaDefinition`).

Categories:
- **Direct mapping** (§1) — value copied verbatim; only the key name or nesting changes.
- **Semantic mapping** (§2) — value/shape needs transformation (rename, restructure, value remap, key remap).
- **Obsolete** (§3) — already covered by the new global file, the new per-class fields, or the existing constants; the old key can be dropped during migration.

Keys with no obvious target in the new schema have been walked through individually; each landed in one of these sections:
- **ADD** (§4) — extend the new schema with a named field and migrate verbatim.
- **REUSE** (§5) — semantics already covered by an existing field; remap and drop the old key.
- **DERIVE** (§6) — drop the YAML key and compute the value in Go from existing data.
- **CONST** (§7) — relocate from YAML into [gen/utils/data/constants.go](../utils/data/constants.go).
- **POSTPONE** (§8) — record and omit the legacy override until its downstream template/test/custom-type work lands.

---

## 1. Direct mapping

Key name unchanged, value unchanged. Only difference is nesting under the new `documentation:` block where applicable.

### Class-level

| Old key | Files | Target |
|---|---:|---|
| `exclude_children` | 6 | `class.exclude_children` |
| `resource_name` | 193 | `class.resource_name` |
| `rn_prepend` | 33 | `class.rn_prepend` |
| `required_as_child` | 1 | `class.required_as_child` |
| `sub_category` | 137 | `class.documentation.sub_category` |
| `ui_locations` | 141 | `class.documentation.ui_locations` |
| `dn_formats` | 22 | `class.documentation.dn_formats` |

### Property-level

| Old key | Files | Target |
|---|---:|---|
| `documentation` (per-property description map) | 122 | `property.documentation.description` |

---

## 2. Semantic mapping

Value or shape needs transformation. Notes describe what each transform does.

### Class-level

| Old key | Files | Target | Transform |
|---|---:|---|---|
| `allow_delete` | 2 | `class.allow_delete` | Value `"false"` → `"never"`; the current provider correction for `fvRsIpslaMonPol` is applied by migration logic without modifying the legacy YAML. |
| `resource_notes` | 11 | `class.documentation.resource.notes` | Rename + nesting move. The new schema splits notes into `documentation.notes` (shared), `documentation.resource.notes` (resource-only), and `documentation.datasource.notes` (datasource-only). Legacy `resource_notes` is resource-only by name; map to `documentation.resource.notes`. |
| `resource_warnings` | 0 v2.19.0 files | `class.documentation.resource.warnings` | Same shape as `resource_notes`. Read by [SetResourceNotesAndWarnigns](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2065) at v2.19.0:2065 but no v2.19.0 YAML file declares it. Loader accepts under the new sibling slot; migration is a no-op today but the loader contract must match the legacy contract for parity. |
| `datasource_notes` | 0 v2.19.0 files | `class.documentation.datasource.notes` | Same shape; read at [v2.19.0:2070](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2070); 0 files use it. Same disposition as `resource_warnings`. |
| `datasource_warnings` | 0 v2.19.0 files | `class.documentation.datasource.warnings` | Same shape; read at [v2.19.0:2075](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2075); 0 files use it. Same disposition as `resource_warnings`. |
| `children` | 7 | `class.include_children` | Rename only (`children` is the older alias). |
| `contained_by` | 37 | `class.include_parents` / `class.exclude_parents` | The old key replaced meta `containedBy`; the canonical fields are additive/subtractive. Remove globally excluded parents, include desired parents absent from meta, and exclude meta parents absent from the desired set. This preserves whether the generated schema exposes `parent_dn`. |
| `class_version` | 1 (`fvRsBDToRelayP`) | `class.supported_versions` | Rename only. |
| `relationship_classes` | 5 | `class.relation_info.to_classes` | Move into nested block; combined with `multi_relationship_class` (size > 1 → drop the flag). For polymorphic-same-type relations (§8.6), union the legacy `parents[].target_classes` values into `to_classes` so the resolved set captures every parent-class scenario (only `fvRsSecInherited` is affected today — `[fvAEPg, fvESg]` becomes `[fvAEPg, fvESg, l3extInstP]`). |
| `migration_version` | 23 | `class.state_upgrades[].prior_schema_version` (+ `class.migration_source: from_sdkv2`) | One `state_upgrades` entry per declared version. See §2.1. |
| `migration_blocks` | 18 | `class.state_upgrades[].attributes` / `.children` | Each `oldKey: newKey` pair becomes an `AttributeUpgradeDefinition` with `legacy_attribute: oldKey` keyed by meta camelCase property name. Dotted values (e.g. `relation_to_static_leafs.deployment_immediacy`) become inner `attributes` / `children` entries. `legacy_type` and `legacy_restriction` remain unset and therefore inherit the current property unless an explicit `type_changes` entry supplies a prior type. See §2.1. |
| `type_changes` | 1 (`netflowMonitorPol`) | `class.state_upgrades[<version>].attributes[<name>].legacy_type` (or `.children[<name>].legacy_type` for block-shape changes) | Move per-attribute type change into the matching `state_upgrades` entry, creating a fresh entry when no `migration_version` lands at that `version`. See §2.1. |

### Property-level

| Old key | Files | Target | Transform |
|---|---:|---|---|
| `overwrites` | 101 | `property.attribute_name` | Old key uses snake_case attribute names (`bd_enforced_enable: bd_enforcement`); new key is the meta camelCase name (`bdEnforcedEnable`) under `properties.<metaName>.attribute_name`. |
| `default_values` | 34 | `property.default_values` | Shape change: old format is `metaPropName: value`; new format is `properties.<metaName>.default_values: { value: versionRange }`. Empty version range string applies to all versions. |
| `resource_required` | 17 | `property.restriction: required` | Each list entry becomes one property with `restriction: required`. |
| `ignores` | 15 | `property.restriction: exclude` (or `global.exclude_properties`) | Class-scoped excludes go to per-property `restriction`; cross-class excludes belong in the new global. |
| `read_only_properties` | 6 | `property.restriction: read_only` | List of meta property names → individual `restriction` entries. |
| `remove_valid_values` | 12 | `property.remove_valid_values` | Shape change: old format is `metaPropName: [values]` at the top of the file; new format is `properties.<metaName>.remove_valid_values: [values]`. |
| `add_valid_values` | 3 | `property.add_valid_values` | Same shape change as `remove_valid_values`. |
| `type_overwrites` | 3 | `property.value_type` | Move per-property; values must mirror the `ValueTypeEnum` vocabulary. |
| `ignore_properties_in_test` | 5 | `property.test_config.ignore_in_test` + `property.example_config.full` | Each entry becomes `ignore_in_test: true`; its authored value is retained separately for the all-attributes example because ignored properties have no normalized test scenario to reuse. |
| `parents` | 66 files / 82 entries | `class.test_config.dependencies[]` (`role: parent`) — **mostly auto-resolved** | The pipeline derives ordinary Parent dependencies from meta `containedBy` and wires `parent_dn`. Entries are retained when they carry a static/non-canonical reference, config overrides, a non-default dependency chain, an off-meta parent, or parent-to-target selection data. `parent_dependency` maps to recursive `dependencies[]`; `parent_dependency_name` maps to its reference; `target_classes` maps to `dependencies[].target_classes`; `class_in_parent` is covered by the recursive shape and is dropped. See §2.2. |
| `targets` | 60 files / 87 entries | `class.test_config.dependencies[]` (`role: target`) — **auto-resolved only for metadata-backed resource targets; otherwise explicit**, see §2.2 | A single target is derived from `Relation.ToClasses[0]` only when the target is a loaded class with a resource artifact. Authored entries remain for definition-only targets, multi-target relations, config overrides, and special dependency chains. `static: true` is the only signal that keeps a literal APIC DN as the test reference. Other retained targets use their explicit `target_dn_ref` or a generated `aci_<resource>.test_<resource>_<index>.id` reference, while `target_dn_overwrite_docs` maps to sparse `properties.tDn.example_config`. `relation_resource_name`, `overwrite_parent_dn_key`, `shared_classes`, and `parent_dependency_dn_ref` are derivable or have no canonical consumer. |
| `test_values` (`all`/`default`/`update`/`legacy`/`custom_type`/`child_*`/`ignore_in_*`/...) | 95 | `property.test_config.{create,default,update,force_new,legacy}` + class child configuration | Old buckets merge into independently defined scenarios while preserving every list member and the distinction between omission and `[]`. Nested `resource_required` supplies the required-only (`default`) scenario. Required and identifying properties can reuse that scenario in future public examples, so a duplicate `example_config.full` is unnecessary. `test_values_for_parent` becomes `class.test_config.embedded_instances` so correlated child instances survive. `child_*` and `ignore_in_*` map to class-level child/dependency configuration. `datasource_required` and `datasource_non_existing` remain derivable (§6); nested IpAddress custom-type values fold into the ordinary scenarios. |

### Global-level

| Old key | Source | Target | Transform |
|---|---|---|---|
| `contained_by_excludes` | `classes/global.yaml` | `global.exclude_parents` | Rename + already present in new global. Drop the old key. |
| `overwrites` | `properties/global.yaml` | `global.attribute_name_overrides` | Old keys are snake_case; new keys are meta camelCase. New global is already populated — drop the old key. |
| `ignores` | `properties/global.yaml` | `global.exclude_properties` | Already present in new global. Drop the old key. |
| `resource_name_doc_overwrite` | `properties/global.yaml` and `properties/<class>.yaml` | `global.documentation_label_overrides` | Already present in new global; the per-class copies (1 file) can be dropped. |

### 2.1 SDKv2 state-upgrade pipeline

The three SDKv2-migration keys (`migration_version`, `migration_blocks`, `type_changes`) translate into a single `state_upgrades:` tree on the new class definition. The legacy YAML carries the name pairs and explicit type changes. The frozen SDKv2 schema remains at [gen/definitions/schema-git-commit-e21fb3e5.json](../definitions/schema-git-commit-e21fb3e5.json), but the current migration script does not read it. For `migration_blocks`, prior type and restriction therefore inherit the current property; `type_changes` is the only current source of an explicit `legacy_type`.

The JSON is retained unchanged with the legacy inputs so a future migration can enrich `legacy_type` or `legacy_restriction` when a concrete compatibility case requires it. It is not referenced by the active runtime pipeline.

**Current migration flow (per attribute touched)**

For each `migration_blocks.<className>.<oldName>: <newName>` entry:

1. **Resolve `<newName>`** to the current meta camelCase property name via `global.attribute_name_overrides` + per-property `attribute_name`. Dotted values (e.g. `relation_to_static_leafs.deployment_immediacy`) become nested `children.<className>.attributes.<innerMetaName>` entries — the scalar-wrap shape that [AttributeUpgradeDefinition.validateChild](../utils/data/definitions.go) explicitly permits (the outer `children` entry carries no `legacy_attribute`; the inner attribute carries the prior flat scalar name).
2. **Emit `legacy_attribute`** and leave `legacy_type` / `legacy_restriction` unset. The renderer and validator interpret those zero values as "inherit from current."
3. **Merge `type_changes` last** by locating the matching `state_upgrades[<version>]` entry and setting `legacy_type` on the named node. The single existing case (`netflowMonitorPol`) has `version: 0` and no `migration_version`, so the script creates a fresh `prior_schema_version: 0` entry.
4. **Leave `legacy_status` unset.** The loader's zero value is `Functioning` (see [`LegacyStatusEnum`](../utils/data/enums.go)), which matches the v2.19.0 "legacy name still exposed" semantic. `frozen` / `removed` are new editorial choices made in migration logic when required.

**`migration_source` — what it means and what it controls**

`migration_source: from_sdkv2` is an enum tag separate from `state_upgrades`. It records *where the resource came from*, not what the upgrade tree looks like. Three consumers today:

| Consumer | Behaviour | Site |
|---|---|---|
| Docs migration warning | When set, `ClassDocumentation.MigrationWarning` is populated with the SDKv2 migration-guide banner. The future active resource-documentation template will render it; the current legacy/static Markdown pages do not consume this normalized field. A framework-native resource gets no warning. | [`setMigrationWarning`](../utils/data/class_documentation.go) |
| Cross-field validator | A non-zero `migration_source` requires at least one `state_upgrades` entry (otherwise the docs would say "migrated" but `UpgradeResourceState` would have nothing to do). | [`validateStateUpgrades`](../utils/data/class.go) |
| Future-proofing | The enum is extensible. Today only `from_sdkv2` is recognised; adding a second source (e.g. `from_terraform_provider_ciscomso`) is a new iota constant + `String()` / `UnmarshalText()` case, no consumer changes needed. | [`MigrationSourceEnum`](../utils/data/enums.go) |

**Orthogonal to `state_upgrades`** by design: a framework-native resource can still grow `state_upgrades` entries for v0→v1 framework-internal schema bumps without ever setting `migration_source`. The two axes:

| | `migration_source` set | `migration_source` unset |
|---|---|---|
| `state_upgrades` present | SDKv2 resource (possibly with later framework-internal bumps) | Framework-native resource that has gone through one or more schema bumps |
| `state_upgrades` absent | **Invalid** — caught by the validator | Framework-native resource, no upgrades ever |

### 2.2 Test-dependency auto-resolution and override criteria

The legacy `parents:` / `targets:` blocks on each `properties/<class>.yaml` were the only source of test-dependency data in v2.19.0 — the generator did not derive them from meta `containedBy` or `relationInfo.toMo`. The new pipeline reverses that: it derives the common case from meta and treats the YAML override as **additive on top of auto-resolution** (or full replacement when the class sets `test_config.replace_auto_resolved: true`). The migration script therefore filters the legacy entries against what auto-resolution produces; entries that auto-resolve correctly are **dropped**, while entries that do not are emitted to `test_config.dependencies[]`.

**What auto-resolves today (no YAML needed)**

| Slot | Auto-resolver | Inputs |
|---|---|---|
| `Class.Parents` (the parent class set) | [setParents](../utils/data/class.go) | `class.include_parents` ∪ (meta `containedBy` minus `class.exclude_parents` ∪ `global.exclude_parents`). Explicit includes take precedence over excludes. |
| Top-level `TestDependency{Role: Parent}` (first **2** parent classes) | [resolveParentDependencies](../utils/data/class.go) | First parent → `aci_<resource>.test.id` + `aci_<resource>.test_2.id` (two instances for ForceNew testing). Second parent → `aci_<resource>.test.id`. Additional parents are skipped — provide them explicitly if needed. |
| Top-level `TestDependency{Role: Target}` (single-target only) | [resolveTargetDependencies](../utils/data/class.go) | `Relation.ToClasses[0]` → `aci_<resource>.test.id` + `aci_<resource>.test_2.id`. **Multi-target relations (`len(Relation.ToClasses) > 1`) raise a diagnostic** unless an explicit `role: target` dependency is declared. |
| Recursive `TestDependency.Dependencies` (the parent's own parents) | [buildDependency](../utils/data/class.go) | For each auto-built dep, recurses through `depClass.Parents` and emits `aci_<resource>.test.id` for each. The `ReferenceTypeEnum` passed here (`ResourceReference` / `DataSourceReference` / `StaticReference`) is honored downstream by `setParentDn` / `setTargetDn` / the placeholder resolvers so static-DN deps render as `StringValue` and Terraform references render as `ReferenceValue`. |
| `parent_dn` attribute test values | [setParentDn](../utils/data/class.go) | Wires `Create` / `Update` / `Default` from the first Parent dep; `ForceNew` from the second. Skipped if the property has an explicit `test_config`. |
| `tDn` / `<target_resource_name>_name` attribute test values | [setTargetDn](../utils/data/class.go) / [setTargetNameProperty](../utils/data/class.go) | Explicit relations expose `tDn` (full DN); Named relations expose `tn<TargetCap>Name` as the meta property, with `AttributeName` renamed to `<target_resource_name>_name` (e.g. `contract_name`) when the target class has a resolvable `resource_name`. Wired from the auto-resolved Target dep. |

**Per-entry decision tree (run during migration)**

For each legacy `parents[]` entry:

1. Resolve `<class_name>`'s meta `containedBy`. If the entry's `class_name` is **not** in the union of meta `containedBy` ∪ existing `include_parents`, append it to `class.include_parents`. Auto-resolution then covers the dep.
2. If `parent_dn` is `aci_<R>.test.id` where `<R>` matches the auto-derived resource name for `class_name`, the test dependency itself is **redundant — drop the entry**. The recursive `parent_dependency` chain is also redundant if it matches the parent's meta `containedBy[0]`.
3. If `parent_dn` is a static literal (e.g. `uni/infra`), emit `test_config.dependencies[]` with `class_name`, `reference: <literal>`, `reference_type: static`, `role: parent`. Keep the dependency: `parentDn` deliberately skips scalar auto-derivation, and [setParentDn](../utils/data/class.go) needs a Parent dependency to populate its normalized Create/Update/Default scenarios. A schema default does not populate these test scenarios. Public-example normalization separately detects an applicable schema default and omits redundant `parent_dn`; that renderer decision does not make the normalized test dependency redundant.
4. If the legacy `parent_dependency` chain differs from meta (e.g. selects a non-default ancestor for the test parent), emit the explicit dep with a recursive `dependencies[]` entry.
5. If an explicitly selected parent falls beyond the two-parent auto-resolution limit, retain it as an explicit dependency.
6. Preserve `target_classes` as `TestDependency.TargetClasses` on the Parent dependency. This is the parent-to-target mapping reserved for class-scoped example rendering of polymorphic relations. The same values are unioned into `relation_info.to_classes` where the legacy relationship list omitted a concrete target. `class_in_parent` is covered by the recursive dependency shape and is dropped.

For each legacy `targets[]` entry:

1. If `len(Relation.ToClasses) > 1`, at least one explicit `role: target` dependency is mandatory. Emit the retained target instances needed by the existing tests/examples; `properties` maps to `config_overrides`.
2. If single-target and the legacy entry has no `properties:` / `static: true` / `overwrite_parent_dn_key` / non-default `parent_dependency_dn_ref`, the entry is **redundant — drop it**. The new pipeline emits the same dep with a resource reference.
3. If `properties:` is present, emit a `config_overrides` map with the same key/value pairs.
4. If `static: true`, emit `reference_type: static` with the literal `target_dn`. Otherwise use `target_dn_ref` when supplied or derive a distinct `aci_<resource>.test_<resource>_<index>.id` reference. A non-static legacy `target_dn` is test fixture data, not proof that the new dependency should remain a literal.
5. If `overwrite_parent_dn_key` is set (31 entries), verify the target class's pipeline-derived `parent_dn` matches. When it doesn't, emit `config_overrides: { <legacy-key>: <reference> }` to preserve the wire shape.
6. Map `target_dn_overwrite_docs` to `properties.tDn.example_config.minimum` and `.full`, separate from the test reference. Drop `relation_resource_name`, `overwrite_parent_dn_key`, `shared_classes`, and `parent_dependency_dn_ref`; retain `target_dn_ref` as the explicit test reference when present.

**Escape hatch — `test_config.replace_auto_resolved: true`**

When a class needs to suppress *all* auto-resolution (e.g. every dep must be a static system DN, or the test exercises a pathological parent chain), set `replace_auto_resolved: true` on `ClassTestConfigDefinition` and declare the full dependency list explicitly. Default is `false` — the migration script should not set this flag; only manual classes opting in to bespoke test scaffolding need it.

**Per-instance child overrides — orthogonal to `dependencies[]`**

Legacy `targets[].properties` was the only knob to push values into a target's HCL block. The new pipeline exposes two separate axes:
- `test_config.dependencies[].config_overrides` overrides scalar properties on a **top-level** dependency resource.
- `test_config.dependencies[].children` and `test_config.children` override values on **nested child blocks** of the dependency or the resource-under-test. These are new and have no legacy counterpart — the migration script does not write them; they are reserved for manual additions.

---

## 3. Obsolete (already migrated or removed)

Functionality already covered by the new global file, by per-class `resource_name` entries, or moved into Go constants in [gen/utils/data/constants.go](../utils/data/constants.go).

| Old key / file | Files | Why obsolete |
|---|---:|---|
| `multi_relationship_class` | 3 | Implicit when `relation_info.to_classes` has more than one entry; no separate flag in the new schema. |
| `classes/global.yaml: contained_by_excludes` | 1 | Identical content already lives under `global.exclude_parents` in the new global. |
| `properties/global.yaml: overwrites` | 1 | Identical content (with key shape changed to meta names) already lives under `global.attribute_name_overrides`. |
| `properties/global.yaml: ignores` | 1 | Identical content already lives under `global.exclude_properties`. |
| `properties/global.yaml: resource_name_doc_overwrite` | 1 | Identical content already lives under `global.documentation_label_overrides`. |
| `definitions/properties/resource_name_overwrite.yaml` (entire file) | 1 | Canonical dependency resolution uses target class definitions directly. The migration validates all 7 legacy entries against the audited map and then drops the obsolete file. |
| `exclude` | 0 v2.19.0 files | Class-level boolean opt-out flag read by [SetClassExclude](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2019) (`m.Exclude` then drove registry filtering at v2.19.0:1411 and 1420). 0 v2.19.0 files declare it. The opt-out direction is covered by `ClassDefinition.Artifacts []ArtifactEnum` (§4.1): `artifacts: []` suppresses resource- and datasource-scoped jobs. Migration is a no-op today; if a future YAML adds `exclude: true`, translate it to `artifacts: []`. |
| `multi_line` | 0 v2.19.0 files | Property-level list of property names whose test values should be wrapped in HCL heredoc (`<<EOT … EOT`) syntax. Read by `LookupTestValue` at [v2.19.0:745](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L745) via `processMultiLine`. 0 v2.19.0 files declare it. The new test renderer is expected to choose HCL formatting based on the value itself (presence of newlines) and per-property `value_type` rather than a YAML allowlist; no migration needed. |
| `required_by_custom_type_in_test` | 0 v2.19.0 files | Property-level list of property names to include in the custom-type test. Read by `IncludeInCustomTypeTest` at [v2.19.0:3681](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3681). No v2.19.0 file declares it, so there is no data to migrate. If specialized custom-type tests are retained, inclusion can be derived from the normalized property `value_type` rather than another allowlist. |
| `exclude_targets` | 4 (`fvAEPg`, `fvESg`, `l3extInstP`, `fvRsSecInherited`) | Legacy `getExcludeTargets` subtracts entries from an automatic cross-product of child target lists. The new pipeline collects only references actually present in resolved child values, so that cross-product no longer exists. Parent-specific choices are preserved positively through `TestDependency.TargetClasses` (§8.6), making the exclusion list redundant. Migration drops all four entries. |
| `custom_type` (nested under `test_values`) | 11 (10 IpAddress + `vmmDomP.arp_learning`) | Legacy `custom_type:` emits a separate custom-type test step. For the 10 IpAddress files, migration folds the author-chosen IP value into the ordinary scenario. `vmmDomP.arp_learning: "disabled"` duplicates its standard value, so it is dropped while `static_custom_type` supplies the explicit `vmm_arp_learning` value type (§8.1). |

---

## 4. ADD — schema additions

Keys decided as ADD: the semantics could not be expressed by a pre-existing field, derived in Go, or relocated to a constant. Each row maps the old key to the field added to the loader structs in [definitions.go](../utils/data/definitions.go).

### 4.1 Decided ADDs

Reviewed and ratified. Per-field reasoning follows the table.

| Old key | Files | New field |
|---|---:|---|
| `include` | 4 | `ClassDefinition.artifacts []ArtifactEnum` |
| `multi_parents` | 2 | `ClassDefinition.parent_dn_variants []ParentDnVariantDefinition` |
| `example_classes` | 8 | `ClassDocumentationDefinition.example_parent_classes []string` |
| `exclude_from_testing` | 2 | `ClassTestConfigDefinition.ignore_tests []IgnoreTestEnum` (value `child`; resource/datasource values have no legacy driver) |
| `ignore_import_state_verify_in_test` | 1 | `ClassTestConfigDefinition.ignore_import_state_verify bool` |
| `documentation` (in `properties/global.yaml`) | 26 entries / 1 file | `GlobalMetaDefinition.PropertyDocumentationOverrides map[string]string` |

### 4.2 Reasoning per decided field

#### `ClassDefinition.artifacts []ArtifactEnum` — was `include`

Files: `fvCrtrn`, `vmmUplinkPCont`, `vzAny`, `fvSiteAssociated` (4) for the migration. The fifth file (`fvFBRoute`) carries `include: true` but is redundant — see the migration-script note in §9. The same field covers the `topSystem` datasource-only case once the render layer lands (see §8.2).

`ArtifactEnum` values are `resource` and `datasource`. For metadata-backed classes, an omitted list means "auto-derive": select both when `IdentifiedBy` is non-empty and neither when it is empty (the legacy `provider.go.tmpl` default). An explicitly empty list opts a metadata-backed class out of both top-level artifact kinds, while a non-empty list selects exactly the listed kinds. Definitions without local APIC metadata are naming lookups only: they omit `artifacts`, remain outside `DataStore.Classes`, and cannot create render jobs.

- Explicit opt-in for classes with empty `IdentifiedBy`. The four migration classes all carry `identifiedBy: []` in meta but ship as public resources today (`aci_epg_useg_block_statement`, `aci_vmm_uplink_container`, `aci_any`, `aci_associated_site`); each becomes `artifacts: [resource, datasource]`. Active model-wrapper generation retains both artifact kinds, and future implementation and example generation will do the same.
- Artifact selection when the metadata-backed default "both" is wrong. A future normalized `topSystem` class would use `[datasource]`; the shape also permits `[resource]` where required.

Not DERIVE: meta has no signal for "empty-identifier class that should be exposed" or "this class is datasource-only" — both are provider policy, not meta data.

Not REUSE on `is_single_nested_when_defined_as_child`: that flag governs nested-attribute *shape*, not registration in `provider.Resources()` / `provider.DataSources()`.

Not split into a boolean opt-in plus a separate selector: the two concerns collapse cleanly into one list — opt-in is "non-empty", selection is "which entries appear" — avoiding the redundancy of `expose_as_top_level: true` + `artifacts: [resource]` for a hypothetical resource-only-with-empty-`IdentifiedBy` class.

Deprecation path: if the generator ever stops keying inclusion off `IdentifiedBy` (e.g., switches to a meta-driven "is configurable" signal), the opt-in dimension collapses; the artifact-selector dimension still has to live somewhere as long as datasource-only resources exist.

#### `ClassDefinition.parent_dn_variants []ParentDnVariantDefinition` — was `multi_parents`

Files: `pkiKeyRing`, `pkiTP` (2). Members: `parent_class string`, `rn_prepend string`, `wrapper_class string`, `test_platform PlatformTypeEnum`.

Both PKI classes have more than one valid placement and each placement reaches APIC through a different request path. `pkiKeyRing`’s meta `dnFormats` show the two shapes directly:

- `uni/userext/pkiext/keyring-{name}` — system-scoped: parent is `pki:Ep`, single direct POST to `api/mo/uni/userext/pkiext/keyring-{name}.json`.
- `uni/tn-{name}/certstore/keyring-{name}` — tenant-scoped: user-facing parent is `fvTenant`, but the request POSTs to `api/mo/uni/tn-{name}/certstore.json` and nests the keyring as a child of an implicit `cloud:CertStore` container. That container is auto-managed by APIC and is never user-addressable.

The generated resource has to branch on the user’s `parent_dn` and pick the right `(api_endpoint, json_envelope)` pair — see [resource_aci_key_ring.go](../../internal/provider/resource_aci_key_ring.go) which today carries the hand-baked literal `wrapperClassMap := map[string]string{"uni/userext/pkiext": "", "certstore": "cloudCertStore"}`. The override exists to drive that branching from data instead of from a hard-coded map.

Members and their purpose:

| Member | Purpose |
|---|---|
| `parent_class` | User-facing parent the variant exposes (e.g., `fvTenant`). It can differ from the direct parent in meta `containedBy` when APIC inserts an implicit wrapper. Its class metadata supplies the valid runtime parent-DN formats. |
| `rn_prepend` | Intermediate RN segment inserted between the selected parent DN and the class RN. |
| `wrapper_class` | Implicit container the request nests the resource inside. Empty for variants that POST against a real, user-addressable parent. |
| `test_platform` | Platform profile (`apic` / `cloud`) for the variant. Future example and test templates can use it to gate variant-specific scenarios. |

Not REUSE on `include_parents`: that field is a flat `[]string`; flattening loses every `rn_prepend` / `wrapper_class` / `test_platform` association.

Not DERIVE in full: the wrapper-container relationship (`fvTenant` → implicit `cloudCertStore` → `pkiKeyRing`) is not modeled in meta JSON — `containedBy` only lists `cloud:CertStore` and `pki:Ep` as direct parents and gives no signal that one of them is auto-created or that `fvTenant` is the user-facing entry point. The variant must therefore name `parent_class`, `rn_prepend`, and `wrapper_class`. Once `parent_class` is known, its valid parent-DN formats are derived from that class through the DataStore rather than duplicated in YAML.

Deprecation path: only viable if the meta file ever exposes per-parent runtime hints (auto-create flags, user-facing-parent indicators). Today’s meta has neither.

#### `ClassDocumentationDefinition.example_parent_classes []string` — was `example_classes`

Files: `fvRsCons`, `fvRsProv`, `fvRsProtBy`, `fvRsConsIf`, `fvRsIntraEpg`, `fvRsSecInherited`, `tagAnnotation`, `tagTag` (8). All are relation or tag classes whose meta `containedBy` set is huge — every `fv:AEPg`, `fv:ESg`, `l3ext:InstP`, `l2ext:InstP`, `mgmt:InstP`, etc. Rendering one example block per parent in the resource/datasource docs would produce dozens of nearly-identical HCL snippets per page.

The legacy templates used this field through `DocumentationExamples` to emit several parent variants. The normalized field is now ordered so a future example pipeline can select its first resolvable parent and produce one representative example per class. The same curated choices (for example `fvAEPg, fvESg` for `fvRsCons`) remain useful as fallback candidates and for future documentation/test consumers.

Not REUSE: there is no existing list field on `ClassDocumentationDefinition` that names a representative subset of `containedBy`. `dn_formats` (the existing override) is a documentation-only flat list of full DN strings, not a parent-class projection — the renderer needs the `ClassName` so it can resolve the corresponding generated resource name (`getResourceName`) for the example HCL.

Not DERIVE: choosing “which 2–3 parents are illustrative” requires editorial judgement (e.g., for `fvRsCons` the meaningful examples are `fvAEPg` and `fvESg`, not `mgmtInstP`). No meta signal ranks parents by representativeness, and a heuristic like “first N from `containedBy`” would render a different and noisier subset than the current docs.

Deprecation path: only when the docs renderer learns to prune `containedBy` to a representative set on its own — e.g., “one per top-level parent family” — with editorial overrides expressed elsewhere.

#### `ClassTestConfigDefinition.ignore_tests []IgnoreTestEnum` — was `exclude_from_testing`

Files: `vmmRsDomMcastAddrNs`, `vmmRsPrefEnhancedLagPol` (2). Both relation classes have empty `IdentifiedBy` and are nested-child-only on the parent resource (`relation_to_multicast_pool` and `relation_to_lacp_enhanced_lag_policy` on `aci_vmm_domain` — see [resource_aci_vmm_domain.go](../../internal/provider/resource_aci_vmm_domain.go)). They never produce a standalone resource and so have no standalone test of their own.

`IgnoreTestEnum` values are `child`, `resource`, `datasource`. Empty/unset list = no skips. Each value suppresses a distinct test surface:

| Enum value | Scope | What it suppresses |
|---|---|---|
| `child` | Parent’s future child-test iteration | One child scenario in every parent’s generated test inputs. The class’s own resource/data-source test selection is unaffected. |
| `resource` | This class’s own resource | The future `resource_aci_<x>_test.go` render job. Other artifact kinds remain unaffected. |
| `datasource` | This class’s own datasource | The future `data_source_aci_<x>_test.go` render job. Other artifact kinds remain unaffected. |

The two migration classes both carry `ignore_tests: [child]`. The setting is normalized now, but no active test template consumes it yet. The `resource` and `datasource` values have no legacy YAML driver; they let future test generation express the orthogonal case where an artifact is wanted but its generated test is genuinely unrunnable (APIC-side race, hardware/license requirement, or non-deterministic ordering).

Why a list and not three booleans: the three scopes are the same axis (test-emission skip at three different render sites) and a class can legitimately combine them, e.g. `[child, resource]` for a relation class whose generated resource test is also flaky. Mirrors the `Artifacts []ArtifactEnum` shape from earlier in this section and keeps the schema growable — a future `import_step` or `replace_step` value drops in without a new field.

Distinct from `ignore_import_state_verify`: future test templates will use that flag to suppress one assertion inside a test that still runs, while `ignore_tests` will suppress whole test jobs or one child scenario.

Why `child` is needed at all on `vmmDomP`’s children — both have a dependency graph that does not fit a single `terraform apply`:

- `vmmRsDomMcastAddrNs` resolves only when the parent’s `enableAVE` attribute is `yes`. The new `Children` override applies to children of a *dependency*, not to siblings of the resource-under-test, so we cannot mutate the parent’s own attributes from a child’s definition.
- `vmmRsPrefEnhancedLagPol` targets `lacpEnhancedLagPol` at `uni/vmmp-VMware/dom-{x}/vswitchpolcont/enlacplagp-{y}`. The target is a grandchild of the same `vmmDomP` that the relation belongs to, producing a self-cycle that needs two applies to set up.

Neither case is expressible in the new dependency graph without (a) sibling-attribute coupling, or (b) multi-step apply support — features that don’t exist today. Excluding the entry from the parent’s child rendering is the only way to keep the parent test green.

Deprecation path: each value retires independently. `child` retires when the test framework grows multi-step apply support and/or sibling-attribute coupling (both classes move to real dependency descriptions). `resource` / `datasource` retire per-class as the underlying flake is fixed.

#### `ClassTestConfigDefinition.ignore_import_state_verify bool` — was `ignore_import_state_verify_in_test`

Files: `vmmDomP` (1).

APIC returns extra non-roundtrip state on `vmmDomP` (server-populated children that aren’t part of the create payload). `ImportStateVerify` would compare those to the planned state and fail. The flag suppresses just the verification step, keeping the import smoke test. It’s a per-class quirk of the wire data; nothing in the meta indicates it.

Deprecation path: only when APIC stops emitting the extra state (or when the framework grows a per-attribute import-verify ignore list).

#### `GlobalMetaDefinition.PropertyDocumentationOverrides map[string]string` — was `documentation` in `properties/global.yaml`

Entries: 26 in `properties/global.yaml` (e.g., `descr: The description of the %s object.`, `nameAlias: The name alias of the %s object.`). Each entry’s `%s` is interpolated with the class’s humanised resource name at render time.

Legacy precedence (`gen/generator.go:2294` v2.19.0): per-class `documentation` override → global `documentation` override → meta `comment` → meta `label`. The global override sits *above* the meta comment and deliberately replaces it. That’s a documented editorial choice — meta comments are inconsistent in wording, capitalisation, and verbosity across the ~180 classes, and the global stub normalises them to a uniform documentation style (e.g., every class’s `annotation` renders as “The annotation of the X object” rather than 171 distinct meta sentences).

The migration scaffold initially omitted the global layer from [setDescription](../utils/data/property_documentation.go), which would have made regenerated documentation diverge from legacy in 26 × N places. The resolved pipeline now applies `PropertyDocumentationOverrides` after the per-class description has been normalized, while preserving the intended precedence: a per-class override wins; otherwise the global value replaces the meta comment/label fallback. Its name aligns with `AttributeNameOverrides` (`attribute_name_overrides`) and `DocumentationLabelOverrides` (`documentation_label_overrides`). The YAML key is renamed from `documentation` during migration so the three globals share the `*_overrides` suffix and the override-not-default semantic is visible in the key itself.

Implemented resolver behavior: [applyGlobalPropertyDocumentationOverrides](../utils/data/property_documentation.go) applies the global lookup without replacing an explicit per-class description. When the global entry contains `%s`, it interpolates the class’s humanised resource name (the same substitution the legacy generator used).

Not DERIVE: the override values are editorial English. There’s no derivation rule that turns 171 distinct meta sentences for `annotation` into one normalised line — that’s a writing decision, not a computation.

Not REUSE: no existing field carries per-meta-property text overrides. Per-class `PropertyDefinition.Documentation.Description` is the natural per-class escape hatch *above* this global, not a replacement (writing 26 overrides into every one of the ~180 class files would multiply the migration footprint by ~180× for no editorial benefit).

Not POSTPONE: the data shape is an unambiguous sibling of two existing globals (`AttributeNameOverrides`, `DocumentationLabelOverrides`); the renderer change is one branch in `setDescription`; and skipping it produces a visible docs regression on regeneration, not a no-op.

Deprecation path: if the meta upstream is ever rewritten to provide consistent property comments, the global overrides can be pruned entry-by-entry — each removal is a single line plus a docs diff review. If every entry ends up removable, the field drops with the same shape it landed in.

---

## 5. REUSE — remap to existing fields

Keys decided as REUSE: semantics match an existing normalized concept, so the migration script rewrites the value into that target. Consumer alignment required by a remap is noted explicitly below.

| Old key | Files | Target existing field | Migration note |
|---|---:|---|---|
| `static_custom_type` | 11 files / 14 entries | `PropertyDefinition.value_type` | `ip_address` is derived from meta `validateAsIPv4OrIPv6` and the redundant override is dropped. The one named entry, `vmmDomP.arpLearning: vmm_arp_learning`, maps to the explicit `ValueTypeEnum.VMMArpLearning` value while retaining the current hand-written custom type implementation. See §8.1 for its later generic-consolidation path. |
| `max_one_class_allowed` | 2 (`fvFBRoute`, `infraRsHPathAtt`) | `ClassDefinition.is_single_nested_when_defined_as_child` | Every template reference to `.MaxOneClassAllowed` already collapses to `or (not .IdentifiedBy) .MaxOneClassAllowed`, which is exactly how `IsSingleNestedWhenDefinedAsChild` is computed in [class.go](../utils/data/class.go) (`ClassDefinition.IsSingleNestedWhenDefinedAsChild \|\| len(IdentifiedBy) == 0`). The `provider.go.tmpl` registry filter (`hasPrefix .RnFormat "rs"`) is also covered: classes with empty `IdentifiedBy` are already excluded by the outer `and .IdentifiedBy` guard. Set `is_single_nested_when_defined_as_child: true` on the two classes and rename the template lookups. |
| `parent_example_dn` | 1 (`vmmDomP`) | Normalized `parentDn` property data | The same `uni/vmmp-VMware` value is already preserved from the legacy property test values in `PropertyDefinition.TestConfig` and is available to future example normalization. Drop the duplicate standalone key. |
| `remove_from_contains` | 1 (`fvRsPathAtt`, value `l2PortSecurityPol`) | `ClassDefinition.exclude_children` | Legacy `SetClassContains` ([v2.19.0:2157](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2157)) builds `m.Contains` from meta `contains` minus the `remove_from_contains` list; `m.Contains` then feeds `DocumentationChildren` ([v2.19.0:3582](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3582)) — the "Children" link list rendered in the resource markdown. Migration sets `exclude_children: [l2PortSecurityPol]` on `fvRsPathAtt`, and the resolved [setChildren](../utils/data/class_documentation.go) path consults the same `ClassDefinition.ExcludeChildren` list when building documentation links. The one field therefore governs both nested-child generation and documentation filtering. |

---

## 6. DERIVE — compute in Go (drop YAML key)

Keys decided as DERIVE: value is a function of data the generator already resolves. The migration script drops the key; `class.go` / `property.go` compute the same answer at codegen time.

| Old key | Files | Derivation rule |
|---|---:|---|
| `data_source_has_no_name_identifier` | 1 (`vzAny`) | `len(class.IdentifiedBy) == 0`. The flag is true for `vzAny` only because it has no naming property — a condition the new generator already exposes. The datasource template branches on `IdentifiedBy` directly. |
| `datasource_required` (nested under `test_values`) | 38 | The datasource lookup config consists of the applicable `parent_dn` plus the renamed `IdentifiedBy` set (snake-case via global `attribute_name_overrides` + per-property `attribute_name`); values come from normalized Default data already populated from `resource_required` and can be reused by future example normalization. Per-file audit (38 files): 36 carry only renamed-`IdentifiedBy` keys with values identical to `resource_required`; the 2 outliers are `fvCrtrn` (empty `IdentifiedBy`, covered by `artifacts: [resource, datasource]` in §4.1) and `vmmDomP` (`parent_dn` against its static parent, preserved directly in normalized `parentDn` data). Migration script: drop the nested block, with a per-file safety check warning when a key is neither `parent_dn` nor in the renamed `IdentifiedBy`, or when its value diverges from `resource_required`. |
| `datasource_non_existing` (nested under `test_values`) | 42 | Sibling of `datasource_required`: same renamed-`IdentifiedBy` key set, but values are a **type-aware non-matching transform** of `resource_required` so the generated datasource test verifies the "no result" branch. Two derivation branches across the 42 files: (a) string-typed naming properties append `_non_existing` (e.g. `criterion` → `criterion_non_existing`, `"131"` → `131_non_existing`); (b) IpAddress-typed properties pick a non-matching IP that still passes the validator (e.g. `10.0.0.2` → `10.0.1.2`, `2.2.2.3` → `2.2.2.4` — typically the next octet). Migration script: drop the nested block, with a per-file safety check warning when a key is not in the renamed `IdentifiedBy` set or the value is not a non-matching transform of `resource_required` for the property's `ValueType`. Exact value equality with the auto-derived candidate is **not** required — the assertion is "non-matching and validator-valid", since the IP increment choice is editorial. The future datasource test renderer must derive a validator-valid non-matching value at codegen time. |
| `static_parent` | 1 (`vmmDomP`) | True iff the resolved `class.test_config.dependencies[]` for the parent role contains a `reference_type: static` entry. Templates branch on the dependency shape instead of a separate flag. |
| `resource_identifier` | 1 (`fvTenant`, value `tn`) | Equal to the RN-prefix segment of `meta.fvTenant.rnFormat` (`tn-{name}` → `tn`). Legacy `GetOverwriteResourceIdentifier` ([v2.19.0:3133](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3133)) is consulted by `GetMultiParentFormats` ([v2.19.0:3096](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3096)) as a fallback when the hardcoded `resourceIdentifier` table misses. The model renderer now identifies a variant from the referenced `ParentClass`'s meta `dnFormats`, so neither the YAML override nor the hardcoded identifier table is needed. |

---

## 7. CONST — relocate to constants.go

Keys decided as CONST: rendering-pipeline tuning, not per-class data. Move into [constants.go](../utils/data/constants.go) and drop from YAML.

| Old key | Source | Constant |
|---|---|---|
| `docs_examples_amount` | `classes/global.yaml` | New `constMaxExamplesToDisplay = 2` (alongside the existing `constMaxDnFormatsToDisplay`). |
| `docs_parent_dn_amount` | `classes/global.yaml` | Already represented by `constMaxParentDnsToDisplay = 20`. Drop the YAML key; the runtime override path is unused. |

---

## 8. POSTPONE — revisit during template/test rework

The remaining POSTPONE keys depend on a downstream test, documentation, or custom-type rewrite. The migration script logs and omits those legacy entries until their consuming layer lands. Resolved items are retained in the numbered context sections below but are no longer listed in this table.

| Old key | Files | Revisit when |
|---|---:|---|
| `class_version_tests` | 1 (`commPol`) | The test templates are rewritten. If the regenerated `commPol` test still needs an independent version filter, ADD `ClassTestConfigDefinition.supported_versions string` (parsed via the existing `Versions` helper). Otherwise drop the key. |
| `datasource_required` (top-level list) | 1 (`topSystem`) | Resolving this requires migrated datasource templates and normalized metadata for `topSystem` — see §8.2. Definition-only classes are available for name lookups but deliberately do not enter `DataStore.Classes`, so they do not create class-scoped render jobs. The live datasource remains the hand-written [internal/provider/data_source_aci_system.go](../../internal/provider/data_source_aci_system.go). |
| `ignore_custom_type_docs` | 6 (`fvAEPg`, `fvRsCons`, `fvRsConsIf`, `fvRsProv`, `mgmtInstP`, `mgmtRsOoBCons`) | The renderer’s rule for `SemanticEquality` valid-values + range docs is decided — see §8.3. Until then, keep the override; once a uniform rule lands the entries collapse to DERIVE (or, less likely, REUSE under a renamed flag). |
| `custom_test_dependency_name` | 1 (`vnsLDevIf`, value `ImportedVnsLDevVipWithFvTenant`) | The test renderer is built and proves it can emit a separate-tenant target subtree via the existing `test_config.dependencies` mechanism — see §8.5. Until then, the hand-written companion HCL in [internal/provider/test_constants.go](../../internal/provider/test_constants.go) stays untouched. |

### 8.1 `vmmDomP.arpLearning` migration — resolved; generic consolidation deferred

`vmm_arp_learning` is the only non-`ip_address` value across all 11 `static_custom_type` files (14 entries). It is implemented today as a hand-written [custom_types/arpLearning.go](../../internal/custom_types/arpLearning.go) that satisfies `basetypes.StringValuableWithSemanticEquals` with a side-table (`"" ≡ "disabled"`), making it a *variant* of the existing `SemanticEquality` value type rather than an unrelated kind of custom type.

The framework’s auto-derived `SemanticEquality` (triggered in [setValueType](../utils/data/property.go) precedence rule 4 when a property has both `ValidValues` and `Validators`) handles the uniform case where wire and human forms map 1:1 via `validValues`. `vmmDomP.arpLearning` falls outside that contract because its equality reads an additional side-table that the meta does not express — the wire value `"disabled"` aliases to the localName `defaultValue` while `0x0`/`0x1` map to `disabled`/`enabled`, and the empty string normalises to `disabled`. The custom type also overrides the meta `uitype: bitmask` (which would auto-derive to `Set`) to render as a single string, captured today through `type_overwrites: arpLearning: string` alongside the `static_custom_type` entry.

The migration is explicit: `static_custom_type: vmm_arp_learning` becomes `value_type: vmm_arp_learning`, and the existing `arpLearning.go` continues to implement the semantics. The ordinary Create/Default scenarios retain the legacy `disabled` value; the old dedicated `test_values.custom_type` bucket adds no distinct coverage and is dropped.

The separate future question is whether all such named semantic-equality variants should use a generic generated mechanism. If that consolidation lands, the `vmm_arp_learning` enum case and bespoke `customtypes/arpLearning.go` can be replaced together without changing this migration's observable provider behavior.

The other 10 nested `custom_type` entries are IpAddress values; their author-chosen IP values fold into the standard scenarios.

### 8.2 `topSystem` / datasource-only artifact context

`aci_system` is the only public datasource without a matching resource. There is no `gen/meta/topSystem.json`, so the class never enters the active DataStore and cannot receive normal class-scoped render jobs. The hand-written [internal/provider/data_source_aci_system.go](../../internal/provider/data_source_aci_system.go), its test, and `docs/data-sources/system.md` survive regeneration because manifest cleanup removes only exact managed paths. The datasource example also remains static, while its standard `provider.tf` is generated through the explicit system provider-example job. The `Required: true` behavior of `system_id` and `pod_id` remains hand-coded in the datasource.

[gen/definitions/properties/topSystem.yaml](../definitions/properties/topSystem.yaml) was authored alongside the manual plugin-framework migration of `aci_system`. The migration script reads it and emits `gen/definitions/topSystem.yaml`, but logs and omits the top-level `datasource_required:` list. The active DataStore cannot consume that emitted definition without corresponding class metadata.

The mechanism for the new pipeline is decided and partially implemented: `ClassDefinition.artifacts []ArtifactEnum` (see §4.2) controls active model wrappers and will control future class-scoped jobs for every metadata-backed class, while definition-only classes are available for name resolution only. `topSystem` therefore needs normalized metadata (or an equivalent normalized class source) before `artifacts: [datasource]` can select a datasource render job. The `Required: true` migration on `system_id`/`pod_id` then collapses into the existing per-property `restriction: required`; no new enum axis is needed on `PropertyDefinition`.

What remains is the class-input and datasource-template build-out:

- Introduce a standard metadata source for `topSystem` (or an explicitly normalized synthetic class) so the emitted definition can enter the DataStore.
- Migrate the datasource schema, implementation, documentation, and test templates. Manifest ownership already preserves the handwritten system files until replacement jobs claim their exact paths.
- Walk through the `topSystem.yaml` content (`documentation` / `overwrites` / `read_only_properties` fold into the standard schema, `datasource_required` collapses into per-property `restriction: required`) before replacing the manual files.

Deprecation path: this entry retires when `topSystem` enters the DataStore and its datasource is fully expressible through the standard templates.

### 8.3 `ignore_custom_type_docs` / `SemanticEquality` valid-values rendering context

Legacy `resource.md.tmpl` renders custom-typed properties with valid values in two branches (v2.19.0:85, 151, 233):

```go-template
{{- if and .HasCustomType (hasCustomTypeDocs .PkgName .Name $.Definitions) (not .ValidateAsIPv4OrIPv6) }}
  - Valid Values:
    *{{ range .ValidValues}} `"{{ . }}"`{{ end}}
    * Or a value in the range of `{{ .min }}` to `{{ .max }}`.
{{- else if .ValidValues }}
  - Valid Values: {{ range .ValidValues}} `"{{ . }}"`{{ end}}
{{- end}}
```

`hasCustomTypeDocs` is a registered template helper backed by [HasCustomTypeDocs](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3768). It returns `true` by default and `false` when the property is listed in the class’s top-level `ignore_custom_type_docs:` list — flipping the docs from the upper branch (with the `Or a value in the range of X to Y` line) to the lower branch (just the localname list).

Why six files carry it for the QoS `Prio` property: the meta validators say `min: 0, max: 9` but the wire `validValues` are non-contiguous (`0,1,2,3,7,8,9` mapped to `level1–level6, unspecified` — no localnames for wire `4/5/6`). The range line is technically truthful (APIC rejects values outside `0..9`) but misleading because users cannot type `"4"`, `"5"`, or `"6"` even though those integers fall inside the validator range.

The override is applied inconsistently. 15 classes ship the same QoS `prio` property with identical meta; only 6 suppress the range line. The other 9 (`fvAp`, `fvESg`, `l3extInstP`, `vzOOBBrCP`, `vzRsAnyToCons`, `vzRsAnyToProv`, `vzRsAnyToConsIf`, `qosDscpClass`, `qosDot1PClass`) keep it. Side-by-side proof in the current docs: [application_profile.md](../../docs/resources/application_profile.md) renders `priority` as `Valid Values: "level1", …, "unspecified". Or a value in the range of 0 to 9.`; [application_epg.md](../../docs/resources/application_epg.md) renders the same property without the range line. There is no editorial reason for the divergence — the 6 fixes correct the docs for some classes; the 9 unfixed classes still mislead.

The new pipeline already auto-derives `ValueType == SemanticEquality` for every `prio` instance via [setValueType](../utils/data/property.go) precedence rule 4 (`len(ValidValues) > 0 && len(Validators) > 0`), so the renderer can detect “this property has both” without per-class hints. The remaining decision is editorial:

- **Uniform suppress** — drop the range line for every `SemanticEquality` property. Cleanest, but loses informative output for `cos` / `dscp` (wire ints equal localnames, range is informative).
- **Localname=wire DERIVE** — keep the range line only when every meta `validValue.value` equals its `localName` (the property is `SemanticEquality` purely because of an `unspecified` sentinel alias). Drops the line for `Prio` (`level1` ≠ wire `3`) and keeps it for `cos` (`0` == wire `0`). Fully derivable from meta, no per-class flag, no inconsistency.
- **Per-property declarative ADD** — keep the override (renamed to something like `PropertyDocumentationDefinition.hide_validator_range bool`). Lets editorial choices stay per-property but perpetuates the 6-of-15 inconsistency unless every QoS-`prio` class is audited.

All three options collapse to either DERIVE or REUSE with a different rendering rule once the renderer is built. The current `ignore_custom_type` ADD proposal was based on the misreading of this flag as “hide a custom-type cross-reference link” (no such link exists in the templates) rather than “suppress one branch of the valid-values rendering for sparse-wire-mapping properties.”

Deprecation path: this entry retires when the docs renderer picks a uniform rule for `SemanticEquality` valid-values + validator-range output. If the rule is meta-derivable (uniform suppress or localname=wire), the override drops entirely (DERIVE). If editorial control is wanted, the override returns under a clearer name covering the actual semantic.

### 8.4 `example_value_overwrite` migration — resolved; rendering pending

Legacy [`LookupTestValue`](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L780) is the generic helper that templates call to render a property value into an example or test step. For relation classes the lookup walks four buckets in order — `testVars["all"]`, `version_mismatch`, `resource_required`, then (for `target_dn` only) `testVars["targets"][i].target_dn_ref` — before falling through to `example_value_overwrite` as a last resort. The legacy author flagged this fallthrough in a self-deprecating comment:

```go
// Referencing is done based on target_dn logic
// This lookup is created as a workaround to reference in an examples on non target_dn attributes
// Redesign of testing / example creation logic should be done to cover this reference use-case
```

A second consumer, `GetTestValueOverwrite` (v2.19.0:878), is hard-coded to `target_dn` and used in three sites of [resource_example_all_attributes.tf.tmpl](../templates/legacy_templates/resource_example_all_attributes.tf.tmpl) (v2.19.0:104, 154, 170). Both consumers read the same per-class `example_value_overwrite:` map keyed by the post-rename snake_case attribute name.

The ten entries split into four groups:

| Group | Files | Override | Why it exists |
|---|---|---|---|
| `target_dn` example projection | `infraRsVipAddrNs`, `infraRsAccBndlSubgrp`, `vmmRsDomMcastAddrNs`, `vmmRsPrefEnhancedLagPol` | `target_dn: aci_<resource>.example.id` | Test dependencies use static or test-labelled references, while the public example deliberately uses the standard `example` resource label. |
| `target_dn` dead code | `infraRsVlanNs` | `target_dn: aci_vlan_pool.example.id` | `targets[0].target_dn_ref: aci_vlan_pool.test_vlan_pool_1.id` is set, so `LookupTestValue` returns that first and never reaches the override. The entry has no consumer today. |
| `parent_dn` example projection | `l3extInstP`, `l3extConsLbl`, `l3extRsRedistributePol` | `parent_dn: aci_l3_outside.example.id` | `l3extOut` has a Terraform name definition but no local APIC metadata. It cannot participate in automatic artifact-backed parent selection, so the example reference is explicit. |
| Editorial | `vnsLDevIf` (`logical_device: aci_l4_l7_device.example_in_another_tenant.id`), `fvTrackMember` (`scope: aci_bridge_domain.example.id`) | Reference to a non-default resource label or a polymorphic-target choice. | `vnsLDevIf` is the imported device and must reference a device in a *different* tenant; the override deliberately uses `example_in_another_tenant`, not `example`. `fvTrackMember.scope` is polymorphic (Bridge Domain or L3Out per its description) and the override picks BD as the canonical example. |

**Current migration contract and intended rendering**

The five `target_dn` entries become sparse `PropertyDefinition.ExampleConfig` data on `tDn`; the three `parent_dn` entries use the same structure on `parentDn`. Target dependencies remain responsible for acceptance-test values. A later example-normalization commit will independently resolve each public reference into renderer-facing values and final class example artifacts.

Arbitrary publicly routable IPv4 values in the immutable legacy test fixtures
are replaced with IANA documentation-range equivalents while their canonical
`TestConfig` entries are constructed. Automatic IP test-value generation uses
the same ranges. Public examples therefore inherit safe addresses from test
data directly instead of applying an output-only replacement table.

The two non-relation editorial entries use semantic property-reference metadata:

- `vnsLDevIf.logical_device` records `classes: [vnsLDevVip]` and `example_label: example_in_another_tenant`.
- `fvTrackMember.scope` records both valid classes (`fvBD`, `l3extOut`) and selects `fvBD` as its example class.

This resolves the old free-form override without treating every ordinary string as a reference. `Property.Reference` describes only semantic DN properties, while sparse `PropertyDefinition.ExampleConfig` handles intentional presentation exceptions without adding output fields to normalized test entries.

The PKI examples carry a related canonical correction. `pkiKeyRing.cert`,
`pkiKeyRing.key`, and `pkiTP.certChain` use module-relative `file()` expressions,
and the two definitions declare their required `example_files`. The generator
copies those shared certificate and private-key assets from
`gen/assets/examples/` into the matching resource-example directories; the
acceptance-test values remain unchanged.

The immutable legacy inputs plus the migration corrections remain authoritative: rerunning `migrate_class_definitions.go` recreates these normalized fields, so later migration runs cannot silently restore the former literal example values.

### 8.5 `vnsLDevIf` / separate-tenant test-target context

`custom_test_dependency_name: ImportedVnsLDevVipWithFvTenant` is a stopgap that prepends a hand-written Go const to every generated test config for `vnsLDevIf`. The legacy template (`gen/templates/legacy_templates/resource_test.go.tmpl`, 10 occurrences) emits `testConfigVnsLDevIfXxx = testConfigImportedVnsLDevVipWithFvTenant + testConfigFvTenantMin + <generated>` across all Min / All / Reset / Children / CustomType / LegacyAttributes variants. The referenced const lives in [internal/provider/test_constants.go](../../internal/provider/test_constants.go) and carries HCL for `aci_tenant.test_tenant_imported_device` + `aci_physical_domain.test` + `aci_l4_l7_device.test_imported_device` — a **separate tenant** subtree so that `vnsLDevIf.logical_device` (the renamed naming property, itself a DN reference to a `vnsLDevVip`) can be tested against an `lDevVip` that lives outside the test’s own tenant.

The legacy generator could not express this case, so the workaround was: write the companion HCL by hand, register its Go const name on the class, and let the test template prepend it verbatim. That’s why the field originally looked like it belonged on `TestDependencyDefinition`, but it isn’t a dependency at all — it’s a directive to splice an external string into the generated test file.

The new pipeline’s [TestDependencyDefinition](../utils/data/definitions.go) is already expressive enough to describe this case structurally: a top-level `vnsLDevVip` entry with `role: target`, `reference_type: resource`, recursive `dependencies: [fvTenant, physDomP]`, and a distinct `reference` value that the renderer turns into a non-default HCL label (e.g. `aci_tenant.test_tenant_imported_device` instead of `aci_tenant.test`). The data model accepts that today; the gap is the renderer:

- The test renderer must emit distinct HCL labels for each `TestDependency.Reference` (not collapse every `fvTenant` to `aci_tenant.test`), so a target dependency can carry its own tenant subtree without colliding with the resource-under-test’s tenant.
- The same renderer must accept a per-dependency override (or derive it from `Reference`) for the Terraform resource label used in HCL.
- Once both behaviours are in place, the `vnsLDevIf` entry collapses to a normal `test_config.dependencies` chain (REUSE) and the hand-written `testConfigImportedVnsLDevVipWithFvTenant` const retires.

**Why this fits POSTPONE**

- The test renderer doesn’t exist yet. Adding `ClassTestConfigDefinition.custom_test_dependency_name` now bakes the stopgap into the new schema and locks the renderer into a string-splice contract that the redesigned `TestDependency` mechanism is meant to replace.
- The case is a single file. Keeping the hand-written test config alongside the existing custom test file [internal/provider/resource_aci_imported_logical_device_test.go](../../internal/provider/resource_aci_imported_logical_device_test.go) costs nothing until the renderer lands.
- The deprecation path is concrete (rewrite as nested `dependencies[]` once the renderer handles distinct HCL labels), so deferring carries no rediscovery risk.

Deprecation path: this entry retires when the test renderer can emit distinct HCL labels per `TestDependency.Reference`. The `vnsLDevIf` YAML then gains a normal `test_config.dependencies` chain and the hand-written `testConfigImportedVnsLDevVipWithFvTenant` const is deleted from [internal/provider/test_constants.go](../../internal/provider/test_constants.go).

### 8.6 Polymorphic same-type relations — explicit parent-target mapping

A *polymorphic same-type* relation can choose a different target class depending on its parent. `fvRsSecInherited` is the current example: under `fvAEPg` it points at another `fvAEPg`, under `fvESg` at another `fvESg`, and under `l3extInstP` at another `l3extInstP`.

The legacy `parents[].target_classes` values are therefore semantic data rather than a redundant detector input. Migration preserves them on each Parent-role `TestDependency.TargetClasses`. It also unions those classes into `Relation.ToClasses` when the legacy relationship list omitted a concrete target.

The explicit mapping is preferable to inferring `ToClasses ⊆ Parents`: it also represents a future asymmetric relation where parent A allows targets `{A, B}` and parent B allows only `{B}`.

```go
TestDependency{
    Class:         fvAEPg,
    Role:          Parent,
    TargetClasses: []*ClassName{fvAEPg},
}
```

**Verification against the 4 multi-target relations today**

| Class | `Parents` (resolved) | `ToClasses` (resolved) | `ToClasses ⊆ Parents` | Polymorphic? |
|---|---|---|:---:|:---:|
| `fvRsDomAtt` | `{fvAEPg}` | `{vmmDomP, physDomP, fcDomP, l2extDomP}` | no | no |
| `infraRsDomP` | `{infraAttEntityP}` | `{vmmDomP, physDomP, fcDomP, l2extDomP}` | no | no |
| `netflowRsExporterToEPg` | `{netflowExporterPol, netflowExporterPolDef}` | `{fvAEPg, l3extInstP, l2extInstP}` | no | no |
| `fvRsSecInherited` | 20 EPG-like classes | `{fvAEPg, fvESg, l3extInstP}` | yes | **yes (3 scenarios)** |

Only `fvRsSecInherited` has the same-type pattern today. The other multi-target relations can still carry `TargetClasses` where their legacy definitions make a parent-specific target choice; for example, `fvRsDomAtt` maps its `fvAEPg` example parent to `vmmDomP`.

**Current consumers**

- Migration retains each legacy mapping and explicit Target dependency, including distinct test labels, and moves intentional public-value differences into property `example_config`.
- The future example normalizer will select one representative parent context, read that Parent dependency's `TargetClasses`, and choose the first matching Target dependency while constructing a renderer-facing class artifact. This keeps parent/target selection out of the template.
- Additional parent and target instances remain in normalized `TestDependencies` for future acceptance-test templates; future public examples will deliberately render only the single representative context.
- `exclude_targets` remains unnecessary because child dependency collection is value-driven; it does not build the legacy cross-product that the exclusion list previously filtered.

**Remaining test-template work**

The acceptance-test renderer may eventually iterate every mapped parent-target scenario rather than wiring only the first two Target dependencies into Create/Update. It can consume the same `TargetClasses` mapping directly, including asymmetric future cases, without adding another schema field or reintroducing `exclude_targets`.

---

## 9. Migration-script work plan

The current [migrate_class_definitions.go](migrate_class_definitions.go) is built out across §§1–7 for both class-level and property-level YAML (the `knownLegacyKeys` allowlist mirrors the §10 audit one-for-one and every entry is `implemented: true`); unresolved §8 POSTPONE entries are dropped during migration with a per-file log line. The normalized example-value and parent-target mapping work described in §§8.4 and 8.6 is implemented. What remains is the downstream acceptance-test and custom-type work for the entries still listed in §8.

Script inputs are the legacy YAML files under [gen/definitions/classes/](../definitions/classes/) and [gen/definitions/properties/](../definitions/properties/). Those directories and the one-shot SDKv2 schema dump [gen/definitions/schema-git-commit-e21fb3e5.json](../definitions/schema-git-commit-e21fb3e5.json) remain unchanged from `develop`. The current script does not consume the JSON; it is retained for future state-upgrade compatibility work. The active loader is non-recursive, so it reads only flat canonical YAML files and ignores the two legacy directories.

Every successful migration rewrites the complete expected flat definition set and then removes stale flat class definitions, excluding `global.yaml`. Definition-only lookup classes are emitted from the migration's compatibility map rather than being added to the legacy directories. Unknown legacy keys fail the command, ensuring an upstream schema addition requires an explicit migration disposition.

### 9.1 Code changes required *before* the script runs end-to-end

The script writes YAML; consuming it requires three groups of upstream changes. **Status: all three groups landed before §9.2 work began — the subsections below document the contract each addition honours so future extensions stay anchored to it.**

**1. Loader struct fields** in [definitions.go](../utils/data/definitions.go):

- `ClassDefinition.Artifacts []ArtifactEnum` (+ `ArtifactEnum` constants `resource`, `datasource`).
- `ClassDefinition.ExampleFiles []string` for shared files copied beside generated resource examples.
- `ClassDefinition.ParentDnVariants []ParentDnVariantDefinition` (+ `ParentDnVariantDefinition` struct with `parent_class`, `rn_prepend`, `wrapper_class`, `test_platform` and `PlatformTypeEnum`).
- `ClassDocumentationDefinition.ExampleParentClasses []string`.
- `ClassTestConfigDefinition.IgnoreTests []IgnoreTestEnum` (+ `IgnoreTestEnum` constants `child`, `resource`, `datasource`).
- `ClassTestConfigDefinition.IgnoreImportStateVerify bool`.
- `GlobalMetaDefinition.PropertyDocumentationOverrides map[string]string`.

**2. Resolver/setter changes** in [class.go](../utils/data/class.go), [class_documentation.go](../utils/data/class_documentation.go), and [property_documentation.go](../utils/data/property_documentation.go). Each new loader field needs a resolved counterpart on the runtime struct plus a `set*` method called from the existing setup chain (following the same pattern as `setIsSingleNestedWhenDefinedAsChild`):

| New loader field | Resolved runtime field | Setter behaviour | Consumer site |
|---|---|---|---|
| `ClassDefinition.Artifacts` | `Class.Artifacts []ArtifactEnum` (same enum on both sides). | For metadata-backed classes, omitted YAML auto-derives `[resource, datasource]` when `len(IdentifiedBy) > 0` and `[]` otherwise. Explicit `artifacts: []` selects neither; a non-empty list selects exactly those kinds. Definition-only classes omit and do not consume this field; they remain outside `DataStore.Classes` and only resolve names. | Artifact-specific call sites use `HasResourceArtifact` or `HasDatasourceArtifact`; shared model jobs intentionally do not filter on artifacts. The static registry infrastructure receives self-registration from the current implementation files; future replacement resource and datasource templates must retain that behavior only for emitted artifacts. |
| `ClassDefinition.ExampleFiles` | No resolved runtime field yet. | Retain the ordered filenames until the example normalizer is introduced. | The future generator will preload each file from `gen/assets/examples/` before cleanup, copy it beside the class's generated resource examples, and record the destination in the managed-files manifest. |
| `ClassDefinition.ParentDnVariants` | `Class.ParentDnVariants []ParentDnVariant` plus a synthesised `Class.DefaultParentDn *ParentDnVariant` for the meta-derived placement. Each `ParentDnVariant` carries the loader fields plus resolved `ParentClass` and optional `WrapperClass` names. | See §9.1.1 for the full derivation and emission rules. | `model.go.tmpl` derives parent-DN match patterns through `ParentClass` and generates `BuildDN`; later resource/schema/test templates can consume the same resolved variants for API routing, schema behavior, and platform gating. Replaces the hand-coded `wrapperClassMap` in [resource_aci_key_ring.go](../../internal/provider/resource_aci_key_ring.go) and the matching block in `resource_aci_certificate_authority.go`. |
| `ClassDocumentationDefinition.ExampleParentClasses` | `ClassDocumentation.ExampleParentClasses []*ClassName` | Resolve each YAML string into a `*ClassName`. When the override is empty, use the normalized `Class.Parents` list so include/exclude rules are retained. | Future example normalization will select the first resolvable parent from this ordered list. Future resource/data-source documentation templates should consume the same normalized example artifacts for HCL and import examples rather than reconstructing parent choices. Legacy templates remain unchanged migration references. |
| `ClassTestConfigDefinition.IgnoreTests` | `Class.TestConfig.IgnoreTests []IgnoreTestEnum` (same enum on both sides). | One-line passthrough from the loader. Empty list = nothing suppressed. | The field is normalized, but no active test template consumes it yet. Future child/resource/data-source test render jobs must gate on the corresponding enum value. The legacy `testvars.yaml.tmpl` remains unchanged and still uses its former fields. Runtime artifacts, schemas, docs, and examples are unaffected; use `Artifacts` to suppress those. |
| `ClassTestConfigDefinition.IgnoreImportStateVerify` | `Class.TestConfig.IgnoreImportStateVerify bool` | One-line passthrough from the loader. | Future resource-test template branch that emits or skips `ImportStateVerify: true`. Distinct from `IgnoreTests: [resource]` — the test still runs; only the one assertion is suppressed. |
| `GlobalMetaDefinition.PropertyDocumentationOverrides` | Consumed directly by the renderer; no per-class resolved counterpart needed. | n/a (lookup is per-property at render time). | The override is applied by [applyGlobalPropertyDocumentationOverrides](../utils/data/property_documentation.go) after per-class [setDescription](../utils/data/property_documentation.go); it interpolates `%s` with the class's humanised resource name (same `GetResourceNameAsDescription` substitution the legacy generator uses) and wins over the meta-comment fallback. |

**3. Constants** in [constants.go](../utils/data/constants.go):

- Add `constMaxExamplesToDisplay = 2` (§7).

#### 9.1.1 `ParentDnVariants` setter and renderer logic

The two affected classes (`pkiKeyRing`, `pkiTP`) each carry one YAML variant; the *default* placement is the meta-derived one. The setter must synthesise the default from meta, validate the variants, and expose a single uniform list for the renderer.

**Migration (YAML key renames done by the script)**

The legacy YAML uses `contained_by` / `test_type`; the new schema uses `parent_class` / `test_platform` to line up with the other `ClassDefinition` fields:

```yaml
# legacy classes/pkiKeyRing.yaml
multi_parents:
  - contained_by: fvTenant
    rn_prepend: certstore
    test_type: cloud
    wrapper_class: cloudCertStore
```

```yaml
# new pkiKeyRing.yaml
parent_dn_variants:
  - parent_class: fvTenant
    rn_prepend: certstore
    test_platform: cloud
    wrapper_class: cloudCertStore
```

**Setter (`setParentDnVariants` on `Class`)**

1. **Derive the default placement from meta.** Walk `meta.dnFormats`; the entry that does not contain any YAML `rn_prepend` segment is the default (for `pkiKeyRing`: `uni/userext/pkiext/keyring-{name}`). Strip the trailing `class.RnFormat` to get the default parent-DN prefix (`uni/userext/pkiext`). Resolve the default parent class from meta `containedBy` after removing the YAML variants' wrapper classes, and store it as `Class.DefaultParentDn`.
2. **Resolve each YAML variant.** Convert `parent_class` and optional `wrapper_class` into `ClassName` values and copy `rn_prepend` and `test_platform`. The variant intentionally does not store a second parent-DN format: its user-facing parent DN is supplied at runtime.
3. **Assemble** `Class.ParentDnVariants` as `[DefaultParentDn] ++ resolved variants`. The default remains first and YAML variants retain source order.

**Model renderer (implemented)**

For every non-default variant, `model.go.tmpl` looks up the referenced
`ParentClass` in the complete DataStore and converts that class's meta
`dnFormats` into anchored match patterns. `BuildDN` checks the exact default
parent DN first, then the derived patterns. A matching non-default variant
inserts its `rn_prepend`; an unmatched value falls back to direct
`parentDN/RN` construction, preserving the existing provider behavior.

For `pkiKeyRing` and `pkiTP`, the `fvTenant` metadata therefore produces a
pattern for `uni/tn-{name}` without adding another YAML property. Runtime
tests cover the default APIC placement, the tenant/cloud placement, and the
fallback for both generated model shapes.

**Downstream consumers**

The later schema template can use `DefaultParentDn` for its default and use
the same parent-class formats for validation. Resource request routing uses
`WrapperClass` to select the implicit APIC envelope, while test templates use
`TestPlatform` to gate placement-specific scenarios. Those consumers should
reuse this resolved variant data rather than reconstructing the legacy
substring map.

**Deprecation path**

The normalized variant data and generated model now provide the replacement for the two `wrapperClassMap` branches and the schema default in [resource_aci_key_ring.go](../../internal/provider/resource_aci_key_ring.go). The current hand-written `resource_aci_key_ring.go` / `resource_aci_certificate_authority.go` adapters retain that logic until their replacement templates are implemented, at which point they should consume the generated model behavior rather than reconstructing the maps.

### 9.2 Script work, by section

All ten steps below are implemented in [migrate_class_definitions.go](migrate_class_definitions.go); the descriptions document the contract each step honours so a future re-run or extension stays anchored to the disposition catalog. Deliberately deferred downstream work is called out inline.

1. **§1 direct mapping** — extend the loader to also copy `rn_prepend`, `required_as_child`, `resource_name`, and `dn_formats` (under `documentation:`), plus the per-property `documentation` map (121 files).
2. **§2 semantic mapping, class-level** — implement the shape-changing transforms: `resource_notes` (+ `resource_warnings` / `datasource_notes` / `datasource_warnings` sibling slots, no v2.19.0 data), `children` → `include_children`, legacy replacement-style `contained_by` → the minimal `include_parents` / `exclude_parents` delta against meta, `class_version` → `supported_versions`, `relationship_classes` (+ drop `multi_relationship_class`) → `relation_info.to_classes`, and `migration_version` + `migration_blocks` + `type_changes` → `state_upgrades` (see §2.1; sets `migration_source: from_sdkv2`, while unspecified `legacy_type` / `legacy_restriction` inherit the current property).
3. **§2 semantic mapping, property-level** — add a second pass that ingests `properties/<class>.yaml` and folds the entries into `NewClassDefinition.Properties` keyed by meta camelCase name. Covers the property-level transforms, the per-class `parents` / `targets` aggregation into class-level `test_config.dependencies[]`, and the `test_values` scenario merge. Preserve `parents[].target_classes` on Parent dependencies, union missing concrete targets into `relation_info.to_classes`, convert non-static target fixtures to resource references, retain explicit public-example differences in `example_config`, and preserve `test_values_for_parent` as `embedded_instances`.
4. **§3 obsolete** — drop the 4 global keys, `multi_relationship_class`, the audited `definitions/properties/resource_name_overwrite.yaml` compatibility map, and `exclude_targets` (4 files) with a per-file log line. The latter is redundant because child dependencies are value-driven and the positive parent-to-target mapping is retained through `TargetClasses` (§8.6).
5. **§4 ADD** — emit the six migrated fields into the output YAML. Verbatim values; one value remap (legacy `exclude_from_testing: true` → `ignore_tests: [child]`) plus the YAML key renames listed in the §4.1 table (e.g. `documentation` in `properties/global.yaml` → `property_documentation_overrides` in `global.yaml`). The `resource` and `datasource` enum values of `ignore_tests` have no legacy driver and so are not written by the script — they exist in the schema as future opt-ins and the loader accepts them when present.
6. **§5 REUSE** — apply the remaps: redundant `static_custom_type: ip_address` drops after a meta-flag assertion, while the named `vmm_arp_learning` value maps explicitly; `max_one_class_allowed` → `is_single_nested_when_defined_as_child`; duplicate `parent_example_dn` drops after its value is preserved in normalized `parentDn` data; `remove_from_contains` → `exclude_children`.
7. **§6 DERIVE** — drop the derivable datasource and identifier keys with a one-line comment per file. Preserve nested `test_values.resource_required` as the property's Default scenario before dropping the old wrapper. Required and identifying properties can use Default in both future public examples, preserving the legacy full-example fallback without duplicating it in `example_config.full`. For `datasource_required` nested, emit a warning when a key isn’t in the renamed `IdentifiedBy` or its value diverges from that required-only scenario. For `datasource_non_existing` nested, emit a warning when a key isn’t in the renamed `IdentifiedBy` or its value is not a non-matching transform. `resource_identifier` remains derivable from meta `rnFormat`.
8. **§7 CONST** — drop both keys from `classes/global.yaml`.
9. **§8 resolution / POSTPONE** — migrate `example_value_overwrite` to dependency or property semantic references. Drop the four still-deferred entries (`class_version_tests`, `datasource_required` top-level, `ignore_custom_type_docs`, `custom_test_dependency_name`) with a per-file log and retain their rationale in §8.
10. **One-off cleanup** — drop the redundant `include: true` from `fvFBRoute` during migration. The class has `IdentifiedBy = ["fbrPrefix"]` and `RnFormat = "pfx-..."`, so the `provider.go.tmpl` registry conditional `or (and .IdentifiedBy (not (and .MaxOneClassAllowed (hasPrefix .RnFormat "rs")))) .Include` is already satisfied by the first branch.

Coverage today: the loader's `UnmarshalStrict` returns clean against every emitted YAML in the corpus; the L1 unknown-key tally is empty and the L2 per-section totals match the disposition catalog.

### 9.3 Open decision points

None block the migration script. The walk-through resolved every key in §4–§7 to a concrete action; the §8 POSTPONE entries are deferred to downstream work, listed here for traceability:

| POSTPONE entry | Blocked on | Owner area |
|---|---|---|
| `datasource_required` top-level (§8.2) | `topSystem` receives a standard metadata source and the datasource schema/implementation/test/docs templates are migrated. Manifest ownership already preserves the handwritten files. | Data model / datasource templates. |
| `ignore_custom_type_docs` (§8.3) | Docs renderer picks a uniform rule for `SemanticEquality` valid-values + validator-range output (uniform suppress, localname=wire DERIVE, or per-property ADD). | Docs renderer / editorial. |
| `custom_test_dependency_name` (§8.5) | Test renderer emits distinct HCL labels per `TestDependency.Reference` so a target dependency can carry its own tenant subtree. | Test renderer. |
| `class_version_tests` (§8 table only) | Test templates are rewritten; re-evaluate whether `commPol` still needs an independent version filter. | Test templates. |

---

## 10. v2.19.0 coverage audit

Proof that §1–§8 catalog every key the v2.19.0 generator accepts. The audit set is the union of two ground-truth sources:

1. **YAML-declared keys** — every distinct top-level key that appears in any file under `v2.19.0:gen/definitions/`.
2. **Code-consumed keys** — every literal in `v2.19.0:gen/generator.go` matched by `key == "<name>"`, plus the nested-map string-literal lookups inside `parents` / `targets` / `multi_parents` / `test_values` entries.

Both lists were extracted with `git show v2.19.0:…` so the audit is reproducible against the tag without checking it out:

```sh
# YAML-declared keys (class scope)
git ls-tree -r v2.19.0 --name-only gen/definitions/classes/ \
  | xargs -I {} sh -c 'git show v2.19.0:{} 2>/dev/null' \
  | grep -E '^[a-z_]+:' | sort -u

# YAML-declared keys (property scope)
git ls-tree -r v2.19.0 --name-only gen/definitions/properties/ \
  | xargs -I {} sh -c 'git show v2.19.0:{} 2>/dev/null' \
  | grep -E '^[a-z_]+:' | sort -u

# Code-consumed top-level keys
git show v2.19.0:gen/generator.go | grep -oE 'key == "[a-z_][a-z_0-9]*"' | sort -u
```

The set difference between code consumers and YAML files surfaced ten accepted-but-unused keys (six top-level, four nested inside `targets` / `test_values`); all ten are now in §2 / §3 (see the **Closing gaps** table below).

### 10.1 Coverage tally

| Source set | Distinct keys at v2.19.0 | Sections that cover them |
|---|---:|---|
| Class-level YAML keys declared in files | 32 | §1, §2, §3, §4, §5, §6, §7, §8 |
| Class-level keys read by code but absent from files | 4 (`resource_warnings`, `datasource_notes`, `datasource_warnings`, `exclude`) | §2 (3) + §3 (1) |
| Property-level YAML keys declared in files | 20 (+ 7 entries inside `definitions/properties/resource_name_overwrite.yaml`) | §1, §2, §3, §4, §5, §6, §8 |
| Property-level keys read by code but absent from files | 2 (`multi_line`, `required_by_custom_type_in_test`) | §3 |
| Global keys (`classes/global.yaml` + `properties/global.yaml`) | 7 | §2 / §3 / §7 |
| Nested keys inside `parents` / `targets` / `multi_parents` entries | 20 | §2 (parents row), §4 (`ParentDnVariants`), §10.3 |
| Nested keys read by code but absent from files | 4 (`target_resource_name`, `version_mismatch`, `child_versions`, `deletable_child`) | §3 (all four — see §10.2) |

**100 % of the v2.19.0 contract** (both declared YAML and accepted-but-unused code paths) is dispositioned. The migration script only acts on keys actually present in files; the audit guarantees that any future YAML adding one of the ten 0-file keys still maps onto a known target.

### 10.2 Closing gaps — keys added during this audit

Each row in the table below is a key that the v2.19.0 generator accepts but no v2.19.0 YAML file declares. They were added to the catalog during the final audit pass; the dispositions match the legacy semantics so the new loader is contract-compatible.

| Key | Legacy consumer | Disposition | Section |
|---|---|---|---|
| `resource_warnings` | [SetResourceNotesAndWarnigns](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2065) (v2.19.0:2065) → `m.ResourceWarnings` | Semantic mapping → `class.documentation.resource.warnings` (sibling of `resource_notes`). | §2 |
| `datasource_notes` | [SetResourceNotesAndWarnigns](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2070) (v2.19.0:2070) → `m.DatasourceNotes` | Semantic mapping → `class.documentation.datasource.notes`. | §2 |
| `datasource_warnings` | [SetResourceNotesAndWarnigns](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2075) (v2.19.0:2075) → `m.DatasourceWarnings` | Semantic mapping → `class.documentation.datasource.warnings`. | §2 |
| `exclude` | [SetClassExclude](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2019) (v2.19.0:2019) → `m.Exclude` (registry filter at v2.19.0:1411, 1420) | Obsolete — covered by `ClassDefinition.Artifacts = []` (§4.1). | §3 |
| `multi_line` | `LookupTestValue` chain at [v2.19.0:745](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L745) (`processMultiLine` heredoc wrap) | Obsolete — new test renderer chooses HCL formatting from the value itself (presence of newlines) plus per-property `value_type`. | §3 |
| `required_by_custom_type_in_test` | [IncludeInCustomTypeTest](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3681) (v2.19.0:3681) | Obsolete — no v2.19.0 definition uses it; any future specialized custom-type test can derive inclusion from the normalized property `value_type`. | §3 |
| `target_resource_name` (nested under `targets[]`) | [GetTestTargetValue](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L366) (v2.19.0:366) — overrides the test-resource label in `aci_<X>.test_<X>_<i>.<attr>` references | Obsolete — 0 v2.19.0 files declare it. The new pipeline derives the test-resource label from the resolved dependency's auto-resolved resource name; no per-target override slot is needed. | §3 |
| `version_mismatch` (nested under `test_values[<prop>]` and `test_values[].children[].<child>`) | [LookupTestValue](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L800) (v2.19.0:800, 929) — alternate test-value bucket selected when the target APIC version doesn't match the property's primary value | Obsolete — 0 v2.19.0 files declare it. The new test renderer expresses version-conditional values via the per-property `TestValueEntry.Versions` field (already on the struct); no nested bucket needed. | §3 |
| `child_versions` (nested under `test_values[].children[]`) | [GetChildVersion](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L231) (v2.19.0:231) — per-child supported-version filter for nested test blocks | Obsolete — 0 v2.19.0 files declare it. The new `class.test_config.children[]` shape carries version filters at the entry level when needed; no separate bucket. | §3 |
| `deletable_child` (nested under `test_values[].children[].<child>`) | [CheckDeletableChild](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L444) (v2.19.0:444, recursive scan) — marks a nested child block as exercising the delete path | Obsolete — 0 v2.19.0 files declare it. The new test renderer derives deletion-testability from the child's own meta `isDeletable` / `allow_delete` flag (mechanism for the deletion test step is part of the test scenario framework rewrite). | §3 |

### 10.3 Nested-key inventory (parents / targets / multi_parents)

Listed here so the parent / target / multi-parent migrations have a single reference. Each nested key already lands in §2 or §4; this is a spot-check, not a separate disposition.

| Outer key | Nested keys read by v2.19.0 | New target |
|---|---|---|
| `parents` ([v2.19.0:3198–3219](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3198)) | `class_name` (entry key), `parent_dependency`, `class_in_parent`, `parent_dependency_name`, `parent_dn`, `target_classes`, `properties` | `class.test_config.dependencies[]` shape (§2 `parents` row). `properties` → `config_overrides`; `target_classes` → `TestDependency.TargetClasses`; `parent_dependency` / `parent_dependency_name` → recursive dependency; `class_in_parent` is derivable. |
| `targets` ([v2.19.0:3282–3320](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3282)) | `class_name` (entry key), `relation_resource_name`, `shared_classes`, `parent_dependency`, `parent_dependency_dn_ref`, `overwrite_parent_dn_key`, `static`, `target_dn`, `target_dn_ref`, `target_dn_overwrite_docs`, `properties` | Folded into `relation_info.to_classes` + `test_config.dependencies[]`. `static: true` preserves a literal reference; otherwise `target_dn_ref` or a generated resource reference supplies the test value. `target_dn_overwrite_docs` → `properties.tDn.example_config`; `properties` → `config_overrides`. The remaining bookkeeping keys are derived or dropped as described in §2.2. |
| `multi_parents` ([v2.19.0:3086–3088](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3086)) | `test_type`, `wrapper_class` | `ClassDefinition.ParentDnVariants[]` (§4.1 `ParentDnVariantDefinition`). |
