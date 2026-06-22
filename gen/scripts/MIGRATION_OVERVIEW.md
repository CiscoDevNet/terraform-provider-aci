# Definition Migration Overview

Status of every top-level key found across [gen/scripts/legacy_definitions/classes/](legacy_definitions/classes/), [gen/scripts/legacy_definitions/properties/](legacy_definitions/properties/), and the two old `global.yaml` files, mapped against the new combined `ClassDefinition` / `PropertyDefinition` / `GlobalMetaDefinition` format defined in [gen/utils/data/definitions.go](../utils/data/definitions.go).

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
- **POSTPONE** (§8) — keep the override; the disposition depends on downstream template/test/custom-type work that hasn’t happened yet.

---

## 1. Direct mapping

Key name unchanged, value unchanged. Only difference is nesting under the new `documentation:` block where applicable.

### Class-level

| Old key | Files | Target |
|---|---:|---|
| `exclude_children` | 6 | `class.exclude_children` |
| `resource_name` | 183 | `class.resource_name` |
| `rn_prepend` | 33 | `class.rn_prepend` |
| `required_as_child` | 1 | `class.required_as_child` |
| `sub_category` | 139 | `class.documentation.sub_category` |
| `ui_locations` | 141 | `class.documentation.ui_locations` |
| `dn_formats` | 22 | `class.documentation.dn_formats` |

### Property-level

| Old key | Files | Target |
|---|---:|---|
| `documentation` (per-property description map) | 123 | `property.documentation.description` |

---

## 2. Semantic mapping

Value or shape needs transformation. Notes describe what each transform does.

### Class-level

| Old key | Files | Target | Transform |
|---|---:|---|---|
| `allow_delete` | 2 | `class.allow_delete` | Value `"false"` → `"never"`. |
| `resource_notes` | 11 | `class.documentation.resource.notes` | Rename + nesting move. The new schema splits notes into `documentation.notes` (shared), `documentation.resource.notes` (resource-only), and `documentation.datasource.notes` (datasource-only). Legacy `resource_notes` is resource-only by name; map to `documentation.resource.notes`. |
| `resource_warnings` | 0 v2.19.0 files | `class.documentation.resource.warnings` | Same shape as `resource_notes`. Read by [SetResourceNotesAndWarnigns](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2065) at v2.19.0:2065 but no v2.19.0 YAML file declares it. Loader accepts under the new sibling slot; migration is a no-op today but the loader contract must match the legacy contract for parity. |
| `datasource_notes` | 0 v2.19.0 files | `class.documentation.datasource.notes` | Same shape; read at [v2.19.0:2070](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2070); 0 files use it. Same disposition as `resource_warnings`. |
| `datasource_warnings` | 0 v2.19.0 files | `class.documentation.datasource.warnings` | Same shape; read at [v2.19.0:2075](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2075); 0 files use it. Same disposition as `resource_warnings`. |
| `children` | 7 | `class.include_children` | Rename only (`children` is the older alias). |
| `contained_by` | 37 | `class.include_parents` | Old key supplied the full `containedBy` list; new key only adds entries on top of meta. Drop entries already present in meta. |
| `class_version` | 1 (`fvRsBDToRelayP`) | `class.supported_versions` | Rename only. |
| `relationship_classes` | 5 | `class.relation_info.to_classes` | Move into nested block; combined with `multi_relationship_class` (size > 1 → drop the flag). For polymorphic-same-type relations (§8.6), union the legacy `parents[].target_classes` values into `to_classes` so the resolved set captures every parent-class scenario (only `fvRsSecInherited` is affected today — `[fvAEPg, fvESg]` becomes `[fvAEPg, fvESg, l3extInstP]`). |
| `migration_version` | 23 | `class.state_upgrades[].prior_schema_version` (+ `class.migration_source: from_sdkv2`) | One `state_upgrades` entry per declared version. See §2.1 for the two-source merge with [schema-git-commit-e21fb3e5.json](legacy_definitions/schema-git-commit-e21fb3e5.json). |
| `migration_blocks` | 18 | `class.state_upgrades[].attributes` / `.children` | Each `oldKey: newKey` pair becomes an `AttributeUpgradeDefinition` with `legacy_attribute: oldKey` keyed by meta camelCase property name. Dotted values (e.g. `relation_to_static_leafs.deployment_immediacy`) become inner `attributes` / `children` entries. The `legacy_type` and `legacy_restriction` fields are read from [schema-git-commit-e21fb3e5.json](legacy_definitions/schema-git-commit-e21fb3e5.json) and omitted when they match the current property. See §2.1. |
| `type_changes` | 1 (`netflowMonitorPol`) | `class.state_upgrades[<version>].attributes[<name>].legacy_type` (or `.children[<name>].legacy_type` for block-shape changes) | Move per-attribute type change into the matching `state_upgrades` entry, creating a fresh entry when no `migration_version` lands at that `version`. See §2.1. |

### Property-level

| Old key | Files | Target | Transform |
|---|---:|---|---|
| `overwrites` | 102 | `property.attribute_name` | Old key uses snake_case attribute names (`bd_enforced_enable: bd_enforcement`); new key is the meta camelCase name (`bdEnforcedEnable`) under `properties.<metaName>.attribute_name`. |
| `default_values` | 34 | `property.default_values` | Shape change: old format is `metaPropName: value`; new format is `properties.<metaName>.default_values: { value: versionRange }`. Empty version range string applies to all versions. |
| `resource_required` | 17 | `property.restriction: required` | Each list entry becomes one property with `restriction: required`. |
| `ignores` | 16 | `property.restriction: exclude` (or `global.exclude_properties`) | Class-scoped excludes go to per-property `restriction`; cross-class excludes belong in the new global. |
| `read_only_properties` | 6 | `property.restriction: read_only` | List of meta property names → individual `restriction` entries. |
| `remove_valid_values` | 12 | `property.remove_valid_values` | Shape change: old format is `metaPropName: [values]` at the top of the file; new format is `properties.<metaName>.remove_valid_values: [values]`. |
| `add_valid_values` | 3 | `property.add_valid_values` | Same shape change as `remove_valid_values`. |
| `type_overwrites` | 3 | `property.value_type` | Move per-property; values must mirror the `ValueTypeEnum` vocabulary. |
| `ignore_properties_in_test` | 5 | `property.test_config.ignore_in_test` | Each entry becomes `ignore_in_test: true` on that property. |
| `parents` | 66 files / 82 entries | `class.test_config.dependencies[]` (`role: parent`) — **mostly auto-resolved** | New pipeline auto-resolves Parent test dependencies from meta `containedBy` via [setParents](../utils/data/class.go) (class.go:524) + [resolveParentDependencies](../utils/data/class.go) (class.go:1233), then auto-wires `parent_dn` via [setParentDn](../utils/data/class.go) (class.go:1420). The legacy entry is **redundant — drop it** when (a) `class_name` is already in meta `containedBy`, (b) `parent_dn` is `aci_<auto-resolved-resource>.test.id`, and (c) the entry's `parent_dependency` chain reproduces via recursive auto-resolution (the parent's own meta `containedBy`). The legacy entry must be **kept** when its `parent_dn` is a static literal (`uni/infra`, etc. — emit with `reference_type: static`), or when its `class_name` is not in meta (add to `class.include_parents`; auto-resolution then covers the rest), or when it carries an explicit `properties` block (emit as `config_overrides`), or when the class needs more than two parent classes (auto-resolution stops at 2). The orthogonal sub-keys map as: `parent_dependency` → recursive `dependencies[]` entry (only if the legacy chain differs from meta); `parent_dependency_name` → `dependencies[].reference`; `target_classes` → drop (filter now lives on `Relation.ToClasses`); `class_in_parent` → drop (covered by the recursive dependency shape). See §2.2 for the per-entry decision tree. |
| `targets` | 53 files / 73 entries | `class.test_config.dependencies[]` (`role: target`) — **mostly auto-resolved (single-target); required for multi-target**, see §2.2 | New pipeline auto-resolves Target deps from `Relation.ToClasses[0]` via [resolveTargetDependencies](../utils/data/class.go) (class.go:1268) and auto-wires `tDn` / `tn<Cap>Name` via [setTargetDn](../utils/data/class.go) (class.go:1488) / [setTargetNameProperty](../utils/data/class.go) (class.go:1555); the named-relation path renames `AttributeName` from the snake_case `tn_<target_cap>_name` fallback to `<target_resource_name>_name` (e.g. `contract_name`) when the target class resolves to a real `resource_name`. The new tests use `aci_<resource>.test.id` references, **superseding** the legacy static system DNs that 77 of the 83 legacy `target_dn` entries carry (67 `uni/...` plus 10 `topology/...`; the remaining 6 entries already use `aci_<resource>.test.id` or are empty). The legacy entry is **redundant — drop it** when (a) the relation is single-target, (b) no `properties:` overrides, (c) `static: false` (or absent), and (d) `overwrite_parent_dn_key` is absent or the target's own pipeline-derived `parent_dn` already matches. The legacy entry must be **kept** when the relation is multi-target (auto-resolution refuses ambiguity; `len(Relation.ToClasses) > 1` is a **hard error** without explicit deps — see class.go:1281), or when it carries `properties` (→ `config_overrides`), or when `static: true` and the target is a true system DN (`uni/infra`, `uni/vmmp-VMware/...`) emitted as `reference_type: static`. Drop `relation_resource_name` (auto-derived), `target_dn_overwrite_docs` (§8.4 POSTPONE), `shared_classes` and `target_dn_ref` (no consumer in new pipeline — confirm during migration). See §2.2 for the per-entry decision tree. |
| `test_values` (`all`/`default`/`update`/`legacy`/`custom_type`/`child_*`/`ignore_in_*`/...) | 95 | `property.test_config.{create,default,update,force_new,legacy}` + `class.test_config.children` | Old buckets merge: `all` → `create`, `default` → `default`. `child_*` and `ignore_in_*` move to class-level `test_config.children` / dependency-scoped overrides. The nested sub-keys land outside this row: `resource_required` is the property-level row above (`property.restriction: required`); `datasource_required` is DERIVE (§6); `datasource_non_existing` is DERIVE (§6, sibling row using the same `IdentifiedBy` + type-aware non-matching transform); nested `custom_type` is Obsolete (§3) — 10 IpAddress files fold their IPv6 author choice into the standard `create:` / `default:` buckets, and the 1 `vmmDomP.arp_learning` entry rides along with §8.1. The `update`, `force_new`, `child_*` and `ignore_in_*` legacy buckets carry no values in any v2.19.0 file — the new-schema slots exist for forward compatibility, but the migration script has no legacy values to translate. |

### Global-level

| Old key | Source | Target | Transform |
|---|---|---|---|
| `contained_by_excludes` | `classes/global.yaml` | `global.exclude_parents` | Rename + already present in new global. Drop the old key. |
| `overwrites` | `properties/global.yaml` | `global.attribute_name_overrides` | Old keys are snake_case; new keys are meta camelCase. New global is already populated — drop the old key. |
| `ignores` | `properties/global.yaml` | `global.exclude_properties` | Already present in new global. Drop the old key. |
| `resource_name_doc_overwrite` | `properties/global.yaml` and `properties/<class>.yaml` | `global.documentation_label_overrides` | Already present in new global; the per-class copies (1 file) can be dropped. |

### 2.1 SDKv2 state-upgrade pipeline

The three SDKv2-migration keys (`migration_version`, `migration_blocks`, `type_changes`) translate into a single `state_upgrades:` tree on the new class definition, but they need a **second input file** to be lossless. The legacy YAML only carries name pairs (old snake_case ↔ new attribute name); the prior framework attribute type and the prior `{required, optional, computed}` triplet live in [gen/scripts/legacy_definitions/schema-git-commit-e21fb3e5.json](legacy_definitions/schema-git-commit-e21fb3e5.json) — a one-shot `terraform providers schema -json` dump of the SDKv2 provider, frozen at git commit `e21fb3e5`. The legacy generator opens it at every run via the [schema-git-commit-* glob](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L1157) and stores the parsed map at `Definitions.Migration`, where [SetMigrationClassTypes](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2517) and [SetLegacyAttributes](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2576) read it.

The migration script inlines the JSON data into each class's `state_upgrades:` tree at migration time, after which **the JSON file is no longer referenced by the runtime pipeline** — `grep -RnE 'schema-git-commit|Definitions.Migration|provider_schemas' gen/utils gen/templates` returns zero hits. The file is preserved under `gen/scripts/legacy_definitions/` so re-running the migration script remains reproducible if future adjustments are needed.

**Two-source merge (per attribute touched)**

For each `migration_blocks.<className>.<oldName>: <newName>` entry:

1. **Resolve `<newName>`** to the current meta camelCase property name via `global.attribute_name_overrides` + per-property `attribute_name`. Dotted values (e.g. `relation_to_static_leafs.deployment_immediacy`) become nested `children.<className>.attributes.<innerMetaName>` entries — the scalar-wrap shape that [AttributeUpgradeDefinition.validateChild](../utils/data/definitions.go) explicitly permits (the outer `children` entry carries no `legacy_attribute`; the inner attribute carries the prior flat scalar name).
2. **Look up `<oldName>` in the JSON** under `provider_schemas["registry.terraform.io/ciscodevnet/aci"].resource_schemas["aci_" + <resource_name>].block`:
   - `block.attributes[<oldName>].type` → mapped to `LegacyAttributeTypeEnum` (e.g. JSON `"string"` → `string_attribute`; `["set", ["object", {…}]]` → `set_nested_attribute`).
   - `block.attributes[<oldName>].{required, optional, computed}` → mapped to `RestrictionEnum`.
   - For renamed blocks, the equivalent fields under `block.block_types[<oldName>]` (`nesting_mode` → `list_nested_attribute` / `set_nested_attribute`).
3. **Compare to the current resolved `Property`** and write `legacy_type` / `legacy_restriction` only when they diverge from the current attribute. Both fields are documented as "Omit when the prior type matches the current type" in [StateUpgradeDefinition](../utils/data/definitions.go); the renderer and validator treat the zero value as "inherit from current." For the common case of a pure name rename with no type or restriction change, the resulting node carries `legacy_attribute` only.
4. **Merge `type_changes` last** by locating the matching `state_upgrades[<version>]` entry and setting `legacy_type` on the named node. The single existing case (`netflowMonitorPol`) has `version: 0` and no `migration_version`, so the script creates a fresh `prior_schema_version: 0` entry.
5. **Leave `legacy_status` unset.** The loader's zero value is `Functioning` (see [LegacyStatusEnum](../utils/data/enums.go) at enums.go:89), which matches the v2.19.0 "legacy name still exposed" semantic. `frozen` / `removed` are new editorial choices — author them by hand after migration when a class wants to freeze or drop a legacy alias.

**`migration_source` — what it means and what it controls**

`migration_source: from_sdkv2` is an enum tag separate from `state_upgrades`. It records *where the resource came from*, not what the upgrade tree looks like. Three consumers today:

| Consumer | Behaviour | Site |
|---|---|---|
| Docs migration warning | When set, [ClassDocumentation.MigrationWarning](../utils/data/class_documentation.go) is populated with the SDKv2 migration-guide banner; the resource docs template renders it. A framework-native resource gets no warning. | [class_documentation.go:348](../utils/data/class_documentation.go) |
| Cross-field validator | A non-zero `migration_source` requires at least one `state_upgrades` entry (otherwise the docs would say "migrated" but `UpgradeResourceState` would have nothing to do). | [class.go:2180](../utils/data/class.go) `validateStateUpgrades` |
| Future-proofing | The enum is extensible ([MigrationSourceEnum](../utils/data/enums.go) at enums.go:340). Today only `from_sdkv2` is recognised; adding a second source (e.g. `from_terraform_provider_ciscomso`) is a new iota constant + `String()` / `UnmarshalText()` case, no consumer changes needed. | [enums.go:335-369](../utils/data/enums.go) |

**Orthogonal to `state_upgrades`** by design: a framework-native resource can still grow `state_upgrades` entries for v0→v1 framework-internal schema bumps without ever setting `migration_source`. The two axes:

| | `migration_source` set | `migration_source` unset |
|---|---|---|
| `state_upgrades` present | SDKv2 resource (possibly with later framework-internal bumps) | Framework-native resource that has gone through one or more schema bumps |
| `state_upgrades` absent | **Invalid** — caught by the validator | Framework-native resource, no upgrades ever |

### 2.2 Test-dependency auto-resolution and override criteria

The legacy `parents:` / `targets:` blocks on each `properties/<class>.yaml` were the only source of test-dependency data in v2.19.0 — the generator did not derive them from meta `containedBy` or `relationInfo.toMo`. The new pipeline reverses that: it derives the common case from meta and treats the YAML override as **additive on top of auto-resolution** (or full replacement when the class sets `test_config.replace_auto_resolved: true`). The migration script must therefore filter the legacy entries against what auto-resolution will produce; entries that auto-resolve correctly are **dropped**, entries that do not are emitted to `test_config.dependencies[]`.

**What auto-resolves today (no YAML needed)**

| Slot | Auto-resolver | Inputs |
|---|---|---|
| `Class.Parents` (the parent class set) | [setParents](../utils/data/class.go) (class.go:524) | meta `containedBy` ∪ `class.include_parents`, minus `class.exclude_parents` ∪ `global.exclude_parents`. |
| Top-level `TestDependency{Role: Parent}` (first **2** parent classes) | [resolveParentDependencies](../utils/data/class.go) (class.go:1233) | First parent → `aci_<resource>.test.id` + `aci_<resource>.test_2.id` (two instances for ForceNew testing). Second parent → `aci_<resource>.test.id`. Additional parents are skipped — provide them explicitly if needed. |
| Top-level `TestDependency{Role: Target}` (single-target only) | [resolveTargetDependencies](../utils/data/class.go) (class.go:1268) | `Relation.ToClasses[0]` → `aci_<resource>.test.id` + `aci_<resource>.test_2.id`. **Multi-target relations (`len(Relation.ToClasses) > 1`) raise a diagnostic** unless an explicit `role: target` dependency is declared. |
| Recursive `TestDependency.Dependencies` (the parent's own parents) | [buildDependency](../utils/data/class.go) (class.go:1300) | For each auto-built dep, recurses through `depClass.Parents` and emits `aci_<resource>.test.id` for each. The `ReferenceTypeEnum` passed here (`ResourceReference` / `DataSourceReference` / `StaticReference`) is honored downstream by `setParentDn` / `setTargetDn` / the placeholder resolvers so static-DN deps render as `StringValue` and Terraform references render as `ReferenceValue`. |
| `parent_dn` attribute test values | [setParentDn](../utils/data/class.go) (class.go:1420) | Wires `Create` / `Update` / `Default` from the first Parent dep; `ForceNew` from the second. Skipped if the property has an explicit `test_config`. |
| `tDn` / `<target_resource_name>_name` attribute test values | [setTargetDn](../utils/data/class.go) (class.go:1488) / [setTargetNameProperty](../utils/data/class.go) (class.go:1555) | Explicit relations expose `tDn` (full DN); Named relations expose `tn<TargetCap>Name` as the meta property, with `AttributeName` renamed to `<target_resource_name>_name` (e.g. `contract_name`) when the target class has a resolvable `resource_name`. Wired from the auto-resolved Target dep. |

**Per-entry decision tree (run during migration)**

For each legacy `parents[]` entry:

1. Resolve `<class_name>`'s meta `containedBy`. If the entry's `class_name` is **not** in the union of meta `containedBy` ∪ existing `include_parents`, append it to `class.include_parents`. Auto-resolution then covers the dep.
2. If `parent_dn` is `aci_<R>.test.id` where `<R>` matches the auto-derived resource name for `class_name`, the test dependency itself is **redundant — drop the entry**. The recursive `parent_dependency` chain is also redundant if it matches the parent's meta `containedBy[0]`.
3. If `parent_dn` is a static literal (e.g. `uni/infra`), emit `test_config.dependencies[]` with `class_name`, `reference: <literal>`, `reference_type: static`, `role: parent`. **Future optimisation — keep the dep for now.** The static dep is only needed because [setParentDn](../utils/data/class.go) (class.go:1420) returns early when no Parent dep exists, leaving `parent_dn.TestValues` nil — the dep is what populates the Create/Update/Default `TestValueEntry` with `ConfigInclude: true` so the generated HCL carries an explicit `parent_dn = "uni/infra"`. When the property already carries `default_values: { "uni/infra": "" }`, the optional-property branch in [generateDefault](../utils/data/property.go) (property.go:705) already produces a valid test (config-omit + assert-default), so the static dep is redundant for correctness — the class doesn't need to be created either, since `uni/infra` is a system DN. Migration keeps the dep for parity with the legacy YAML; a future change to `setParentDn` that falls back to `default_values` when no Parent dep is present would let `migrate_class_definitions.go` drop the static dep for any class whose `default_values.parent_dn` already covers the value. Tracked here so the optimisation surfaces when the test renderer is touched.
4. If the legacy `parent_dependency` chain differs from meta (e.g. selects a non-default ancestor for the test parent), emit the explicit dep with a recursive `dependencies[]` entry.
5. If the class declares 3+ parents (only `fvRsSecInherited` does today, with 3), emit the 3rd as an explicit dep. Retires once the polymorphic-same-type auto-detector (§8.6) lands and bypasses the 2-parent cap for matching classes.
6. Drop `target_classes` and `class_in_parent` from the *output* YAML unconditionally (the new schema has no equivalent slot). The script still **reads** `target_classes` to detect the polymorphic-same-type pattern (§8.6): when every `target_classes` entry across a class's `parents:` list matches its own `class_name` 1:1 (only `fvRsSecInherited` does today), union all distinct target-class values into `relation_info.to_classes` so the resolved `Class.Relation.ToClasses ⊆ Class.Parents`. The polymorphic auto-detector then drives multi-scenario test rendering with no further YAML hint. `class_in_parent` is fully covered by the recursive dependency shape; warn if either key carries data that is not reproducible by the new shape.

For each legacy `targets[]` entry:

1. If `len(Relation.ToClasses) > 1` (multi-target — 13 files), the explicit `role: target` dep is **mandatory** today. This requirement retires for polymorphic-same-type relations once the auto-detector and scenario-iterating renderer land (§8.6); until then, `fvRsSecInherited` migrates with explicit deps like any other multi-target class. Emit one entry per target class with `properties` → `config_overrides` and `target_dn` → `reference` (`reference_type: static` for literals, `resource` for `aci_<R>.test.id`).
2. If single-target and the legacy entry has no `properties:` / `static: true` / `overwrite_parent_dn_key` / non-default `parent_dependency_dn_ref`, the entry is **redundant — drop it**. The new pipeline emits the same dep with a resource reference.
3. If `properties:` is present, emit a `config_overrides` map with the same key/value pairs.
4. If `static: true` (17 entries), emit `reference_type: static` with the literal `target_dn`. Confirm the literal still resolves on a live APIC; if the underlying class now has a resource, prefer `reference_type: resource` (`aci_<R>.test.id`).
5. If `overwrite_parent_dn_key` is set (31 entries), verify the target class's pipeline-derived `parent_dn` matches. When it doesn't, emit `config_overrides: { <legacy-key>: <reference> }` to preserve the wire shape.
6. Drop `relation_resource_name` (auto-derived from `Relation.ToClasses[0]`), `target_dn_overwrite_docs` (§8.4 POSTPONE — kept until the example renderer lands), `shared_classes` and `target_dn_ref` (no consumer in new pipeline — assert empty during migration; warn otherwise).

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
| `legacy_definitions/properties/resource_name_overwrite.yaml` (entire file) | 1 | The 7 relation → resource-name entries are now expressed as per-class `resource_name` on each relation class definition. |
| `exclude` | 0 v2.19.0 files | Class-level boolean opt-out flag read by [SetClassExclude](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2019) (`m.Exclude` then drives registry filtering at v2.19.0:1411 and 1420). 0 v2.19.0 files declare it. The opt-out direction is covered by the new `ClassDefinition.Artifacts []ArtifactEnum` (§4.1): `artifacts: []` excludes the class from `provider.Resources()` / `provider.DataSources()`. Migration is a no-op today; if a future YAML adds `exclude: true`, translate to `artifacts: []`. |
| `multi_line` | 0 v2.19.0 files | Property-level list of property names whose test values should be wrapped in HCL heredoc (`<<EOT … EOT`) syntax. Read by `LookupTestValue` at [v2.19.0:745](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L745) via `processMultiLine`. 0 v2.19.0 files declare it. The new test renderer is expected to choose HCL formatting based on the value itself (presence of newlines) and per-property `value_type` rather than a YAML allowlist; no migration needed. |
| `required_by_custom_type_in_test` | 0 v2.19.0 files | Property-level list of property names to include in the custom-type test. Read by `IncludeInCustomTypeTest` at [v2.19.0:3681](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3681). 0 v2.19.0 files declare it. The new pipeline expresses custom-type test inclusion via the per-property `test_config.custom_type` bucket (already covered by §2 `test_values`); no migration needed. |
| `exclude_targets` | 4 (`fvAEPg`, `fvESg`, `l3extInstP`, `fvRsSecInherited`) | Legacy `getExcludeTargets` ([v2.19.0:3259](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3259)) subtracts entries from `SetModelTestDependencies`'s automatic union of child-class `targets:` lists ([v2.19.0:3245](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3245)). All 4 entries are restatements of the polymorphic-same-type rule (§8.6): each EPG class excludes its peer EPG-like targets so that `fvRsSecInherited` rendered as a child of (say) `fvAEPg` only points its `tDn` at another `fvAEPg`; `fvRsSecInherited.exclude_targets: [l3extInstP]` removes the 2 l3extInstP entries from its own 6-entry `targets:` list as a safety against cross-type selection in legacy logic. Under the new pipeline two factors collapse the override to redundant: (a) [collectChildDrivenDependencies](../utils/data/class.go) (line 1963) is **value-driven** — it never auto-unions child `targets:` lists, so the 3 EPG files' subtraction targets don't exist to begin with; (b) the polymorphic-same-type auto-detector (§8.6) derives the parent-class → target-class match directly from `Parents ∩ ToClasses` plus the rendering site, so neither `exclude_targets` nor `parents[].target_classes` are needed to encode same-type filtering. Migration: drop the key from all 4 files. The migration script verifies that each file's exclusion list matches the auto-detector's prediction (catches editorial divergence from the rule) and warns otherwise; see §9.2 step 4. |
| `custom_type` (nested under `test_values`) | 11 (10 IpAddress + `vmmDomP.arp_learning`) | Legacy `custom_type:` bucket emits an extra `testConfig<Class>CustomType` step that exercises the custom-typed property against an alternate format — for IpAddress, an IPv6 literal next to the IPv4 `all:` value. The new pipeline has no separate "custom-type test" concept: the test scenario framework that would re-introduce multi-format coverage is unresolved (paired with the test renderer rewrite). Migration: for the 10 IpAddress files (`fvBD`, `fvEpIpTag`, `fvFBRMember`, `fvFBRoute`, `fvIpAttr`, `fvTrackMember`, `mgmtSubnet`, `mplsNodeSidP`, `netflowExporterPol`, `pimRouteMapEntry`), fold the author-chosen IPv6 value into the same property's `create:` (or `default:`) entry — author control over the test format is preserved; the auto-derived IPv4 value for that property is replaced. The 1 `vmmDomP.arp_learning: "disabled"` entry stays untouched and rides along with §8.1 (named-SemanticEquality variant framework). |

---

## 4. ADD — schema additions

Keys decided as ADD: the semantics cannot be expressed by an existing field, derived in Go, or relocated to a constant. Each row maps the old key to a new field that needs to be added to the loader structs in [definitions.go](../utils/data/definitions.go).

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

`ArtifactEnum` values are `resource` and `datasource`. Empty/unset list means "auto-derive": expose as both when `IdentifiedBy` is non-empty, expose as nothing when it is empty (the legacy `provider.go.tmpl` default). A non-empty list overrides the auto-derivation and acts as both axes at once:

- Explicit opt-in for classes with empty `IdentifiedBy`. The four migration classes all carry `identifiedBy: []` in meta but ship as public resources today (`aci_epg_useg_block_statement`, `aci_vmm_uplink_container`, `aci_any`, `aci_associated_site`); each becomes `artifacts: [resource, datasource]`. Without it they would silently disappear from the regenerated provider — a breaking change.
- Artifact selector when the auto-derived "both" is wrong. `topSystem` is datasource-only (`aci_system`); `artifacts: [datasource]` suppresses the resource. There is no current resource-only class, but the shape allows `artifacts: [resource]`.

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
| `parent_class` | User-facing parent the variant exposes (e.g., `fvTenant`). Listed in meta `containedBy`, but `containedBy` does not distinguish a system-scoped placement from a wrapped tenant-scoped placement. |
| `rn_prepend` | Intermediate RN segment that selects this variant. The generated resource matches it against the user’s `parent_dn` to route the API call. |
| `wrapper_class` | Implicit container the request nests the resource inside. Empty for variants that POST against a real, user-addressable parent. |
| `test_platform` | Which platform profile (`apic` / `cloud`) exercises this variant in tests. Lets `testvars.yaml.tmpl` gate variant-specific test cases without splitting the test file. |

Not REUSE on `include_parents`: that field is a flat `[]string`; flattening loses every `rn_prepend` / `wrapper_class` / `test_platform` association.

Not DERIVE: the wrapper-container relationship (`fvTenant` → implicit `cloudCertStore` → `pkiKeyRing`) is not modeled in meta JSON — `containedBy` only lists `cloud:CertStore` and `pki:Ep` as direct parents and gives no signal that one of them is auto-created or that `fvTenant` is the user-facing entry point. It is a generator-side convention captured nowhere else.

Deprecation path: only viable if the meta file ever exposes per-parent runtime hints (auto-create flags, user-facing-parent indicators). Today’s meta has neither.

#### `ClassDocumentationDefinition.example_parent_classes []string` — was `example_classes`

Files: `fvRsCons`, `fvRsProv`, `fvRsProtBy`, `fvRsConsIf`, `fvRsIntraEpg`, `fvRsSecInherited`, `tagAnnotation`, `tagTag` (8). All are relation or tag classes whose meta `containedBy` set is huge — every `fv:AEPg`, `fv:ESg`, `l3ext:InstP`, `l2ext:InstP`, `mgmt:InstP`, etc. Rendering one example block per parent in the resource/datasource docs would produce dozens of nearly-identical HCL snippets per page.

The field drives `DocumentationExamples` in the templates: [resource_example.tf.tmpl](../templates/resource_example.tf.tmpl), [datasource_example.tf.tmpl](../templates/datasource_example.tf.tmpl), [resource_example_all_attributes.tf.tmpl](../templates/resource_example_all_attributes.tf.tmpl), [resource.md.tmpl](../templates/resource.md.tmpl), and [testvars.yaml.tmpl](../templates/testvars.yaml.tmpl) all iterate `DocumentationExamples` to emit one example per chosen parent. The override picks a small, illustrative subset (e.g., `fvAEPg, fvESg` for `fvRsCons`; `fvTenant, fvAEPg` for `tagAnnotation`) so the rendered docs stay readable.

Not REUSE: there is no existing list field on `ClassDocumentationDefinition` that names a representative subset of `containedBy`. `dn_formats` (the existing override) is a documentation-only flat list of full DN strings, not a parent-class projection — the renderer needs the `ClassName` so it can resolve the corresponding generated resource name (`getResourceName`) for the example HCL.

Not DERIVE: choosing “which 2–3 parents are illustrative” requires editorial judgement (e.g., for `fvRsCons` the meaningful examples are `fvAEPg` and `fvESg`, not `mgmtInstP`). No meta signal ranks parents by representativeness, and a heuristic like “first N from `containedBy`” would render a different and noisier subset than the current docs.

Deprecation path: only when the docs renderer learns to prune `containedBy` to a representative set on its own — e.g., “one per top-level parent family” — with editorial overrides expressed elsewhere.

#### `ClassTestConfigDefinition.ignore_tests []IgnoreTestEnum` — was `exclude_from_testing`

Files: `vmmRsDomMcastAddrNs`, `vmmRsPrefEnhancedLagPol` (2). Both relation classes have empty `IdentifiedBy` and are nested-child-only on the parent resource (`relation_to_multicast_pool` and `relation_to_lacp_enhanced_lag_policy` on `aci_vmm_domain` — see [resource_aci_vmm_domain.go](../../internal/provider/resource_aci_vmm_domain.go)). They never produce a standalone resource and so have no standalone test of their own.

`IgnoreTestEnum` values are `child`, `resource`, `datasource`. Empty/unset list = no skips. Each value suppresses a distinct test surface:

| Enum value | Scope | What it suppresses |
|---|---|---|
| `child` | Parent’s testvars iteration | One child block in every parent’s generated `testvars.yaml`. The class’s own resource / datasource tests still emit. |
| `resource` | This class’s own resource | The `resource_aci_<x>_test.go` file. Runtime resource artifact (`resource_aci_<x>.go`), schema, docs, and examples are still generated. |
| `datasource` | This class’s own datasource | The `data_source_aci_<x>_test.go` file. Runtime datasource, schema, docs, and examples are still generated. |

The two migration classes both carry `ignore_tests: [child]`. The `resource` and `datasource` values have no legacy YAML driver — they exist in the enum so the schema can express the orthogonal case where the artifact is wanted but its generated test is genuinely unrunnable (APIC-side race, hardware/license requirement, non-deterministic ordering). Today that case is handled out-of-band: hand-deleting the generated test file or carrying a `Skip(…)` shim after every regeneration. The flag lets the YAML record the decision so regeneration honours it.

Why a list and not three booleans: the three scopes are the same axis (test-emission skip at three different render sites) and a class can legitimately combine them, e.g. `[child, resource]` for a relation class whose generated resource test is also flaky. Mirrors the `Artifacts []ArtifactEnum` shape from earlier in this section and keeps the schema growable — a future `import_step` or `replace_step` value drops in without a new field.

Distinct from `ignore_import_state_verify`: that flag suppresses one assertion inside a test that still runs; `ignore_tests` suppresses whole test files (or one child block of a peer’s test file).

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

The new pipeline’s [setDescription](../utils/data/property_documentation.go) currently runs only the per-class override → meta comment → meta label chain. The global layer was dropped silently during the migration scaffold, so today regenerated docs would diverge from legacy in 26 × N places without any per-class change being made. Naming alignment with existing siblings on the same struct: `AttributeNameOverrides` (`attribute_name_overrides`), `DocumentationLabelOverrides` (`documentation_label_overrides`), and now `PropertyDocumentationOverrides` (`property_documentation_overrides`). The YAML key is renamed from `documentation` during migration so the three globals share the `*_overrides` suffix and the override-not-default semantic is visible in the key itself.

Renderer change: insert a global lookup between (1) per-class override and (2) meta comment in [setDescription](../utils/data/property_documentation.go). When the global entry contains `%s`, interpolate with the class’s humanised resource name (the same `GetResourceNameAsDescription` substitution the legacy generator does).

Not DERIVE: the override values are editorial English. There’s no derivation rule that turns 171 distinct meta sentences for `annotation` into one normalised line — that’s a writing decision, not a computation.

Not REUSE: no existing field carries per-meta-property text overrides. Per-class `PropertyDefinition.Documentation.Description` is the natural per-class escape hatch *above* this global, not a replacement (writing 26 overrides into every one of the ~180 class files would multiply the migration footprint by ~180× for no editorial benefit).

Not POSTPONE: the data shape is an unambiguous sibling of two existing globals (`AttributeNameOverrides`, `DocumentationLabelOverrides`); the renderer change is one branch in `setDescription`; and skipping it produces a visible docs regression on regeneration, not a no-op.

Deprecation path: if the meta upstream is ever rewritten to provide consistent property comments, the global overrides can be pruned entry-by-entry — each removal is a single line plus a docs diff review. If every entry ends up removable, the field drops with the same shape it landed in.

---

## 5. REUSE — remap to existing fields

Keys decided as REUSE: semantics already match a field that exists in the new schema. No code changes; the migration script rewrites the value into the existing target.

| Old key | Files | Target existing field | Migration note |
|---|---:|---|---|
| `static_custom_type` (`ip_address` entries) | 11 files / 14 entries | `PropertyDefinition.value_type: ip_address` | The new `ValueTypeEnum.IpAddress` member is auto-derived from meta `validateAsIPv4OrIPv6: true` in [setValueType](../utils/data/property.go) (precedence rule 3), so every `static_custom_type: ip_address` entry is a redundant restatement of the meta. Migration script: drop the override entirely after asserting the meta carries the flag (warn on any mismatch — catches a class where the override exists *because* the meta lacks it). The remaining non-`ip_address` entry (`vmmDomP.arpLearning: vmm_arp_learning`) is POSTPONE — see §8. |
| `max_one_class_allowed` | 2 (`fvFBRoute`, `infraRsHPathAtt`) | `ClassDefinition.is_single_nested_when_defined_as_child` | Every template reference to `.MaxOneClassAllowed` already collapses to `or (not .IdentifiedBy) .MaxOneClassAllowed`, which is exactly how `IsSingleNestedWhenDefinedAsChild` is computed in [class.go](../utils/data/class.go) (`ClassDefinition.IsSingleNestedWhenDefinedAsChild \|\| len(IdentifiedBy) == 0`). The `provider.go.tmpl` registry filter (`hasPrefix .RnFormat "rs"`) is also covered: classes with empty `IdentifiedBy` are already excluded by the outer `and .IdentifiedBy` guard. Set `is_single_nested_when_defined_as_child: true` on the two classes and rename the template lookups. |
| `parent_example_dn` | 1 (`vmmDomP`) | `ClassTestConfigDefinition.dependencies[].reference` (with `reference_type: static`) | Same datum already lives on the static parent dependency for `vmmDomP` (`uni/vmmp-VMware`). Doc renderer should consume that dependency’s `reference` string. Drop the standalone key. |
| `remove_from_contains` | 1 (`fvRsPathAtt`, value `l2PortSecurityPol`) | `ClassDefinition.exclude_children` | Legacy `SetClassContains` ([v2.19.0:2157](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L2157)) builds `m.Contains` from meta `contains` minus the `remove_from_contains` list; `m.Contains` then feeds `DocumentationChildren` ([v2.19.0:3582](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3582)) — the "Children" link list rendered in the resource markdown. The new pipeline’s [setChildren](../utils/data/class_documentation.go) (docs side) builds `Documentation.Children` from meta `rnMap` but currently only excludes children already embedded as nested attributes — it does **not** consult `ExcludeChildren`. Migration sets `exclude_children: [l2PortSecurityPol]` on `fvRsPathAtt` **and** adds one line to the docs-side `setChildren` to also skip entries in `class.ClassDefinition.ExcludeChildren`. The field already covers nested-children generation; this one-line extension makes it cover docs rendering too — unifying the two legacy keys (`exclude_children` + `remove_from_contains`) into the single new field. |

---

## 6. DERIVE — compute in Go (drop YAML key)

Keys decided as DERIVE: value is a function of data the generator already resolves. The migration script drops the key; `class.go` / `property.go` compute the same answer at codegen time.

| Old key | Files | Derivation rule |
|---|---:|---|
| `data_source_has_no_name_identifier` | 1 (`vzAny`) | `len(class.IdentifiedBy) == 0`. The flag is true for `vzAny` only because it has no naming property — a condition the new generator already exposes. The datasource template branches on `IdentifiedBy` directly. |
| `datasource_required` (nested under `test_values`) | 38 | The datasource test’s lookup config = the renamed `IdentifiedBy` set (snake-case via global `attribute_name_overrides` + per-property `attribute_name`); the values = the `Default` bucket (already populated from `resource_required`). Per-file audit (38 files): 36 carry only renamed-`IdentifiedBy` keys with values identical to `resource_required`; the 2 outliers are `fvCrtrn` (empty `IdentifiedBy`, covered by `artifacts: [resource, datasource]` in §4.1) and `vmmDomP` (`parent_dn` against the static parent, covered by `parent_example_dn` in §5 / `static_parent` below). Migration script: drop the nested block, with a per-file safety check warning when a key is not in the renamed `IdentifiedBy` or its value diverges from `resource_required`. |
| `datasource_non_existing` (nested under `test_values`) | 42 | Sibling of `datasource_required`: same renamed-`IdentifiedBy` key set, but values are a **type-aware non-matching transform** of `resource_required` so the generated datasource test verifies the "no result" branch. Two derivation branches across the 42 files: (a) string-typed naming properties append `_non_existing` (e.g. `criterion` → `criterion_non_existing`, `"131"` → `131_non_existing`); (b) IpAddress-typed properties pick a non-matching IP that still passes the validator (e.g. `10.0.0.2` → `10.0.1.2`, `2.2.2.3` → `2.2.2.4` — typically the next octet). Migration script: drop the nested block, with a per-file safety check warning when a key is not in the renamed `IdentifiedBy` set or the value is not a non-matching transform of `resource_required` for the property's `ValueType`. Exact value equality with the auto-derived candidate is **not** required — the assertion is "non-matching and validator-valid", since the IP increment choice is editorial. The new datasource test renderer applies the same transform at codegen time. |
| `static_parent` | 1 (`vmmDomP`) | True iff the resolved `class.test_config.dependencies[]` for the parent role contains a `reference_type: static` entry. Templates branch on the dependency shape instead of a separate flag. |
| `resource_identifier` | 1 (`fvTenant`, value `tn`) | Equal to the RN-prefix segment of `meta.fvTenant.rnFormat` (`tn-{name}` → `tn`). Legacy `GetOverwriteResourceIdentifier` ([v2.19.0:3133](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3133)) is consulted by `GetMultiParentFormats` ([v2.19.0:3096](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3096)) as a fallback when the hardcoded `resourceIdentifier` table miss. The new `ParentDnVariants` setter (§9.1.1) derives every variant’s `ParentDn` directly from meta `dnFormats` + `rnFormat`, so the YAML override and the hardcoded table both retire together. Drop the key; the setter computes `tn-` from meta. |

---

## 7. CONST — relocate to constants.go

Keys decided as CONST: rendering-pipeline tuning, not per-class data. Move into [constants.go](../utils/data/constants.go) and drop from YAML.

| Old key | Source | Constant |
|---|---|---|
| `docs_examples_amount` | `classes/global.yaml` | New `constMaxExamplesToDisplay = 2` (alongside the existing `constMaxDnFormatsToDisplay`). |
| `docs_parent_dn_amount` | `classes/global.yaml` | Already represented by `constMaxParentDnsToDisplay = 20`. Drop the YAML key; the runtime override path is unused. |

---

## 8. POSTPONE — revisit during template/test rework

Keys decided as POSTPONE: the right answer depends on a downstream rewrite (test renderer, docs renderer, custom-type framework) that hasn’t happened yet. The migration script preserves the entries unchanged with a TODO comment; the ADD-vs-drop decision is taken when the dependent work lands.

| Old key | Files | Revisit when |
|---|---:|---|
| `class_version_tests` | 1 (`commPol`) | The test templates are rewritten. If the regenerated `commPol` test still needs an independent version filter, ADD `ClassTestConfigDefinition.supported_versions string` (parsed via the existing `Versions` helper). Otherwise drop the key. |
| `static_custom_type` (named-type entries) | 1 (`vmmDomP.arpLearning: vmm_arp_learning`) | A generic mechanism for custom semantic-equality variants is decided — see §8.1. Until then, keep the override and the hand-written [custom_types/arpLearning.go](../../internal/custom_types/arpLearning.go) untouched. |
| `datasource_required` (top-level list) | 1 (`topSystem`) | The render layer learns to honour `ClassDefinition.artifacts` (mechanism decided in §4.2) — see §8.2. Until then, `properties/topSystem.yaml` stays as aspirational placeholder data with no consumer; the live datasource is the hand-written [internal/provider/data_source_aci_system.go](../../internal/provider/data_source_aci_system.go). |
| `ignore_custom_type_docs` | 6 (`fvAEPg`, `fvRsCons`, `fvRsConsIf`, `fvRsProv`, `mgmtInstP`, `mgmtRsOoBCons`) | The renderer’s rule for `SemanticEquality` valid-values + range docs is decided — see §8.3. Until then, keep the override; once a uniform rule lands the entries collapse to DERIVE (or, less likely, REUSE under a renamed flag). |
| `example_value_overwrite` | 7 (`infraRsAccBndlSubgrp`, `infraRsVipAddrNs`, `infraRsVlanNs`, `vmmRsDomMcastAddrNs`, `vmmRsPrefEnhancedLagPol`, `vnsLDevIf`, `fvTrackMember`) | The example/test renderer is built and the relation-target derivation rule is decided — see §8.4. The 5 `target_dn` cases collapse to DERIVE; the 2 editorial cases need a narrower field shape than the original proposal. |
| `custom_test_dependency_name` | 1 (`vnsLDevIf`, value `ImportedVnsLDevVipWithFvTenant`) | The test renderer is built and proves it can emit a separate-tenant target subtree via the existing `test_config.dependencies` mechanism — see §8.5. Until then, the hand-written companion HCL in [internal/provider/test_constants.go](../../internal/provider/test_constants.go) stays untouched. |

### 8.1 `vmmDomP.arpLearning` context

`vmm_arp_learning` is the only non-`ip_address` value across all 11 `static_custom_type` files (14 entries). It is implemented today as a hand-written [custom_types/arpLearning.go](../../internal/custom_types/arpLearning.go) that satisfies `basetypes.StringValuableWithSemanticEquals` with a side-table (`"" ≡ "disabled"`), making it a *variant* of the existing `SemanticEquality` value type rather than an unrelated kind of custom type.

The framework’s auto-derived `SemanticEquality` (triggered in [setValueType](../utils/data/property.go) precedence rule 4 when a property has both `ValidValues` and `Validators`) handles the uniform case where wire and human forms map 1:1 via `validValues`. `vmmDomP.arpLearning` falls outside that contract because its equality reads an additional side-table that the meta does not express — the wire value `"disabled"` aliases to the localName `defaultValue` while `0x0`/`0x1` map to `disabled`/`enabled`, and the empty string normalises to `disabled`. The custom type also overrides the meta `uitype: bitmask` (which would auto-derive to `Set`) to render as a single string, captured today through `type_overwrites: arpLearning: string` alongside the `static_custom_type` entry.

There is no precedent yet for how the framework should declare *named* semantic-equality variants. Possible shapes — all open — include a `value_type: semantic_equality.vmm_arp_learning` enum constant plus a registry of value-map providers; a per-property `equality_map` definition that drives a generated custom type; or pulling the side-tables out of `customtypes/` into a YAML registry. The right design needs its own walk-through and is out of scope for the migration script. Until that lands, the override stays where it is and `arpLearning.go` keeps doing the work by hand.

Deprecation path: once the generic mechanism is in place, this entry collapses into `value_type` (REUSE) and the bespoke `customtypes/arpLearning.go` retires — along with any other future deviations that can’t be derived from meta alone.

The sibling `test_values.arp_learning.custom_type: "disabled"` entry in `properties/vmmDomP.yaml` rides along with this POSTPONE: it is the test-side value for the same custom type, and the migration script leaves it in place until the framework decision lands (the other 10 nested `custom_type:` entries — all IpAddress — are §3 Obsolete and fold into standard buckets; see §3).

### 8.2 `topSystem` / datasource-only artifact context

`aci_system` is the only public datasource without a matching resource. The legacy generator handled it out-of-band: there is no `gen/meta/topSystem.json`, so the class never enters [getClassModels](../generator.go) and never reaches the render gate. The hand-written [internal/provider/data_source_aci_system.go](../../internal/provider/data_source_aci_system.go) survives regeneration because of an explicit whitelist on `cleanDirectory(providerPath, []string{… "data_source_aci_system.go", "data_source_aci_system_test.go" …})`, paired with `cleanDirectory(datasourcesDocsPath, []string{"system.md"})` and `cleanDirectory(datasourcesExamplesPath, []string{"aci_system"})`. The `Required: true` on `system_id` and `pod_id` is hand-coded in that file.

[gen/scripts/legacy_definitions/properties/topSystem.yaml](legacy_definitions/properties/topSystem.yaml) was authored alongside the manual plugin-framework migration of `aci_system` as a placeholder for the day the generator can produce this datasource. Today nothing reads it: the top-level `datasource_required:` list is not in the closed set of keys the legacy generator recognises (the matching template lookup `$.datasource_required` reads the unrelated `test_values.datasource_required` *bucket*, not the top-level list), and the new pipeline has no consumer either.

The mechanism for the new pipeline is decided: `ClassDefinition.artifacts []ArtifactEnum` (see §4.2) covers both registration opt-in and artifact selection. `topSystem` would carry `artifacts: [datasource]` once `gen/meta/topSystem.json` lands and the class enters `loadClasses` like any other. The `Required: true` migration on `system_id`/`pod_id` is a no-op once the resource side is suppressed: the existing per-property `restriction: required` covers the same intent without a new enum axis on `PropertyDefinition`.

What remains is the render-layer build-out:

- The renderer must read `Class.Artifacts` and emit only the listed artifacts (schema, examples, docs, tests).
- The legacy `cleanDirectory(...)` whitelist that preserves hand-written `data_source_aci_system.go` / `system.md` / `examples/aci_system` needs an equivalent in the new pipeline so the hand-written file survives regeneration until the generated replacement is ready.
- Once `topSystem.json` exists and the renderer honours `artifacts: [datasource]`, walk through the `topSystem.yaml` content (`documentation` / `overwrites` / `read_only_properties` fold into the standard schema, `datasource_required` collapses into per-property `restriction: required`), drop the manual file, and retire the whitelist entries.

Deprecation path: this entry retires when the renderer honours `artifacts` and `topSystem.yaml` is fully expressible in the standard schema.

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

`hasCustomTypeDocs` is a registered template helper backed by [HasCustomTypeDocs](../generator.go) (v2.19.0:3768). It returns `true` by default and `false` when the property is listed in the class’s top-level `ignore_custom_type_docs:` list — flipping the docs from the upper branch (with the `Or a value in the range of X to Y` line) to the lower branch (just the localname list).

Why six files carry it for the QoS `Prio` property: the meta validators say `min: 0, max: 9` but the wire `validValues` are non-contiguous (`0,1,2,3,7,8,9` mapped to `level1–level6, unspecified` — no localnames for wire `4/5/6`). The range line is technically truthful (APIC rejects values outside `0..9`) but misleading because users cannot type `"4"`, `"5"`, or `"6"` even though those integers fall inside the validator range.

The override is applied inconsistently. 15 classes ship the same QoS `prio` property with identical meta; only 6 suppress the range line. The other 9 (`fvAp`, `fvESg`, `l3extInstP`, `vzOOBBrCP`, `vzRsAnyToCons`, `vzRsAnyToProv`, `vzRsAnyToConsIf`, `qosDscpClass`, `qosDot1PClass`) keep it. Side-by-side proof in the current docs: [application_profile.md](../../docs/resources/application_profile.md) renders `priority` as `Valid Values: "level1", …, "unspecified". Or a value in the range of 0 to 9.`; [application_epg.md](../../docs/resources/application_epg.md) renders the same property without the range line. There is no editorial reason for the divergence — the 6 fixes correct the docs for some classes; the 9 unfixed classes still mislead.

The new pipeline already auto-derives `ValueType == SemanticEquality` for every `prio` instance via [setValueType](../utils/data/property.go) precedence rule 4 (`len(ValidValues) > 0 && len(Validators) > 0`), so the renderer can detect “this property has both” without per-class hints. The remaining decision is editorial:

- **Uniform suppress** — drop the range line for every `SemanticEquality` property. Cleanest, but loses informative output for `cos` / `dscp` (wire ints equal localnames, range is informative).
- **Localname=wire DERIVE** — keep the range line only when every meta `validValue.value` equals its `localName` (the property is `SemanticEquality` purely because of an `unspecified` sentinel alias). Drops the line for `Prio` (`level1` ≠ wire `3`) and keeps it for `cos` (`0` == wire `0`). Fully derivable from meta, no per-class flag, no inconsistency.
- **Per-property declarative ADD** — keep the override (renamed to something like `PropertyDocumentationDefinition.hide_validator_range bool`). Lets editorial choices stay per-property but perpetuates the 6-of-15 inconsistency unless every QoS-`prio` class is audited.

All three options collapse to either DERIVE or REUSE with a different rendering rule once the renderer is built. The current `ignore_custom_type` ADD proposal was based on the misreading of this flag as “hide a custom-type cross-reference link” (no such link exists in the templates) rather than “suppress one branch of the valid-values rendering for sparse-wire-mapping properties.”

Deprecation path: this entry retires when the docs renderer picks a uniform rule for `SemanticEquality` valid-values + validator-range output. If the rule is meta-derivable (uniform suppress or localname=wire), the override drops entirely (DERIVE). If editorial control is wanted, the override returns under a clearer name covering the actual semantic.

### 8.4 `example_value_overwrite` / relation-target example rendering context

Legacy [generator.go:780 `LookupTestValue`](../generator.go) (v2.19.0) is the generic helper that templates call to render a property value into an example or test step. For relation classes the lookup walks four buckets in order — `testVars["all"]`, `version_mismatch`, `resource_required`, then (for `target_dn` only) `testVars["targets"][i].target_dn_ref` — before falling through to `example_value_overwrite` as a last resort. The legacy author flagged this fallthrough in a self-deprecating comment:

```go
// Referencing is done based on target_dn logic
// This lookup is created as a workaround to reference in an examples on non target_dn attributes
// Redesign of testing / example creation logic should be done to cover this reference use-case
```

A second consumer, `GetTestValueOverwrite` (v2.19.0:878), is hard-coded to `target_dn` and used in three sites of [resource_example_all_attributes.tf.tmpl](../templates/resource_example_all_attributes.tf.tmpl) (v2.19.0:104, 154, 170). Both consumers read the same per-class `example_value_overwrite:` map keyed by the post-rename snake_case attribute name.

The seven entries split into three groups:

| Group | Files | Override | Why it exists |
|---|---|---|---|
| `target_dn` workaround | `infraRsVipAddrNs`, `infraRsAccBndlSubgrp`, `vmmRsDomMcastAddrNs`, `vmmRsPrefEnhancedLagPol` | `target_dn: aci_<resource>.example.id` | The `targets:` list has no `target_dn_ref` (or is commented out for `ignore_tests: [child]` classes), so the fallthrough fires. |
| `target_dn` dead code | `infraRsVlanNs` | `target_dn: aci_vlan_pool.example.id` | `targets[0].target_dn_ref: aci_vlan_pool.test_vlan_pool_1.id` is set, so `LookupTestValue` returns that first and never reaches the override. The entry has no consumer today. |
| Editorial | `vnsLDevIf` (`logical_device: aci_l4_l7_device.example_in_another_tenant.id`), `fvTrackMember` (`scope: aci_bridge_domain.example.id`) | Reference to a non-default resource label or a polymorphic-target choice. | `vnsLDevIf` is the imported device and must reference a device in a *different* tenant; the override deliberately uses `example_in_another_tenant`, not `example`. `fvTrackMember.scope` is polymorphic (Bridge Domain or L3Out per its description) and the override picks BD as the canonical example. |

**Implications for the new pipeline**

Five of the seven (the `target_dn` group plus the dead-code entry) are mechanically derivable: `"aci_" + getResourceName(meta.relationInfo.toMo) + ".example.id"`. The new pipeline already exposes the same datum via [Class.Relation.ToClasses](../utils/data/class.go), and `vmmRsDomMcastAddrNs` / `vmmRsPrefEnhancedLagPol` keep the meta `relationInfo.toMo` even though their `targets:` lists are commented out for the `ignore_tests: [child]` test gap (§4.2). Once the example renderer lands and applies that rule, all five collapse to DERIVE.

The two editorial entries need a narrower field shape than the original `PropertyDocumentationDefinition.example_value string` proposal:

- `vnsLDevIf.logical_device` only diverges from the default DERIVE in the resource label (`example_in_another_tenant` vs `example`). A boolean or label override on the relation property (e.g., `example_resource_label: example_in_another_tenant`) captures it without re-introducing free-form value strings.
- `fvTrackMember.scope` needs a target-class selector among multiple `to_classes`. That datum belongs on the polymorphic relation, not on a generic doc override.

**Why this fits POSTPONE**

- The example/test renderer doesn’t exist yet. Picking the field shape now risks creating dead data: 5 of 7 entries become DERIVE the moment the renderer learns the simple `relationInfo.toMo` rule, and the 2 editorial entries need a different field than the one originally proposed.
- `infraRsVlanNs` is already dead code today; migrating it as-is would carry forward an entry with no consumer.
- The legacy code author explicitly called this a workaround pending redesign of the test/example renderer.

Deprecation path: 5 of 7 retire when the renderer applies the meta-derived `relationInfo.toMo` rule (DERIVE). The 2 editorial entries (`vnsLDevIf.logical_device`, `fvTrackMember.scope`) return under a narrower field shape that’s decided alongside the polymorphic-relation handling (`fvTrackMember.scope_dn` targets BD or L3Out) and the cross-tenant example-label override (`vnsLDevIf`).

### 8.5 `vnsLDevIf` / separate-tenant test-target context

`custom_test_dependency_name: ImportedVnsLDevVipWithFvTenant` is a stopgap that prepends a hand-written Go const to every generated test config for `vnsLDevIf`. The legacy template (`gen/templates/resource_test.go.tmpl`, 10 occurrences) emits `testConfigVnsLDevIfXxx = testConfigImportedVnsLDevVipWithFvTenant + testConfigFvTenantMin + <generated>` across all Min / All / Reset / Children / CustomType / LegacyAttributes variants. The referenced const lives in [internal/provider/test_constants.go](../../internal/provider/test_constants.go) at line 149 and carries HCL for `aci_tenant.test_tenant_imported_device` + `aci_physical_domain.test` + `aci_l4_l7_device.test_imported_device` — a **separate tenant** subtree so that `vnsLDevIf.logical_device` (the renamed naming property, itself a DN reference to a `vnsLDevVip`) can be tested against an `lDevVip` that lives outside the test’s own tenant.

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

### 8.6 Polymorphic same-type relations — auto-detection pattern

A *polymorphic same-type* relation is a relation class whose `Relation.ToClasses` is a subset of its own `Parents`, with cardinality ≥ 2. Only `fvRsSecInherited` matches today: `Parents` = 20 EPG-like classes (meta `containedBy` minus excludes), `ToClasses` = `{fvAEPg, fvESg, l3extInstP}` (after the §2 `relationship_classes` union rule). The relation expresses "this resource references another instance of one of its own valid parent types" — a `fvRsSecInherited` rendered under `fvAEPg` points its `tDn` at another `fvAEPg`, under `fvESg` at another `fvESg`, under `l3extInstP` at another `l3extInstP`.

**Detector**

The pattern is fully derivable from already-resolved fields on `Class`:

```go
func (c *Class) isPolymorphicSameType() bool {
    if !c.Relation.RelationalClass || len(c.Relation.ToClasses) < 2 {
        return false
    }
    parentSet := make(map[string]bool, len(c.Parents))
    for _, p := range c.Parents {
        parentSet[p.String()] = true
    }
    for _, t := range c.Relation.ToClasses {
        if !parentSet[t.String()] {
            return false
        }
    }
    return true
}
```

Rule: `ToClasses ⊆ Parents` with `|ToClasses| ≥ 2`. Not set-equality — `Parents` (from meta `containedBy` minus excludes) is typically a strict superset of the editorially-chosen `ToClasses`.

**Verification against the 4 multi-target relations today**

| Class | `Parents` (resolved) | `ToClasses` (resolved) | `ToClasses ⊆ Parents` | Polymorphic? |
|---|---|---|:---:|:---:|
| `fvRsDomAtt` | `{fvAEPg}` | `{vmmDomP, physDomP, fcDomP, l2extDomP}` | no | no |
| `infraRsDomP` | `{infraAttEntityP}` | `{vmmDomP, physDomP, fcDomP, l2extDomP}` | no | no |
| `netflowRsExporterToEPg` | `{netflowExporterPol, netflowExporterPolDef}` | `{fvAEPg, l3extInstP, l2extInstP}` | no | no |
| `fvRsSecInherited` | 20 EPG-like classes | `{fvAEPg, fvESg, l3extInstP}` | yes | **yes (3 scenarios)** |

Only `fvRsSecInherited` triggers. The three other multi-target relations have disjoint parent and target sets (domain types are never parents of `fvRsDomAtt`; netflow exporters are never EPGs); the detector correctly stays silent.

**What the detector replaces**

1. **`parents[].target_classes`** — the per-parent target-class filter (3 entries on `fvRsSecInherited`) becomes the pairing `Parent[i] ↔ ToClasses[i]` for matched names; no explicit list needed.
2. **`exclude_targets` on the 3 EPG parents** — the same-type filter for child rendering becomes a renderer rule "select target whose class matches the rendering parent class"; no explicit exclusion list needed.
3. **`exclude_targets: [l3extInstP]` on `fvRsSecInherited`** — the safety against cross-type selection in legacy logic becomes a no-op once parent-class → target-class is enforced explicitly.
4. **Multi-target target-dep mandate** (§2.2 targets step 1) — `fvRsSecInherited` becomes auto-resolvable like a single-target relation, emitting N scenarios instead of 1.

**Code changes once the renderer lands**

- `isPolymorphicSameType()` on `Class` (~10 lines, data-only).
- [resolveParentDependencies](../utils/data/class.go) (class.go:1233) — relax the 2-parent cap when polymorphic; emit all N parents.
- [resolveTargetDependencies](../utils/data/class.go) (class.go:1268) — replace the multi-target diagnostic with a polymorphic branch that emits one target dep per `ToClass`, with distinct HCL labels (e.g. `aci_<resource>.test_master.id`) so parent and target of the same class don't collide in HCL.
- [setParentDn](../utils/data/class.go) (class.go:1420) / [setTargetDn](../utils/data/class.go) (class.go:1488) — bail when polymorphic; the renderer iterates scenarios instead of receiving single-step wiring.
- [buildTestChildren](../utils/data/class.go) (class.go:1755) — when the child class is polymorphic-same-type and the calling parent is in its `ToClasses`, override `tDn` in the built instance with the target dep matching the parent's class. This replaces the per-EPG `exclude_targets` lists for child rendering.

Total: ~25–30 lines in `class.go`, no schema changes, no new struct fields.

**Renderer work — POSTPONE (§8.5)**

The renderer must iterate N parent-target pairs when a class is polymorphic; emit one test-step group per scenario; and tolerate distinct HCL labels per dependency (the same capability needed for `vnsLDevIf`'s separate-tenant target subtree). Both converge on the same renderer milestone.

**Migration script impact today**

- Union `parents[].target_classes` into `relation_info.to_classes` (§2 `relationship_classes` row). For `fvRsSecInherited` this expands `[fvAEPg, fvESg]` to `[fvAEPg, fvESg, l3extInstP]`.
- Drop `exclude_targets` from all 4 files (§3). Polymorphic auto-detection covers the filter once it lands; the legacy union mechanism it subtracted from doesn't exist in the new pipeline.
- Drop `parents[].target_classes` from the output YAML (§2.2 step 6). The detector reconstructs the per-parent target choice from `Parents ∩ ToClasses`.
- Until the renderer iterates scenarios, `fvRsSecInherited`'s test stays in its legacy shape via the §2.2 targets step 1 multi-target mandate: the migration emits 3 explicit `role: target` deps so today's renderer can produce a working single-scenario test. When the renderer + auto-detector land, those explicit deps drop and the YAML's `test_config` reduces to nothing.

**Deprecation path**

Retires when the test renderer can iterate polymorphic scenarios and emit distinct HCL labels per `TestDependency.Reference` (the same milestone as §8.5). The 4 `exclude_targets` entries and the 3 `parents[].target_classes` blocks stay dropped (already removed at migration); the `fvRsSecInherited` `test_config.dependencies[]` shrinks to empty once the auto-detector covers it.

**Edge case noted**

Asymmetric polymorphic (parent A → targets `{A, B}`; parent B → targets `{B}` only) is not covered by the subset rule — the detector would emit cross-product scenarios that don't exist on APIC. No such class exists today; if one ever appears, the YAML can opt out via `test_config.replace_auto_resolved: true` and hand-author the dependency chain.

---

## 9. Migration-script work plan

The current [migrate_class_definitions.go](migrate_class_definitions.go) is built out across §§1–7 for both class-level and property-level YAML (the `knownLegacyKeys` allowlist mirrors the §10 audit one-for-one and every entry is `implemented: true`); §8 POSTPONE entries are dropped during migration with a per-file log line. The §9.1 loader struct fields, setters, and constants are all landed. What remains is the downstream renderer / custom-type work the §8 POSTPONE entries are blocked on, the renderer-side polymorphic-same-type detector (§8.6), and a handful of post-pass refinements noted inline (`exclude_targets` editorial verification, optional-property `default_values` fallback in `setParentDn`).

Script inputs are the legacy YAML files under [gen/scripts/legacy_definitions/](legacy_definitions/) plus the one-shot SDKv2 schema dump [gen/scripts/legacy_definitions/schema-git-commit-e21fb3e5.json](legacy_definitions/schema-git-commit-e21fb3e5.json). The JSON is consumed only when translating §2's three SDKv2-migration keys (see §2.1). The legacy YAML files and JSON dump are preserved under `gen/scripts/legacy_definitions/` after migration completes so future adjustments can re-run the script; nothing in the new pipeline reads them.

### 9.1 Code changes required *before* the script runs end-to-end

The script writes YAML; consuming it requires three groups of upstream changes. **Status: all three groups landed before §9.2 work began — the subsections below document the contract each addition honours so future extensions stay anchored to it.**

**1. Loader struct fields** in [definitions.go](../utils/data/definitions.go):

- `ClassDefinition.Artifacts []ArtifactEnum` (+ `ArtifactEnum` constants `resource`, `datasource`).
- `ClassDefinition.ParentDnVariants []ParentDnVariantDefinition` (+ `ParentDnVariantDefinition` struct with `parent_class`, `rn_prepend`, `wrapper_class`, `test_platform` and `PlatformTypeEnum`).
- `ClassDocumentationDefinition.ExampleParentClasses []string`.
- `ClassTestConfigDefinition.IgnoreTests []IgnoreTestEnum` (+ `IgnoreTestEnum` constants `child`, `resource`, `datasource`).
- `ClassTestConfigDefinition.IgnoreImportStateVerify bool`.
- `GlobalMetaDefinition.PropertyDocumentationOverrides map[string]string`.

**2. Resolver/setter changes** in [class.go](../utils/data/class.go), [class_documentation.go](../utils/data/class_documentation.go), and [property_documentation.go](../utils/data/property_documentation.go). Each new loader field needs a resolved counterpart on the runtime struct plus a `set*` method called from the existing setup chain (same pattern as the existing `setIsSingleNestedWhenDefinedAsChild` at [class.go:482](../utils/data/class.go)):

| New loader field | Resolved runtime field | Setter behaviour | Consumer site |
|---|---|---|---|
| `ClassDefinition.Artifacts` | `Class.Artifacts []ArtifactEnum` (same enum on both sides). | Empty YAML → auto-derive: `[resource, datasource]` when `len(IdentifiedBy) > 0`; `[]` otherwise (class is still loaded so it can be referenced as a child/parent/relation target, but not registered as a top-level artifact). Non-empty YAML → copy verbatim, overriding the auto-derivation. | `provider.go.tmpl` registry filter switches from `or (and .IdentifiedBy (not (and .MaxOneClassAllowed (hasPrefix .RnFormat "rs")))) .Include` to a `has .Artifacts "resource"` / `has .Artifacts "datasource"` lookup (one template helper, used in both branches). The four §4.2 opt-in classes carry `[resource, datasource]`; `topSystem` carries `[datasource]`; nested-only relation classes resolve to `[]`. |
| `ClassDefinition.ParentDnVariants` | `Class.ParentDnVariants []*ParentDnVariant` plus a synthesised `Class.DefaultParentDn *ParentDnVariant` for the meta-derived placement. Each `ParentDnVariant` carries the loader fields plus a resolved `ParentClass *ClassName` (so templates can call `getResourceName`) and a resolved `WrapperClass *ClassName` (so the renderer can validate it exists in meta). | See §9.1.1 for the full derivation and emission rules. | Three template sites: `parent_dn` schema attribute (default + validator), Create / Update API routing in `resource_aci_<x>.go.tmpl`, and `testvars.yaml.tmpl` platform gating. Replaces the hand-coded `wrapperClassMap` in [resource_aci_key_ring.go](../../internal/provider/resource_aci_key_ring.go) and the matching block in `resource_aci_certificate_authority.go`. |
| `ClassDocumentationDefinition.ExampleParentClasses` | `ClassDocumentation.ExampleParentClasses []*ClassName` | Resolve each YAML string into a `*ClassName`. When the override is empty, fall back to the existing meta-`containedBy` projection. | The parent-example projection in [class_documentation.go](../utils/data/class_documentation.go) and the `DocumentationExamples` iteration sites in `resource.md.tmpl` / `*_example.tf.tmpl` / `testvars.yaml.tmpl` read the override when populated. |
| `ClassTestConfigDefinition.IgnoreTests` | `Class.TestConfig.IgnoreTests []IgnoreTestEnum` (same enum on both sides). | One-line passthrough from the loader. Empty list = nothing suppressed. | Three render-site gates keyed by enum membership: `child` is consumed by [testvars.yaml.tmpl:270](../templates/testvars.yaml.tmpl) `{{- range $key, $value := .Children}}{{if not (has .TestConfig.IgnoreTests "child")}}…` (replaces the legacy `.ExcludeFromTesting` lookup); `resource` skips emission of `resource_aci_<x>_test.go` in `resource_test.go.tmpl`; `datasource` skips emission of `data_source_aci_<x>_test.go` in `datasource_test.go.tmpl`. The runtime resource / datasource, schema, docs, and examples are unaffected — use `Artifacts` to drop those too. |
| `ClassTestConfigDefinition.IgnoreImportStateVerify` | `Class.TestConfig.IgnoreImportStateVerify bool` | One-line passthrough from the loader. | Resource-test template branch that emits or skips `ImportStateVerify: true`. Distinct from `IgnoreTests: [resource]` — the test still runs; only the one assertion is suppressed. |
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

1. **Derive the default placement from meta.** Walk `meta.dnFormats`; the entry that doesn't start with any YAML `rn_prepend` is the default (for `pkiKeyRing`: `uni/userext/pkiext/keyring-{name}`). Strip the trailing identifying RN segment via `class.RnFormat` to get the default `parent_dn` prefix (`uni/userext/pkiext`). Resolve `meta.containedBy` for that DN to a `*ClassName` (`pki:Ep`) and store as `Class.DefaultParentDn = &ParentDnVariant{ParentClass: pkiEp, ParentDn: "uni/userext/pkiext", RnPrepend: "", WrapperClass: nil, TestPlatform: apic}`.
2. **Resolve each YAML variant.** For each `parent_dn_variants[]` entry:
   - Resolve `parent_class` and `wrapper_class` (when set) into `*ClassName` instances via the existing `loadClasses` cache; if either is unknown emit a generator error (typo catch).
   - Cross-check: `parent_class` must be present in meta `containedBy` for at least one of the meta `dnFormats`; warn otherwise.
   - Cross-check: `rn_prepend` must appear as a literal segment in at least one meta `dnFormat`; warn otherwise.
   - Compute the variant's `ParentDn` from the matching `dnFormat` (everything before `{name}`), e.g. `""` for the tenant-scoped `keyring` placement — the user-facing `parent_dn` is the tenant's DN, supplied at runtime, and `rn_prepend` is the only fixed suffix added between the tenant DN and the RN.
3. **Assemble** `Class.ParentDnVariants` = `[DefaultParentDn] ++ resolved variants`. Single ordered list keeps the renderer trivial.

**Renderer (one template helper, three emission sites)**

A. **`parent_dn` schema default + validator** in `resource_aci_<x>.go.tmpl`. Today hand-coded as `stringdefault.StaticString("uni/userext/pkiext")` at [resource_aci_key_ring.go:191](../../internal/provider/resource_aci_key_ring.go). Generated as:

```go
"parent_dn": schema.StringAttribute{
    Optional: true,
    Computed: true,
    Default:  stringdefault.StaticString({{ .DefaultParentDn.ParentDn | quote }}),
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            regexp.MustCompile(`^(?:{{ range $i, $v := .ParentDnVariants }}{{ if $i }}|{{ end }}{{ $v.ParentDn }}{{ if $v.RnPrepend }}/[^/]+/{{ $v.RnPrepend }}{{ end }}{{ end }})`),
            "parent_dn must reference one of the supported parent placements",
        ),
    },
},
```

B. **Create / Update API routing.** Replaces the hand-coded `wrapperClassMap` literal:

```go
// generated from .ParentDnVariants
parentDnRouting := []struct{ Marker, Wrapper string }{
{{- range .ParentDnVariants }}
    { Marker: {{ if .RnPrepend }}{{ .RnPrepend | quote }}{{ else }}{{ .ParentDn | quote }}{{ end }}, Wrapper: {{ .WrapperClass.PkgName | quote }} },
{{- end }}
}
for _, route := range parentDnRouting {
    if !strings.Contains(data.Id.ValueString(), route.Marker) { continue }
    if route.Wrapper != "" {
        DoRestRequest(ctx, &resp.Diagnostics, r.client,
            fmt.Sprintf("api/mo/%s%s.json", strings.Split(data.Id.ValueString(), route.Marker)[0], route.Marker),
            "POST", jsonPayload)
    } else {
        DoRestRequest(ctx, &resp.Diagnostics, r.client,
            fmt.Sprintf("api/mo/%s.json", data.Id.ValueString()),
            "POST", jsonPayload)
    }
    break
}
```

The order matters: variants with a non-empty `rn_prepend` (the more specific marker, e.g. `certstore`) must be tested before the default (the meta DN prefix, e.g. `uni/userext/pkiext`), otherwise the default would match first for the wrapped placement too. The setter already returns the slice in default-then-variants order, so the template iterates in reverse (`{{ range $i := slice .ParentDnVariants ... }}`) or the setter swaps the order — trivial either way, but the test must cover both placements for both classes to lock the contract.

C. **Test platform gating** in `testvars.yaml.tmpl`. The setter already exposes `TestPlatform` per variant; the template emits the variant's test block only when the current render's platform matches:

```
{{- range .ParentDnVariants }}
{{- if or (eq .TestPlatform $.RenderPlatform) (eq .TestPlatform "") }}
# variant block for {{ .ParentClass.PkgName }} on platform {{ .TestPlatform }}
{{- end }}
{{- end }}
```

For `pkiKeyRing`/`pkiTP` the YAML variant carries `test_platform: cloud`, so its block only appears in cloud-platform testvars; the default (apic) variant always appears.

**Deprecation path**

Everything in [resource_aci_key_ring.go](../../internal/provider/resource_aci_key_ring.go) lines 463–471 and 538–549 (the two `wrapperClassMap` literals plus the hand-coded `Default:` at line 191) is replaced by generated code driven from the YAML. The hand-written `resource_aci_key_ring.go` / `resource_aci_certificate_authority.go` files retire once the generator emits a matching file.

### 9.2 Script work, by section

All ten steps below are implemented in [migrate_class_definitions.go](migrate_class_definitions.go); the descriptions document the contract each step honours so a future re-run or extension stays anchored to the disposition catalog. Refinements still in flight are called out inline.

1. **§1 direct mapping** — extend the loader to also copy `rn_prepend`, `required_as_child`, `resource_name`, and `dn_formats` (under `documentation:`), plus the per-property `documentation` map (121 files).
2. **§2 semantic mapping, class-level** — implement the shape-changing transforms: `resource_notes` (+ `resource_warnings` / `datasource_notes` / `datasource_warnings` sibling slots, no v2.19.0 data), `children` → `include_children`, `contained_by` → `include_parents` (subtract meta `containedBy` first), `class_version` → `supported_versions`, `relationship_classes` (+ drop `multi_relationship_class`) → `relation_info.to_classes`, and `migration_version` + `migration_blocks` + `type_changes` → `state_upgrades` (two-source merge with [schema-git-commit-e21fb3e5.json](legacy_definitions/schema-git-commit-e21fb3e5.json), see §2.1; sets `migration_source: from_sdkv2` and omits `legacy_type` / `legacy_restriction` when they match the current property).
3. **§2 semantic mapping, property-level** — add a second pass that ingests `properties/<class>.yaml` and folds the entries into `NewClassDefinition.Properties` keyed by meta camelCase name. Covers the 12 property-level transforms (~125 files), the per-class `parents` / `targets` aggregation into class-level `test_config.dependencies[]` (apply the §2.2 decision tree — drop entries that auto-resolve from meta `containedBy` / `Relation.ToClasses[0]`; emit only the irreducible overrides; promote unknown parent classes into `class.include_parents` then auto-resolution covers them), and the `test_values` bucket merge. Also detect the polymorphic-same-type pattern (§8.6) — when every `parents[].target_classes` entry matches its own `class_name` 1:1, union the distinct target-class values into the class's `relation_info.to_classes` so future auto-detection has the data it needs.
4. **§3 obsolete** — drop the 4 global keys, `multi_relationship_class`, the entire `legacy_definitions/properties/resource_name_overwrite.yaml` file, and `exclude_targets` (4 files) with a per-file log line. For each `exclude_targets` file, verify that the exclusion list matches the polymorphic-same-type auto-detector's prediction (§8.6) — i.e. each EPG class excludes exactly its peer same-domain classes, and `fvRsSecInherited` excludes exactly the classes in its own `targets:` list that are not in `relation_info.to_classes`. Warn on any divergence (catches editorial choices that diverge from the rule).
5. **§4 ADD** — emit the six migrated fields into the output YAML. Verbatim values; one value remap (legacy `exclude_from_testing: true` → `ignore_tests: [child]`) plus the YAML key renames listed in the §4.1 table (e.g. `documentation` in `properties/global.yaml` → `property_documentation_overrides` in `global.yaml`). The `resource` and `datasource` enum values of `ignore_tests` have no legacy driver and so are not written by the script — they exist in the schema as future opt-ins and the loader accepts them when present.
6. **§5 REUSE** — apply the four remaps: `static_custom_type: ip_address` drop with meta-flag assertion; `max_one_class_allowed` → `is_single_nested_when_defined_as_child`; `parent_example_dn` drop (covered by the static dependency); `remove_from_contains` → `exclude_children` (one-line addition to the docs-side `setChildren` so the field covers docs rendering too).
7. **§6 DERIVE** — drop the five keys with a one-line comment per file. For `datasource_required` nested, emit a warning when a key isn’t in the renamed `IdentifiedBy` or its value diverges from `resource_required` (catches drift). For `datasource_non_existing` nested, emit a warning when a key isn’t in the renamed `IdentifiedBy` or its value is not a non-matching transform of `resource_required` for the property’s `ValueType` (catches drift). For `resource_identifier`, drop without further checks — the value is derivable from meta `rnFormat` by the `ParentDnVariants` setter.
8. **§7 CONST** — drop both keys from `classes/global.yaml`.
9. **§8 POSTPONE** — drop the 6 entries (`class_version_tests`, `static_custom_type` named-type, `datasource_required` top-level, `ignore_custom_type_docs`, `example_value_overwrite`, `custom_test_dependency_name`) with a per-file `POSTPONE:` / `DROP:` log line and the matching `prop:`-prefixed disposition in the L2 tally. The canonical schema has no slot for these today; their dispositions are tracked in §8 for the downstream renderer / custom-type work that unblocks them. (The original plan was to preserve them with `# TODO(postpone-§8.N)` markers, but the loader's strict unmarshal would reject unknown keys; dropping at migration time and keeping the design intent in §8 is the chosen approach.)
10. **One-off cleanup** — drop the redundant `include: true` from `fvFBRoute` during migration. The class has `IdentifiedBy = ["fbrPrefix"]` and `RnFormat = "pfx-..."`, so the `provider.go.tmpl` registry conditional `or (and .IdentifiedBy (not (and .MaxOneClassAllowed (hasPrefix .RnFormat "rs")))) .Include` is already satisfied by the first branch.

Coverage today: the loader's `UnmarshalStrict` returns clean against every emitted YAML in the corpus; the L1 unknown-key tally is empty and the L2 per-section totals match the file counts in §§1–7. The only remaining drift is editorial — the `exclude_targets` polymorphic-prediction verification and the optional-property `default_values` fallback in `setParentDn` are tracked in their own bullets above.

### 9.3 Open decision points

None block the migration script. The walk-through resolved every key in §4–§7 to a concrete action; the §8 POSTPONE entries are deferred to downstream work, listed here for traceability:

| POSTPONE entry | Blocked on | Owner area |
|---|---|---|
| `static_custom_type` named-type (§8.1) | Design for named semantic-equality variants (registry, enum, or YAML-driven custom types). | Custom-type framework. |
| `datasource_required` top-level (§8.2) | Renderer learns `Class.Artifacts` + `cleanDirectory` whitelist equivalent. | Render layer. |
| `ignore_custom_type_docs` (§8.3) | Docs renderer picks a uniform rule for `SemanticEquality` valid-values + validator-range output (uniform suppress, localname=wire DERIVE, or per-property ADD). | Docs renderer / editorial. |
| `example_value_overwrite` (§8.4) | Example renderer applies the `relationInfo.toMo` rule (collapses 5 of 7); narrower fields decided for the 2 editorial cases alongside polymorphic-relation handling. | Example/test renderer. |
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
| Property-level YAML keys declared in files | 20 (+ 7 entries inside `legacy_definitions/properties/resource_name_overwrite.yaml`) | §1, §2, §3, §4, §5, §6, §8 |
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
| `required_by_custom_type_in_test` | [IncludeInCustomTypeTest](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3681) (v2.19.0:3681) | Obsolete — covered by the per-property `test_config.custom_type` bucket (already in the §2 `test_values` row). | §3 |
| `target_resource_name` (nested under `targets[]`) | [GetTestTargetValue](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L366) (v2.19.0:366) — overrides the test-resource label in `aci_<X>.test_<X>_<i>.<attr>` references | Obsolete — 0 v2.19.0 files declare it. The new pipeline derives the test-resource label from the resolved dependency's auto-resolved resource name; no per-target override slot is needed. | §3 |
| `version_mismatch` (nested under `test_values[<prop>]` and `test_values[].children[].<child>`) | [LookupTestValue](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L800) (v2.19.0:800, 929) — alternate test-value bucket selected when the target APIC version doesn't match the property's primary value | Obsolete — 0 v2.19.0 files declare it. The new test renderer expresses version-conditional values via the per-property `TestValueEntry.Versions` field (already on the struct); no nested bucket needed. | §3 |
| `child_versions` (nested under `test_values[].children[]`) | [GetChildVersion](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L231) (v2.19.0:231) — per-child supported-version filter for nested test blocks | Obsolete — 0 v2.19.0 files declare it. The new `class.test_config.children[]` shape carries version filters at the entry level when needed; no separate bucket. | §3 |
| `deletable_child` (nested under `test_values[].children[].<child>`) | [CheckDeletableChild](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L444) (v2.19.0:444, recursive scan) — marks a nested child block as exercising the delete path | Obsolete — 0 v2.19.0 files declare it. The new test renderer derives deletion-testability from the child's own meta `isDeletable` / `allow_delete` flag (mechanism for the deletion test step is part of the test scenario framework rewrite). | §3 |

### 10.3 Nested-key inventory (parents / targets / multi_parents)

Listed here so the parent / target / multi-parent migrations have a single reference. Each nested key already lands in §2 or §4; this is a spot-check, not a separate disposition.

| Outer key | Nested keys read by v2.19.0 | New target |
|---|---|---|
| `parents` ([v2.19.0:3198–3219](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3198)) | `class_name` (entry key), `parent_dependency`, `class_in_parent`, `parent_dependency_name`, `parent_dn`, `target_classes`, `properties` | `class.test_config.dependencies[]` shape (§2 `parents` row). `properties` → `config_overrides` (§2.2 step 3); `target_classes` is **read** by the script for polymorphic-same-type detection (§8.6) before being dropped from output (§2.2 parents step 6). |
| `targets` ([v2.19.0:3282–3320](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3282)) | `class_name` (entry key), `relation_resource_name`, `shared_classes`, `parent_dependency`, `parent_dependency_dn_ref`, `overwrite_parent_dn_key`, `static`, `target_dn`, `target_dn_ref`, `target_dn_overwrite_docs`, `properties` | Folded into the per-class `relation_info.to_classes` + `test_config.dependencies[]` shapes (§2 `relationship_classes` / `targets` row). `target_dn` → `reference` (§2.2 targets step 1); `properties` → `config_overrides` (§2.2 targets step 3); `parent_dependency_dn_ref` → §2.2 targets step 2; `target_dn_overwrite_docs` is §8.4 POSTPONE; `target_dn_ref` and `shared_classes` are asserted-empty drops (§2.2 targets step 6). |
| `multi_parents` ([v2.19.0:3086–3088](https://github.com/CiscoDevNet/terraform-provider-aci/blob/v2.19.0/gen/generator.go#L3086)) | `test_type`, `wrapper_class` | `ClassDefinition.ParentDnVariants[]` (§4.1 `ParentDnVariantDefinition`). |
