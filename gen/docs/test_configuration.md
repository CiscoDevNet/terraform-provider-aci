# Test Configuration System

> **Work in progress.** This document is actively being drafted and the structure of this `gen/docs/` folder (file layout, naming, cross-linking) is still to be determined. Content and organization may change.

This document describes the test data model used by the ACI provider code generator to produce acceptance tests. It covers property test values, class-level test dependencies, child test values, the loading pipeline, and when to provide manual YAML overrides.

---

## 1. Overview

The generator gathers per-property and per-class **test data** describing the values, dependencies, and child blocks needed to exercise a resource. Future test templates consume this data to render acceptance tests in `resource_aci_<name>_test.go`. This document covers the data model only; concrete step ordering, step count, and per-scenario assertions are out of scope and are tracked in Section 9.

The two primitives used by a future test renderer are:

- If `ConfigInclude == true` → render `ConfigValue` in HCL config.
- Always assert `AssertValue` in state checks.

---

## 2. Property TestValues

### 2.1 TestValueEntry

Each property's test data is captured as ordered entries within a lifecycle scenario:

```go
type TestValueEntry struct {
    ConfigValue   string
    ConfigInclude bool
    AssertValue   string
    ValueType     ValueRenderTypeEnum
    Versions      *Versions
}
```

### 2.2 TestValues

```go
type TestValueScenario struct {
    Defined bool
    Entries []TestValueEntry
}

type TestValues struct {
    Create   TestValueScenario
    Update   TestValueScenario
    Default  TestValueScenario
    ForceNew TestValueScenario
    Legacy   TestValueScenario
}
```

- For scalar properties, a defined scenario normally has exactly one entry.
- For set-typed properties, a scenario can have multiple entries, one per set member.
- `Defined` distinguishes an omitted scenario from an explicitly defined empty set. That distinction cannot be represented by inspecting `len(Entries)` alone.

### 2.3 TestValues buckets

Each scenario on `TestValues` is a **data bucket**, not a literal test step. A future template may consume one, several, or none of these buckets when rendering a scenario.

- `Create` — full set of attribute values for an all-attributes configuration.
- `Update` — alternate values, used wherever a scenario needs the resource to change without destroy+recreate of the parent.
- `Default` — per-attribute data for "required-only" scenarios: required properties remain configured; optional properties carry server-default assertions with `ConfigInclude=false`.
- `ForceNew` — alternate values for scenarios that intentionally change a replace-triggering attribute. Today `parent_dn` uses the second Parent dependency when one is available; other properties normally retain their Create values.
- `Legacy` — alternate values for scenarios that exercise the property under its deprecated legacy alias instead of its current name. Auto-derived from `Create` when the legacy alias has the same Terraform type; explicit YAML is required when types diverge. Undefined when the property has no testable legacy alias. See §2.4 "Legacy bucket auto-derivation" and `state_upgrades.md` §11 for a worked example.

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

The generator first uses a non-empty APIC default when one exists. Otherwise it applies these deterministic conventions:

- `annotation` uses `"annotation"` / `"annotation_2"`.
- IP-address properties use `"192.0.2.1"` / `"198.51.100.2"` (or the corresponding prefix form) from the IANA documentation ranges.
- Numeric properties use the validator minimum and, when possible, minimum + 1.
- Other strings use `"<attribute_name>_1"` / `"<attribute_name>_2"`.

Distinct Create and Update values keep repeated child instances on distinct APIC DNs without per-instance string manipulation. Explicit scenario values always take precedence. The migration script also replaces arbitrary publicly routable IPv4 fixtures in the immutable legacy definitions with documentation-range equivalents, making the emitted canonical test values safe for future examples as well.

#### Default bucket auto-derivation
- **Required properties:** ConfigInclude=true, ConfigValue and AssertValue copied from Create.
- **Required unconstrained strings:** use `"test_<meta_property_name>"` in Default to preserve the established required-only convention.
- **Optional properties:** ConfigInclude=false, AssertValue = server default:
  - If `DefaultValues` defined in documentation: use that value.
  - Otherwise: AssertValue = `""` (empty string assumed as server default).

#### ForceNew bucket auto-derivation
Non-parent properties copy their Create entry into ForceNew. Only `parent_dn` is explicitly populated with a second reference, drawn from the second Parent-role dependency. Relation `tDn` uses its first target in ForceNew and its second target in Update. Per-property `force_new` overrides exist in YAML but are rarely needed.

#### Legacy bucket auto-derivation

Driven by `state_upgrades` entries on the class definition (see `state_upgrades.md`). For each property:

1. **Explicit YAML wins.** `test_config.legacy` always replaces auto-derivation. This is the sole path when the legacy alias has a different Terraform type than the current attribute (auto-derive cannot guess a valid HCL shape).
2. **Skip when not testable.** If the property has no `TestValues`, has `IgnoreInTest=true`, or is `ReadOnly`, Legacy stays undefined.
3. **Skip when no testable legacy alias exists.** A legacy alias is testable when at least one `StateUpgradeValue` entry has `Status` `Functioning` or `Frozen` AND renames the attribute. Same-name entries (type-only or restriction-only changes) are not separately testable because they share the current attribute name. `Removed` entries are migration-only and never produce a Legacy bucket.
4. **Skip with a generator warning when types diverge.** When any non-`Removed` legacy alias declares a `legacy_type` different from the current property's type, auto-derivation refuses to guess (cloning Create would produce HCL of the wrong shape). The generator logs a `Warn` and requires explicit `test_config.legacy` in the property definition. Templates can detect this case from `TestValues.Legacy.Defined == false` despite the property carrying a `Functioning`/`Frozen` `StateUpgradeValue` with a renamed alias.
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
          value_type: "string"    # default; "reference" or "expression" render unquoted
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
- `value_type`: `"string"` (default, quoted), `"reference"` (a semantic Terraform address), or `"expression"` (another unquoted HCL expression).
- `legacy`: optional. Overrides Legacy bucket auto-derivation (see §2.4). Required when the legacy alias declares a `legacy_type` different from the current attribute's type, because auto-derive cannot guess the prior HCL shape.

#### Independent scenario overrides

`create`, `update`, `default`, and `force_new` are independent overrides. Omitted scenarios are derived from metadata and the merged Create scenario; an author can therefore provide only the exceptional values. For example, defining only `create` causes Default and ForceNew to be re-derived from that value while an omitted Update retains its normal generated value.

An explicitly empty list is different from omission:

```yaml
properties:
  encapMode:
    test_config:
      update: []
```

For a Set property this means “clear the set” and remains a defined scenario with zero entries. Explicitly empty scalar scenarios are rejected during completeness validation. `legacy` remains independent and can be supplied alone.

#### Future example metadata

The canonical definition contract already records sparse `example_config` overrides and semantic references, although active example rendering is intentionally left for a later commit. Use these fields only when a future public example cannot be inferred from normalized test data:

```yaml
properties:
  scope:
    reference:
      classes: [fvBD, l3extOut]
      example_class: fvBD
      example_attribute: name
  cert:
    example_config:
      full:
        - value: file("certificate.pem")
          value_type: expression
```

`reference.classes` records every valid APIC target class. `example_class` selects the representative Terraform artifact and `example_attribute` defaults to `id`. `example_config.minimum` and `.full` are independent sparse value scenarios; they preserve omission versus an explicitly empty set in the same way as test scenarios. The later example-normalization commit will consume this metadata without changing test inputs or assertions.

---

## 3. Class Test Dependencies

### 3.1 TestDependency Struct

```go
type TestDependency struct {
    Class           *ClassName
    Reference       string
    ReferenceType   ReferenceTypeEnum
    Role            TestDependencyRoleEnum
    TargetClasses   []*ClassName
    Dependencies    []*TestDependency
    ConfigOverrides map[string]string
    Children        map[string]*TestChild
}
```

`TargetClasses` is meaningful on a Parent dependency and records which target classes are compatible with that parent in a polymorphic relation. It is retained for future renderer normalization; it does not replace the dependency's acceptance-test `Reference`.

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
- **First resolvable parent class:** 2 instances for ForceNew testing:
  - `aci_<resource_name>.test.id` (Role=Parent)
  - `aci_<resource_name>.test_2.id` (Role=Parent)
- **Second resolvable parent class (if one exists):** 1 instance for compatibility testing:
  - `aci_<resource_name_2>.test.id` (Role=Parent)
- **3+ resolvable parents:** Only the first 2 resolvable parent classes in `Class.Parents` order are auto-resolved. Unresolvable entries are skipped while searching for those two; later resolvable parent classes are silently dropped (Trace-level log only, no diagnostic). Add them via explicit `test_config.dependencies` if they need to appear in HCL — note that even then, `parent_dn` auto-wiring still only uses the first 2 entries (see Section 6).
- **Globally-excluded singletons** (e.g., `polUni`, `fabricInst`): filtered out of `c.Parents` upstream by `GlobalMetaDefinition.ExcludeParents` in `setParents`, before dependency resolution runs. They never reach the dependency resolver.
- **Parent with no resource artifact:** silently skipped (Trace log only) at the dependency-resolution step. This includes definition-only classes: they can supply Terraform names for explicit references, but they do not prove that a generated resource exists and therefore cannot participate in automatic dependency resolution. This is a fallback safety net; in practice global `ExcludeParents` should cover such classes.

#### Targets (from `Class.Relation.ToClasses`)
- **Single-target** (`len(ToClasses) == 1`): 2 instances for toggle testing:
  - `aci_<resource_name>.test.id` (Role=Target)
  - `aci_<resource_name>.test_2.id` (Role=Target)
- **Multi-target** (`len(ToClasses) > 1`): No auto-resolution. At least one explicit `test_config.dependencies` entry with Role=Target is required. `setTargetDn` wires the first Target into Create/Default/ForceNew and the second, when present, into Update. Additional Target dependencies remain available to a future renderer but are not assigned to a lifecycle bucket automatically.

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
      target_classes: [vmmDomP]     # optional compatible targets for polymorphic relations
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
- `IsSingleNestedWhenDefinedAsChild == false` → 2 instances (set, e.g., `tagAnnotation`)

### 4.2 Auto-Derivation

- **Instance 0:** Uses child class's own `TestValues.Create` values.
- **Instance `i > 0`** (set-type only): Uses child class's own `TestValues.Update` values (falls back to Create when Update is undefined).

No per-instance disambiguation is applied. String-typed naming attributes (members of `IdentifiedBy`) auto-derive to distinct Create/Update values via §2.4 (`"<attr>_1"` and `"<attr>_2"`), so set-type children (`tagAnnotation`, `tagTag`, …) land on distinct APIC DNs naturally. Reference-typed identifiers are already disambiguated upstream (e.g. `aci_contract.test.name` vs `aci_contract.test_2.name`) and flow through verbatim. Explicit `test_config` overrides on either bucket take precedence.

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
      instances:                   # controls count; entries are sparse overlays
        - properties:
            target_dn:
              - config_value: "{{aci_bridge_domain.test.id}}"
          children:                # recursive grandchild overrides
            tagAnnotation:
              instances:
                - properties:
                    key:
                      - config_value: "custom_key_0"
                    value:
                      - config_value: "custom_value_0"
                - properties:
                    key:
                      - config_value: "custom_key_1"
                    value:
                      - config_value: "custom_value_1"
    tagAnnotation:
      instances:
        - properties:
            key:
              - config_value: "single_key"
            value:
              - config_value: "single_value"
```

When one or more `instances` are listed, their count replaces the auto-derived instance count. Each listed entry is a sparse positional overlay on the corresponding auto-derived instance: omitted properties and grandchildren remain derived, while explicitly listed values replace their matching fields. If more entries are listed than the normal cardinality, later entries derive their base property values using the same non-zero-instance behavior as instance 1. An omitted or explicitly empty list leaves auto-derived instances unchanged.

Each property value uses the same ordered entry-list shape as a lifecycle scenario. This preserves every member of a Set and allows an explicit `[]` to represent an empty set.

There is no separate `instance_count`. Auto-derived children use `IsSingleNestedWhenDefinedAsChild` to choose one or two instances; when one or more `instances` are supplied, the number of listed entries is the replacement count.

When a child class needs the same correlated instances in every parent, put those instances in the child class's `test_config.embedded_instances`. This is the canonical form of legacy `test_values_for_parent`; parent-specific `test_config.children` remains available for genuine parent overrides.

### 4.5 Placeholder Resolution in Child Overrides

Placeholders in child override property values (`{{<reference>}}`) are resolved against the **parent class's** `TestDependencies` (recursive DAG search) — which at this point contains both own and child-collected dependencies.

---

## 5. Pipeline / Loading Order

`NewDataStore` runs a meta-retrieval phase followed by the ordered normalization phases below:

1. **Meta retrieval** (before any loops): `setMetaHost`, `retrieveEnvMetaClassesFromRemote`, `refreshMetaFiles` populate `gen/meta/` from the remote pubhub host on demand.
2. **Loop 1 — Load classes** (`loadClasses` → per file `loadClass` → `NewClass` → `setClassData`): parse each meta file and run the class-level normalization chain. Class loading may recursively load a related class when resource-name derivation needs it. Within a class, parents are resolved before properties so the synthetic `parentDn` property can be added; relation, resource-name, RN, test-configuration, and version data are then completed before the global post-processing loops.
3. **Loop 2 — Documentation:** for each loaded class, call `setDocumentation(ds)`. Runs after Loop 1 globally because doc rendering reads child class info from the DataStore.
4. **Loop 3 — Test dependencies + property test values** (`setTestData`, first loop): per class, call `setTestDependencies(ds)` then `setPropertyTestValues(ds)` back-to-back.
   - `setTestDependencies` resolves the `TestDependencies` DAG (shared pointers via `seen`) and resolves `ConfigOverrides` placeholders against the DAG.
   - `setPropertyTestValues` wires `parent_dn` from Parent dependencies, wires `tDn` or a named target property from Target dependencies, and resolves `{{<reference>}}` placeholders in property `TestValues`.
   - These two share a loop because `setPropertyTestValues` reads only the class's own `TestDependencies` (just populated) — no cross-class read into another class's test data.
5. **Loop 4 — Child test values + completeness validation** (`setTestData`, second loop): per class, call `setChildTestValues(ds)` then `validateTestCompleteness(ctx)`.
   - `setChildTestValues` builds `TestChildren` from child class `TestValues` or the child's `EmbeddedInstances`, applies parent overrides from `ClassDefinition.TestConfig.Children`, collects child-driven dependencies into the parent's `TestDependencies`, and resolves `{{<reference>}}` placeholders in child instance properties.
   - `validateTestCompleteness` reports unresolved placeholders, undefined standard scenarios, and explicitly empty non-Set scenarios.

Loop 4 must follow Loop 3 globally because `buildChildInstance` reads property `TestValues` (including wired `tDn`/`parentDn`) from **other** classes, so every class must have finished Loop 3 first.

---

## 6. When to Provide Manual YAML

| Scenario | What to provide |
|----------|----------------|
| Multi-target relations | `test_config.dependencies` with Role=target for the target instances exercised by tests; at least one explicit Target is required |
| Definition-only parent or target | Explicit `test_config.dependencies`; a name-only definition is never treated as proof of a generated artifact |
| Parent-specific target selection | `target_classes` on that Parent dependency |
| Static DN references | `test_config.dependencies` with `reference_type: "static"` |
| Plain DN/reference property | `property.reference.classes` and, when needed by future examples, `example_class`, `example_attribute`, or `example_label` |
| Correlated child instances reused by all parents | `class.test_config.embedded_instances` on the child class |
| Dependency needing child-driven targets | Manual `dependencies` on the dependency (NOT auto-collected for dependencies, only for resource-under-test) |
| Dependency needing children to be valid | Explicit `children` block on the dependency definition (children are NOT auto-populated for dependency resources, even if the dependency's own class has children defined) |
| Non-standard test values | `test_config` on PropertyDefinition (server normalization, special characters) |
| Server default differs from empty string | `test_config.default` on PropertyDefinition with both `assert_value: "<expected>"` AND explicit `config_include: false`. Omitting `config_include` defaults to `true` and would force the value into HCL config, defeating the point. |
| 3+ parent classes | `test_config.dependencies` retains an additional parent for a future test renderer to emit as an HCL prerequisite, but `parent_dn` auto-wiring still only uses the first 2 Parent-role entries (Create/Update/Default use [0], ForceNew uses [1]). To switch to a 3rd parent type, also override the `parent_dn` property's `test_config` directly with a `{{<reference>}}` placeholder pointing at the desired entry. |
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
                tnFvCtxName:
                  - config_value: "{{aci_vrf.test.id}}"
      dependencies:
        - class_name: fvTenant
          reference: "aci_tenant.test.id"
          reference_type: "resource"
        - class_name: fvCtx
          reference: "aci_vrf.test.id"
          reference_type: "resource"
```

### 7.5 Multi-target override example (fvRsDomAtt)

`fvRsDomAtt` is a multi-target relation: its meta `toMo` is the abstract `infra:DomP`, while `relation_info.to_classes` names the four concrete domain types. Multi-target relations disable target auto-resolution, so the definition declares the target instances represented in normalized test data. The retained `target_classes` value associates the selected `fvAEPg` parent with `vmmDomP`.

```yaml
test_config:
  dependencies:
    - class_name: fvAEPg
      reference: "aci_application_epg.test.id"
      role: "parent"
      target_classes: [vmmDomP]
    - class_name: vmmDomP
      reference: "uni/vmmp-VMware/dom-domain_1"
      reference_type: "static"
      role: "target"
    - class_name: vmmDomP
      reference: "uni/vmmp-VMware/dom-domain_2"
      reference_type: "static"
      role: "target"
```

`tDn` is wired from the first two declared targets:

- Create/Default/ForceNew = `uni/vmmp-VMware/dom-domain_1`
- Update = `uni/vmmp-VMware/dom-domain_2`

Additional domain types can be added as explicit Target dependencies when a test needs to exercise them. The future example normalizer can use `target_classes` to replace these static test fixtures with a semantic Terraform reference without mutating the test lifecycle values.

---

## 8. Custom Test Cases

### 8.1 File naming convention

- **Generated** (overwritten on regen): `resource_aci_<name>_test.go`
- **User-maintained** (never touched): `resource_aci_<name>_custom_test.go`

User-maintained `*_custom_test.go` files must never be added to generated render jobs or the managed-files manifest. Manifest cleanup deletes only exact recorded paths, so those files remain untouched without a filename scan.

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
- **ForceNew / parent-switch step** — apply a config that swaps `parent_dn` to the second Parent-role dependency from `TestValues.ForceNew`; assert destroy+recreate of the resource-under-test. Relation `tDn` remains on its Create target in this step.
- **Import step** — `ImportState`/`ImportStateVerify` after one of the above apply steps; assert imported state matches applied state.
- **Data source step** — apply the resource then read it via `data.aci_<name>`; assert every attribute (including children) matches the resource's state.
- **Cleanup** — final destroy step (typically implicit via the test framework, but the template must not leave orphan dependencies).

### Import testing gaps
- Import with non-default attribute values (legacy generated tests import after their reset/default step)
- Import after child modifications (partial child presence)
- Verify `parent_dn` reconstruction correctness explicitly (legacy generated tests cover it only indirectly through `ImportStateVerify`)

### Child lifecycle gaps
- Child attribute updates (changing a value within an existing child block between steps, not add/remove)
- Incremental child addition (adding children one at a time, not all at once)
- Child ordering stability verification (reorder set-typed children, verify no diff)
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
- Explicit plan-shape assertion that a replace-trigger step shows destroy+create (legacy generated tests infer this only from apply success).
- RequiresReplace coverage strategy:
  - Implicit coverage when a property's Create and Update values differ (Terraform auto-plans replace as part of the update step). Free-form string properties normally satisfy this because §2.4 auto-derives distinct `"<attr>_1"` / `"<attr>_2"` values.
  - **Gap to address:** RequiresReplace properties whose `test_config` override pins Create == Update, or whose `ValidValues` only contain a single member. The (future) template Update step should detect these and either widen the value set or emit an explicit replace-trigger scenario.

### Generator code follow-ups (not template scope)

- **Silent-skip diagnostics for dropped parents:** Both parents encountered after two resolvable parent classes and parents without a resource artifact are dropped with only a Trace-level log. Consider emitting a Warning-level diagnostic in either case so coverage gaps are visible to whoever runs the generator.
- **Multi-target lifecycle data only selects the first two `Target`-role entries via `tDn` Create/Update.** For relations with 3+ concrete targets, further explicit dependencies remain available in `TestDependencies` but are not assigned to a lifecycle bucket. The future acceptance-test renderer must decide whether to emit an additional scenario per remaining target.

### ForceNew / parent_dn
- Multi-parent type compatibility (verify resource works under different parent types within same test)

### Attribute update
- Value-to-value update (All_v1 → All_v2 with different non-default values; legacy generated tests use All → Min → Reset)
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
