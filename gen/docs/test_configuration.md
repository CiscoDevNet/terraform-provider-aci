# Test Configuration System

> **Work in progress.** This document is actively being drafted and the structure of this `gen/docs/` folder (file layout, naming, cross-linking) is still to be determined. Content and organization may change.

This document describes the test data model used by the ACI provider code generator to produce acceptance tests. It covers property test values, class-level test dependencies, child test values, the loading pipeline, and when to provide manual YAML overrides.

---

## 1. Overview

The generator gathers per-property and per-class **test data** describing the values, dependencies, and child blocks needed to exercise a resource. Future test templates consume this data to render acceptance tests in `resource_aci_<name>_test.go`. This document covers the data model only; concrete step ordering, step count, and per-scenario assertions are out of scope and are tracked in Section 9.

The two primitives used by every consumer of the data:
- If `ConfigInclude == true` → render `ConfigValue` in HCL config.
- Always assert `AssertValue` in state checks.

---

## 2. Property TestValues

### 2.1 TestValueEntry

Each property's test data is captured per-step as a slice of `TestValueEntry`:

```go
type TestValueEntry struct {
    ConfigValue   string          // Value written into HCL config
    ConfigInclude bool            // Whether to include this entry in config
    AssertValue   string          // Expected value in Terraform state after apply
    ValueType     ValueRenderType // StringValue (quoted) or ReferenceValue (unquoted expression)
}
```

### 2.2 TestValues

```go
type TestValues struct {
    Create   []TestValueEntry  // Data bucket: full attribute values for an all-attributes configuration.
    Update   []TestValueEntry  // Data bucket: alternate values for scenarios that change attributes in place.
    Default  []TestValueEntry  // Data bucket: required-only scenarios. Required props carry Create values; optional props carry server-default assertions with ConfigInclude=false.
    ForceNew []TestValueEntry  // Data bucket: alternate values for scenarios that intentionally change a replace-triggering attribute (today: parent_dn second instance; non-parent props typically copy Create).
    Legacy   []TestValueEntry  // Data bucket: per-step values for legacy alias attributes exposed via state_upgrades (Functioning / Frozen). See §2.4 "Legacy bucket auto-derivation".
}
```

- For scalar properties, each slice has exactly 1 entry.
- For set-typed (bitmask) properties, each slice has multiple entries (one per set member).

### 2.3 TestValues buckets

Each slice on `TestValues` is a **data bucket**, not a literal test step. A future template may consume one, several, or none of these buckets when rendering a scenario.

- `Create` — full set of attribute values for an all-attributes configuration.
- `Update` — alternate values, used wherever a scenario needs the resource to change without destroy+recreate of the parent.
- `Default` — per-attribute data for "required-only" scenarios: required props carry Create values; optional props carry server-default assertions with `ConfigInclude=false`.
- `ForceNew` — alternate values for scenarios that intentionally change a replace-triggering attribute. Today only `parent_dn` (second instance) and `tDn` (single-target relations) are auto-populated differently from Create; other properties typically copy Create.
- `Legacy` — alternate values for scenarios that exercise the property under its deprecated legacy alias instead of its current name. Auto-derived from `Create` when the legacy alias has the same Terraform type; explicit YAML required when types diverge. Nil when the property has no testable legacy alias. See §2.4 "Legacy bucket auto-derivation" and `state_upgrades.md` §11 for a worked example.

### 2.4 Auto-Derivation Rules

When no explicit `test_config` is provided on a property definition, values are auto-derived:

#### From ValidValues (enum properties)
- Pick first two values alphabetically.
- Create = first value, Update = second value.
- If only 1 valid value: Create = Update = that value.

#### Set-typed (bitmask) properties
- Create: first 2 members from ValidValues (alphabetically).
- Update: third and fourth members (or overlap if fewer than 4).

#### Free-form string properties (no ValidValues)
- Create = `"<attribute_name>_1"`, Update = `"<attribute_name>_2"` for both Required and Optional properties. Distinct buckets are the source of truth for downstream consumers: list-type child instance 0 takes Create and instance 1 takes Update (see §4.2), so this scheme keeps sibling instances on distinct APIC Dns without any per-instance value mangling. The (future) Update test step at the parent root chooses which bucket to apply for required vs optional properties; emitting both values here is safe in either case and explicit `test_config` overrides on either bucket continue to win.

#### Default bucket auto-derivation
- **Required properties:** ConfigInclude=true, ConfigValue and AssertValue copied from Create.
- **Optional properties:** ConfigInclude=false, AssertValue = server default:
  - If `DefaultValues` defined in documentation: use that value.
  - Otherwise: AssertValue = `""` (empty string assumed as server default).

#### ForceNew bucket auto-derivation
Non-parent properties copy their Create entry into ForceNew. Only `parent_dn` (and `tDn` for single-target relations) is explicitly populated with a second reference, drawn from the second auto-resolved Parent/Target dependency. Per-property `force_new` overrides exist in YAML but are rarely needed.

#### Legacy bucket auto-derivation

Driven by `state_upgrades` entries on the class definition (see `state_upgrades.md`). For each property:

1. **Explicit YAML wins.** `test_config.legacy` always replaces auto-derivation. This is the sole path when the legacy alias has a different Terraform type than the current attribute (auto-derive cannot guess a valid HCL shape).
2. **Skip when not testable.** If the property has no `TestValues`, has `IgnoreInTest=true`, or is `ReadOnly`, Legacy stays nil.
3. **Skip when no testable legacy alias exists.** A legacy alias is testable when at least one `StateUpgradeValue` entry has `Status` `Functioning` or `Frozen` AND renames the attribute. Same-name entries (type-only or restriction-only changes) are not separately testable because they share the current attribute name. `Removed` entries are migration-only and never produce a Legacy bucket.
4. **Skip with a generator warning when types diverge.** When any non-`Removed` legacy alias declares a `legacy_type` different from the current property's type, auto-derivation refuses to guess (cloning Create would produce HCL of the wrong shape). The generator logs a `Warn` and requires explicit `test_config.legacy` in the property definition. Templates can detect this case by observing `len(TestValues.Legacy) == 0` despite the property carrying a `Functioning`/`Frozen` `StateUpgradeValue` with a renamed alias.
5. **Otherwise: clone Create.** Each entry in `Create` is copied into `Legacy` field-for-field (including `Versions`). The legacy alias re-uses the current attribute's values, just rendered under its prior name.

Legacy runs at the end of Loop 3 (`setPropertyTestValues`), after `parent_dn` / `tDn` placeholder resolution, so any dependency-derived `Create` entries are already wired and safe to clone.

### 2.5 IgnoreInTest vs SupportedVersions

- `IgnoreInTest`: Property is **never** testable regardless of version (e.g., broken API behavior). Set via `test_config.ignore_in_test: true` in YAML.
- `SupportedVersions`: Property is testable only on specific APIC versions. Test templates use this to conditionally include/exclude. No impact on test data generation.

### 2.6 YAML Input Format

On `PropertyDefinition` in the class YAML file:

```yaml
properties:
  descr:
    test_config:
      ignore_in_test: false  # default; set true to skip this property entirely
      create:
        - config_value: "my_description"
          config_include: true    # default when omitted
          assert_value: ""        # default: same as config_value
          value_type: "string"    # default; or "reference" for unquoted expressions
      update:
        - config_value: "updated_description"
      default:
        - config_value: ""
          config_include: false
          assert_value: "default_from_server"
      legacy:
        - config_value: "my_description"   # required only when the legacy type diverges from the current attribute
```

- `config_include`: pointer-based — `nil` defaults to `true`.
- `assert_value`: empty string defaults to `config_value`.
- `value_type`: `"string"` (default, quoted) or `"reference"` (unquoted HCL expression).
- `legacy`: optional. Overrides Legacy bucket auto-derivation (see §2.4). Required when the legacy alias declares a `legacy_type` different from the current attribute's type, because auto-derive cannot guess the prior HCL shape.

#### All-or-nothing for standard buckets

`create`, `update`, `default`, and `force_new` form a single unit at load time. Either supply none of them (the property auto-derives all four) or supply all four (full manual override). Mixing — for example supplying only `create` — fails class validation with an error naming the present and missing buckets:

> `property "<name>" test_config supplies [create] but omits [default force_new update]; all four standard buckets (create, update, default, force_new) must be supplied together`

The rule exists because partial overrides silently nil the unspecified buckets, leaving downstream consumers with no data to render the corresponding scenarios. `legacy` is independent and is not part of this check; it can be supplied alone or omitted regardless of what the standard buckets look like.

---

## 3. Class Test Dependencies

### 3.1 TestDependency Struct

```go
type TestDependency struct {
    Class           *ClassName
    Reference       string                // e.g. "aci_tenant.test.id" or static DN
    ReferenceType   ReferenceTypeEnum     // StaticReference, ResourceReference, DataSourceReference
    Role            TestDependencyRoleEnum // Parent or Target (top-level only; zero for nested)
    Dependencies    []*TestDependency     // Recursive prerequisites
    ConfigOverrides map[string]string     // Property overrides for dependency's HCL config
    Children        map[string]*TestChild // Child block overrides for dependency resource
}
```

### 3.2 Role (Parent/Target)

- **Parent:** Provides the `parent_dn` attribute. Only set on top-level entries.
- **Target:** Provides the `target_dn` attribute (relation classes). Only set on top-level entries.
- Nested dependencies (inside `Dependencies`) are pure prerequisites — Role is always zero.

### 3.3 ReferenceType

| Type | Description | Example |
|------|-------------|---------|
| `ResourceReference` | Terraform resource attribute path | `aci_tenant.test.id` |
| `DataSourceReference` | Terraform data source attribute path | `data.aci_tenant.test.id` |
| `StaticReference` | Hardcoded DN string | `uni/vmmp-VMware/dom-domain_1` |

Auto-resolution **never** produces `StaticReference` — it always requires manual YAML.

### 3.4 Auto-Resolution Rules

#### Parents (from `Class.Parents`)
- **First parent class:** 2 instances for ForceNew testing:
  - `aci_<resource_name>.test.id` (Role=Parent)
  - `aci_<resource_name>.test_2.id` (Role=Parent)
- **Second parent class (if exists):** 1 instance for compatibility testing:
  - `aci_<resource_name_2>.test.id` (Role=Parent)
- **3+ parents:** Only the first 2 parent classes in `Class.Parents` order are auto-resolved. Additional parent classes are silently dropped (Trace-level log only, no diagnostic). Add them via explicit `test_config.dependencies` if they need to appear in HCL — note that even then, `parent_dn` auto-wiring still only uses the first 2 entries (see Section 6).
- **Globally-excluded singletons** (e.g., `polUni`, `fabricInst`): filtered out of `c.Parents` upstream by `GlobalMetaDefinition.ExcludeParents` in `setParents`, before dependency resolution runs. They never reach the dependency resolver.
- **Parent with no resource and no `NoMetaFile` entry:** silently skipped (Trace log only) at the dependency-resolution step. This is a fallback safety net; in practice global `ExcludeParents` should cover such classes.

#### Targets (from `Class.Relation.ToClasses`)
- **Single-target** (`len(ToClasses) == 1`): 2 instances for toggle testing:
  - `aci_<resource_name>.test.id` (Role=Target)
  - `aci_<resource_name>.test_2.id` (Role=Target)
- **Multi-target** (`len(ToClasses) > 1`): No auto-resolution. Must provide explicit `test_config.dependencies`. Note: `setTargetDn` only wires `tDn` from the first two declared Target-role entries (Create/Default/ForceNew = `targets[0]`; Update = `targets[1]`). Additional targets render as HCL prerequisite resources but are never assigned to `tDn` — to exercise them, override the `tDn` property's `test_config` directly with `{{<reference>}}` placeholders or split into multiple test scenarios.

#### Recursive resolution
Each dependency's own parents are resolved recursively as nested `Dependencies` (Role=0). The same `seen` map is shared to prevent infinite recursion and enable DAG deduplication.

### 3.5 DAG Deduplication

A `seen map[string]*TestDependency` (keyed by `Reference` string) ensures shared ancestors are rendered only once. When the same reference appears in both parent and target chains, the same `*TestDependency` pointer is reused. The reference string is inherently unique.

### 3.6 ConfigOverrides

Property overrides for a dependency resource's HCL configuration. When empty, the dependency renders using its class's own `TestValues.Create`. When populated, specified properties override auto-derived values.

Placeholders in `ConfigOverrides` values use the `{{<reference>}}` format and are resolved against the entire DAG `seen` map by reference key.

### 3.7 Placeholder Syntax

All placeholders use the same unified format: `{{<reference>}}`. The reference must exactly match a `TestDependency.Reference` string in the resolution scope. Leading/trailing whitespace inside the braces is trimmed (e.g., `{{ aci_tenant.test.id }}` resolves correctly).

| Context | Resolved Against |
|---------|------------------|
| Dependency `config_overrides` | DAG `seen` map (flat lookup by reference key) |
| Property `test_config` values | Class's TestDependencies (recursive DAG search) |
| Child override `properties` | Parent class's TestDependencies (recursive DAG search) |

### 3.8 YAML Input Format

On `ClassDefinition`:

```yaml
test_config:
  dependencies:
    - class_name: fvTenant
      reference: "aci_tenant.test.id"
      reference_type: "resource"   # "resource" (default), "static", "data_source"
      role: "parent"               # required at top level: "parent" or "target"
      config_overrides:
        name: "test_tenant"
      dependencies:                # nested prerequisites (recursive, same struct)
        - class_name: vmmDomP
          reference: "uni/vmmp-VMware/dom-test_domain"
          reference_type: "static"
          # role must be empty at nested level
    - class_name: fvBD
      reference: "aci_bridge_domain.test.id"
      reference_type: "resource"
      role: "target"
      config_overrides:
        tenant_dn: "{{aci_tenant.test.id}}"  # resolved against DAG
      dependencies:
        - class_name: fvTenant
          reference: "aci_tenant.test.id"
          reference_type: "resource"
```

By default, `test_config.dependencies` is **additive**: explicit definitions are processed first, then auto-resolution fills in the remainder (skipping already-defined references). To skip auto-resolution entirely, set `replace_auto_resolved: true`:

```yaml
test_config:
  replace_auto_resolved: true
  dependencies:
    - ...
```

---

## 4. Child Test Values

### 4.1 Instance Count Logic

- `IsSingleNestedWhenDefinedAsChild == true` → 1 instance (1:1 relation, e.g., `fvRsBd`)
- `IsSingleNestedWhenDefinedAsChild == false` → 2 instances (list/set, e.g., `tagAnnotation`)

### 4.2 Auto-Derivation

- **Instance 0:** Uses child class's own `TestValues.Create` values.
- **Instance `i > 0`** (list-type only): Uses child class's own `TestValues.Update` values (falls back to Create if Update is nil so override-only Update buckets still behave).

No per-instance disambiguation is applied. String-typed naming attributes (members of `IdentifiedBy`) auto-derive to distinct Create/Update values via §3 (`"<attr>_1"` and `"<attr>_2"`), so list-type children (`tagAnnotation`, `tagTag`, …) land on distinct APIC Dns naturally. Reference-typed identifiers are already disambiguated upstream (e.g. `aci_contract.test.name` vs `aci_contract.test_2.name`) and flow through verbatim. Explicit `test_config` overrides on either bucket take precedence.

Recursive: Each child class's own `Children` are resolved the same way and attached as nested `TestChild` entries.

### 4.3 Child-Driven Dependency Collection

During `resolveChildTestValues()`, any child instance property with `ValueType == ReferenceValue` that references a resource not already in the parent's `TestDependencies` is auto-collected:
- Look up the reference in the child class's own `TestDependencies`.
- If found, add the `*TestDependency` (pointer shared) to the parent's `TestDependencies`.

This ensures the parent's test HCL includes all resource blocks needed by its children.

### 4.4 Override Format

```yaml
test_config:
  children:
    fvRsBd:                        # keyed by child class name
      instances:                   # FULL REPLACEMENT when present
        - properties:
            target_dn: "{{aci_bridge_domain.test.id}}"
          children:                # recursive grandchild overrides
            tagAnnotation:
              instances:
                - properties:
                    key: "custom_key_0"
                    value: "custom_value_0"
                - properties:
                    key: "custom_key_1"
                    value: "custom_value_1"
    tagAnnotation:
      instances:
        - properties:
            key: "single_key"
            value: "single_value"
```

**Full replacement semantics:** When `instances` is specified, ALL auto-derived instances are discarded and replaced by exactly what's listed. This avoids the "which instance am I overriding?" problem with positional matching.

No `instance_count` — count is always determined by the child class's `IsSingleNestedWhenDefinedAsChild`. To change the count, modify the child class definition.

### 4.5 Placeholder Resolution in Child Overrides

Placeholders in child override property values (`{{<reference>}}`) are resolved against the **parent class's** `TestDependencies` (recursive DAG search) — which at this point contains both own and child-collected dependencies.

---

## 5. Pipeline / Loading Order

`NewDataStore` runs a meta-retrieval phase followed by four sequential loops over `ds.Classes`:

1. **Meta retrieval** (before any loops): `setMetaHost`, `retrieveEnvMetaClassesFromRemote`, `refreshMetaFiles` populate `gen/meta/` from the remote pubhub host on demand.
2. **Loop 1 — Load classes** (`loadClasses` → per file `loadClass` → `NewClass` → `setClassData`): parse each meta file and run the class-level setters in single-class scope (no cross-class lookups):
   - `setProperties(ds)` — `Property.TestValues` from `test_config` or auto-derive
   - `setRelation()` — `Relation.ToClasses` populated
   - `setResourceName()` — `ResourceName` populated
   - `setChildren()` — `Children` list populated
   - `setParents()` — `Parents` populated, with global `ExcludeParents` applied
   - ... (all other existing setters)
3. **Loop 2 — Documentation:** for each loaded class, call `setDocumentation(ds)`. Runs after Loop 1 globally because doc rendering reads child class info from the DataStore.
4. **Loop 3 — Test dependencies + property test values** (`setTestData`, first loop): per class, call `setTestDependencies(ds)` then `setPropertyTestValues(ds)` back-to-back.
   - `setTestDependencies` resolves the `TestDependencies` DAG (shared pointers via `seen`) and resolves `ConfigOverrides` placeholders against the DAG.
   - `setPropertyTestValues` wires `parent_dn` from the first Parent dependency, wires `tDn` from Target dependencies (single-target only), and resolves `{{<reference>}}` placeholders in property `TestValues`.
   - These two share a loop because `setPropertyTestValues` reads only the class's own `TestDependencies` (just populated) — no cross-class read into another class's test data.
5. **Loop 4 — Child test values + completeness validation** (`setTestData`, second loop): per class, call `setChildTestValues(ds)` then `validateTestCompleteness(ctx)`.
   - `setChildTestValues` builds `TestChildren` from child class `TestValues`, applies child overrides from `ClassDefinition.TestConfig.Children`, collects child-driven dependencies into the parent's `TestDependencies`, and resolves `{{<reference>}}` placeholders in child instance properties.
   - `validateTestCompleteness` reports unresolved placeholders in `ConfigOverrides`, properties, and children.

Loop 4 must follow Loop 3 globally because `buildChildInstance` reads property `TestValues` (including wired `tDn`/`parentDn`) from **other** classes, so every class must have finished Loop 3 first.

---

## 6. When to Provide Manual YAML

| Scenario | What to provide |
|----------|----------------|
| Multi-target relations | `test_config.dependencies` with Role=target for each target |
| Static DN references | `test_config.dependencies` with `reference_type: "static"` |
| Dependency needing child-driven targets | Manual `dependencies` on the dependency (NOT auto-collected for dependencies, only for resource-under-test) |
| Dependency needing children to be valid | Explicit `children` block on the dependency definition (children are NOT auto-populated for dependency resources, even if the dependency's own class has children defined) |
| Non-standard test values | `test_config` on PropertyDefinition (server normalization, special characters) |
| Server default differs from empty string | `test_config.default` on PropertyDefinition with both `assert_value: "<expected>"` AND explicit `config_include: false`. Omitting `config_include` defaults to `true` and would force the value into HCL config, defeating the point. |
| 3+ parent classes | `test_config.dependencies` will render the additional parent resource in HCL as a prerequisite, but `parent_dn` auto-wiring still only uses the first 2 Parent-role entries (Create/Update/Default use [0], ForceNew uses [1]). To actually switch between a 3rd parent type, also override the `parent_dn` property's `test_config` directly with a `{{<reference>}}` placeholder pointing at the desired entry. |
| Legacy alias with divergent type | `test_config.legacy` on PropertyDefinition. Auto-derivation skips with a `Warn` log when any `Functioning`/`Frozen` `StateUpgradeValue` carries a `legacy_type` different from the current attribute's type — supply the prior-shape HCL values explicitly. See §2.4 "Legacy bucket auto-derivation". |

---

## 7. Examples

### 7.1 Simple resource (fvTenant — no parents, no targets)

No YAML needed. Properties auto-derive from ValidValues or free-form rules. No TestDependencies generated: fvTenant's only meta-declared parent is `polUni`, which is filtered out of `c.Parents` by global `ExcludeParents` in `setParents`, so the dependency resolver sees an empty parent list.

### 7.2 Resource with parent chain (fvAEPg → fvAp → fvTenant)

Auto-resolved:
```
TestDependencies:
  [0] aci_application_profile.test.id (Role=Parent)
       └─ aci_tenant.test.id (Role=0, prerequisite)
  [1] aci_application_profile.test_2.id (Role=Parent)
       └─ aci_tenant.test.id (Role=0, shared pointer with [0]'s dep)
```

`parent_dn` property auto-wired:
- Create/Update/Default = `aci_application_profile.test.id` (ReferenceValue)
- ForceNew bucket holds `aci_application_profile.test_2.id`; templates consume it when rendering a parent-switch scenario.

### 7.3 Relation resource with target (fvRsBd → fvBD)

Auto-resolved:
```
TestDependencies:
  [0] aci_application_profile.test.id (Role=Parent)     # from parent chain
  [1] aci_application_profile.test_2.id (Role=Parent)
  [2] aci_bridge_domain.test.id (Role=Target)
       └─ aci_tenant.test.id (prerequisite, shared)
  [3] aci_bridge_domain.test_2.id (Role=Target)
       └─ aci_tenant.test.id (shared pointer)
```

`tDn` property auto-wired:
- Create = `aci_bridge_domain.test.id` (ReferenceValue)
- Update = `aci_bridge_domain.test_2.id` (ReferenceValue)
- Default = `aci_bridge_domain.test.id` (required, stays in config)

### 7.4 Dependency with nested blocks (fvBD needing fvRsCtx → fvCtx)

When a dependency resource itself needs children to be valid, declare them under `children` on the dependency definition. Children are NOT auto-populated for dependency resources — you must explicitly redeclare any required children, even if the dependency's own class has them defined.

```yaml
test_config:
  dependencies:
    - class_name: fvBD
      reference: "aci_bridge_domain.test.id"
      reference_type: "resource"
      role: "target"
      children:
        fvRsCtx:
          instances:
            - properties:
                tnFvCtxName: "{{aci_vrf.test.id}}"
      dependencies:
        - class_name: fvTenant
          reference: "aci_tenant.test.id"
          reference_type: "resource"
        - class_name: fvCtx
          reference: "aci_vrf.test.id"
          reference_type: "resource"
```

### 7.5 Multi-target override example (fvRsDomAtt)

`fvRsDomAtt` is the canonical real-world multi-target relation: its meta `toMo` is the abstract `infra:DomP` and the class targets four concrete domain types — `vmmDomP`, `physDomP`, `fcDomP`, `l2extDomP`. Today fvRsDomAtt's class YAML uses the legacy `relationship_classes` field (consumed only by the SDKv2 templates); once migrated to the framework's `relation_info.to_classes`, `len(Relation.ToClasses) > 1` disables target auto-resolution and every concrete target instance must be defined explicitly. Parents (`fvAEPg` → `fvAp` → `fvTenant`) still auto-resolve in additive mode.

```yaml
test_config:
  dependencies:
    - class_name: vmmDomP
      reference: "aci_vmm_domain.test.id"
      reference_type: "resource"
      role: "target"
    - class_name: physDomP
      reference: "aci_physical_domain.test.id"
      reference_type: "resource"
      role: "target"
    - class_name: fcDomP
      reference: "aci_fc_domain.test.id"
      reference_type: "resource"
      role: "target"
    - class_name: l2extDomP
      reference: "aci_l2_domain.test.id"
      reference_type: "resource"
      role: "target"
```

`tDn` property auto-wired from the first two declared targets:
- Create/Default/ForceNew = `aci_vmm_domain.test.id`
- Update = `aci_physical_domain.test.id`

`fcDomP` and `l2extDomP` render as HCL resource blocks (DAG prerequisites) but are never assigned to `tDn`. To exercise all four, override the `tDn` property's `test_config` with explicit `{{<reference>}}` placeholders, or split into multiple test scenarios.

---

## 8. Custom Test Cases

### 8.1 File naming convention

- **Generated** (overwritten on regen): `resource_aci_<name>_test.go`
- **User-maintained** (never touched): `resource_aci_<name>_custom_test.go`

The generator skips/preserves files matching `*_custom_test.go` during cleanup.

### 8.2 When to use custom tests

- Import edge cases (non-standard DN formats)
- Error validation (invalid attribute values, conflicting attributes)
- State migration testing
- Provider-upgrade scenarios
- Complex multi-resource interaction patterns

### 8.3 Coexistence

Custom files share the same package and can reference generated test helpers (provider config functions, test check builders, etc.).

---

## 9. Future Test Scenarios

No test templates have been implemented yet — the data model in Sections 1–7 is built, but nothing consumes it to render `resource_aci_<name>_test.go` files. The list below covers both the **baseline scenarios** the template layer must render first and the **gaps** to address once the baseline is in place.

### Baseline scenarios (not yet implemented)

The minimum scenarios every generated `*_test.go` should render once templates exist:

- **Create step** — apply a config built from each property's `TestValues.Create` bucket (plus all auto-resolved/explicit dependencies and child blocks); assert state matches `AssertValue` for every entry with `ConfigInclude=true`, and assert server-default `AssertValue` for entries with `ConfigInclude=false`.
- **Update step** — re-apply with values from `TestValues.Update`; assert in-place update succeeds and state matches the new `AssertValue`s. For properties where Create == Update, assert no diff.
- **Required-only / Default step** — apply a config built from `TestValues.Default` (only required props in HCL; optional props omitted); assert optional props carry their server-default `AssertValue` in state.
- **ForceNew / parent-switch step** — apply a config that swaps the `parent_dn` (and `tDn` for single-target relations) to the second instance from `TestValues.ForceNew`; assert destroy+recreate of the resource-under-test.
- **Import step** — `ImportState`/`ImportStateVerify` after one of the above apply steps; assert imported state matches applied state.
- **Data source step** — apply the resource then read it via `data.aci_<name>`; assert every attribute (including children) matches the resource's state.
- **Cleanup** — final destroy step (typically implicit via the test framework, but the template must not leave orphan dependencies).

### Import testing gaps
- Import with non-default attribute values (currently imports after Reset step which has defaults)
- Import after child modifications (partial child presence)
- Verify `parent_dn` reconstruction correctness explicitly (currently only implicit via `ImportStateVerify`)

### Child lifecycle gaps
- Child attribute updates (changing a value within an existing child block between steps, not add/remove)
- Incremental child addition (adding children one at a time, not all at once)
- Child ordering stability verification (reorder list-type children, verify no diff)
- Data source children (currently data source tests don't verify child/nested block attributes)

### Custom type / semantic equivalence gaps
- Custom type update (change from one non-canonical form to another)
- Custom type in children (nested block attributes with custom types)
- Custom type import verification (import state matches normalized form)

### Legacy attribute gaps

The `TestValues.Legacy` data bucket is now populated (auto-derived from `Create` when types match, or supplied via `test_config.legacy` when they diverge — see §2.4). The remaining gaps live in the template layer that consumes it:

- Render a scenario that applies a config under the legacy alias and asserts the current-name state attribute receives the same value (Functioning aliases) or stays at the prior value with no device write (Frozen aliases).
- Legacy-to-new attribute migration in a single test (set via legacy → read via new attribute).
- Conflicting legacy + new attribute error validation (`ConflictsWith` enforcement on Functioning aliases).
- Removed-alias migration coverage: state with the removed name upgrades cleanly to the current schema (driven by `state_upgrades`, no Legacy bucket involvement).

### Step rendering / template scope

The data model in this document gathers test data per scenario; concrete step ordering and assertions live in the (future) test templates layer. Items below belong to that layer:

- Concrete step ordering: all-attributes create → update → required-only → parent switch → import → cleanup.
- Parent-switch step rendering: consumes the `ForceNew` bucket (second Parent-role entry from `TestDependencies`) and asserts destroy+recreate.
- Multi-parent type compatibility scenario when a class has 2+ parent types (renders the second Parent class with a single `.test.id` instance and verifies the resource applies under it).
- Explicit plan-shape assertion that a replace-trigger step shows destroy+create (currently only implicit via apply success).
- RequiresReplace coverage strategy:
  - Implicit coverage when a property's Create and Update values differ (Terraform auto-plans replace as part of the update step). Free-form string properties always satisfy this because §3 auto-derives distinct `"<attr>_1"` / `"<attr>_2"` values.
  - **Gap to address:** RequiresReplace properties whose `test_config` override pins Create == Update, or whose `ValidValues` only contain a single member. The (future) template Update step should detect these and either widen the value set or emit an explicit replace-trigger scenario.

### Generator code follow-ups (not template scope)

- **Silent-skip diagnostics for dropped parents:** Both "3+ parents" (auto-resolver bails at `i >= 2`) and "parent with no resource and no `NoMetaFile` entry" (`getResourceNameForClass` returns empty) are dropped with only a Trace-level log. Consider emitting a Warning-level diagnostic in either case so coverage gaps are visible to whoever runs the generator.
- **Multi-target only exercises the first two `Target`-role entries via `tDn` Create/Update.** For relations with 3+ concrete targets (e.g., `fvRsDomAtt`'s four domain types), the additional targets render in HCL but are never assigned to `tDn`. Decide whether to auto-emit additional scenarios per remaining target, or document/codify the override pattern that walks `tDn` through every target.

### ForceNew / parent_dn
- Multi-parent type compatibility (verify resource works under different parent types within same test)

### Attribute update
- Value-to-value update (All_v1 → All_v2 with different non-default values; currently only All → Min → Reset)
- Target DN change for relation resources (update `target_dn` to different target)
- Individual attribute toggling (isolate single attribute change to detect side effects)

### Error validation
- Invalid attribute values (expect specific error messages)
- Version constraint violations (use attribute on unsupported version)
- Invalid `parent_dn` format
- Conflicting attributes (mutually exclusive validators)

### Sensitive/password attributes
- Explicit sensitive value persistence assertion (verify config value survives Read)
- Sensitive attribute update (change password between steps)
- Import with sensitive attributes (verify `ImportStateVerifyIgnore` is correct)

### Plan accuracy
- `PlanOnly: true` steps to verify diff prediction without apply
- Empty plan after no-op apply (verify no spurious diffs)

### Generator tooling / CI hygiene

The `diff` job in `.github/workflows/checks.yml` already runs `go generate` followed by `git diff --exit-code`, which surfaces YAML strict-unmarshal errors, `genLogger.Fatal` calls, accumulated `ctx.Diagnostics` errors, and uncommitted generated-output drift. Remaining gaps:

- Run `go vet ./gen/...` and `staticcheck` in CI to catch shadowed errors, unused returns, and other static issues that unit tests don't cover. Today only `gofmt` runs in the `build` job.
- Run generator unit tests (`go test ./gen/...`) with `-race` in CI. Today `-race` is only applied to the acceptance job against `internal/provider`; there is no job that exercises the `gen` package's own tests, so the logger's `ResetForTest()` isolation is not actually verified under race detection.
