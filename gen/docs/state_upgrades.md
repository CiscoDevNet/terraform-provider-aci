# State Upgrades

> **Work in progress.** This document is actively being drafted alongside the wider `gen/docs/` restructure. Content and organization may change.

This document describes the `state_upgrades` definition input used by the ACI provider code generator to drive Terraform's `UpgradeResourceState` RPC and the current-schema exposure of transitional legacy attributes. It covers the YAML shape, the recursive node model, the validation rules, and the `legacy_status` lifecycle.

---

## 1. Overview

`state_upgrades` is a per-class YAML list, declared inside a class definition file (e.g. `gen/definitions/fvAEPg.yaml`), that records every prior schema version a resource has ever had and exactly how each prior version differs from the current schema.

It replaces the legacy `migration_blocks` / `migration_version` / `type_changes` keys that lived in older `gen/definitions/classes/*.yaml` files, along with the schema-git-commit JSON. Bulk conversion of those legacy keys is a separate follow-up effort; this design is the input contract for the new template wiring.

Two pieces of information drive the renderer:

- **`state_upgrades`** — the list of upgrade hops. Each entry describes one prior schema version's shape in enough detail to (a) reconstruct the prior schema for the framework's `UpgradeResourceState` RPC and (b) render the current schema's transitional legacy aliases.
- **`migration_source`** — the lineage of the resource (e.g. `from_sdkv2`). Orthogonal to `state_upgrades`. Drives the docs migration warning and any future migration-source-specific codegen.

---

## 2. Single derivation rule

Each `legacy_attribute` declared anywhere under `state_upgrades` is the only input the renderer needs to derive both:

1. The prior-schema reconstruction required by `UpgradeResourceState`.
2. The current-schema exposure (or non-exposure) of that legacy name as a deprecated alias.

The optional `legacy_status` field on each node picks the lifecycle stage. See §6.

---

## 3. YAML shape

```yaml
migration_source: from_sdkv2          # optional; today only "from_sdkv2"

state_upgrades:
  - prior_schema_version: 0           # required; unique across the list
    attributes:                       # optional; map of meta PropertyName -> node
      pcEnfPref:
        legacy_attribute: pc_enf_pref
      removedAttr:
        legacy_attribute: removed_attr
        legacy_type: string_attribute
        legacy_restriction: optional
        legacy_status: removed
    children:                         # optional; map of meta child class name -> node
      fvRsBd:
        legacy_attribute: relation_to_bridge_domain
      fvRsNodeAtt:
        attributes:                   # scalar-wrap: inner attrs carry prior flat scalars
          tDn:
            legacy_attribute: node_dn
          encap:
            legacy_attribute: node_encap
```

Per the Terraform plugin framework contract, each entry is **direct-to-current** — the framework selects the upgrader matching the saved state's version and does not chain intermediate entries.

---

## 4. Attribute upgrade node model

`attributes` and `children` entries share the same recursive node shape:

```go
type AttributeUpgradeDefinition struct {
    LegacyAttribute   string                                // prior TF attribute name
    LegacyType        LegacyAttributeTypeEnum               // prior framework attribute type
    LegacyRestriction RestrictionEnum                       // prior required/optional/read_only
    LegacyStatus      LegacyStatusEnum                      // current-schema exposure
    Attributes        map[string]AttributeUpgradeDefinition // inner scalar attrs
    Children          map[string]AttributeUpgradeDefinition // inner nested blocks
}
```

Recursive and unbounded in depth: any `attributes` or `children` map can itself describe inner attribute and child changes.

### 4.1 Key namespace

- Keys under `attributes` are **always** meta `PropertyName` values (meta-derived like `pcEnfPref` or synthetic like `parentDn` / `tDn`).
- Keys under `children` are **always** meta child class names.
- A `legacy_status: removed` entry may reference a meta-only property name that no longer appears in the resolved current-schema set.

### 4.2 Field semantics

| Field | Required? | Meaning |
|---|---|---|
| `legacy_attribute` | when status is `removed`; otherwise optional | The Terraform attribute name as it appeared in the prior schema. Omit to carry the current name forward. |
| `legacy_type` | when status is `removed`; otherwise optional | The framework attribute type in the prior schema. Omit to inherit from the current attribute. |
| `legacy_restriction` | when status is `removed`; otherwise optional | The schema restriction (`required` / `optional` / `read_only`) in the prior schema. Omit to inherit from the current attribute. |
| `legacy_status` | optional | Current-schema exposure of the legacy name. See §6. |

---

## 5. Disambiguation tables

### 5.1 `attributes` entries

| Fields present | Meaning |
|---|---|
| `legacy_attribute` | Rename only. |
| `legacy_type` | Type change only. |
| Both | Rename + type change. |
| `legacy_status: removed` + all three legacy fields | Attribute existed in prior schema, no current-schema home. Entry retained to drive migration only. |

### 5.2 `children` entries

| Fields present | Meaning |
|---|---|
| `legacy_attribute` only on the child | Block existed in prior schema under a different name, same type. |
| `legacy_type` only on the child | Block existed under the same name, different type. |
| Both on the child | Existed under a different name AND a different type. |
| Neither on the child, only inner annotations | Block existed under the same name/type but inner properties had renames/type changes. |
| Neither on the child, at least one inner attribute carries `legacy_attribute` | Block is NEW in current schema; inner attribute(s) receive the prior flat scalar value(s) (scalar-wrap). |

---

## 6. `legacy_status` lifecycle

`legacy_status` controls how a legacy attribute is exposed in the **current** schema (separate from how it is consumed during state upgrade).

| Value | Current schema exposure | Device round-trip | Use case |
|---|---|---|---|
| `functioning` *(default; zero value)* | Legacy name exposed alongside the replacement, `Optional+Computed+Deprecated+ConflictsWith` | Yes — value piped into the replacement | Just renamed: both names work today. |
| `frozen` | Legacy name exposed with `Deprecated`, state-preserving plan modifier | No | APIC has dropped support but configs keep the old name without a diff. |
| `removed` | Legacy name NOT exposed | N/A — migration-only | Post schema-version bump: prior state still needs to be readable. |

The zero value is `functioning` because that is the natural meaning of a freshly renamed attribute. Omitting the YAML field reaches the zero value via Go zero-init.

A field present-but-empty (`legacy_status: ""`) is treated as a typo and errors at parse time.

### 6.1 Typical progression

Edited in-place on the YAML by flipping `legacy_status`:

1. Rename introduced → add `state_upgrades` entry with `legacy_attribute: pc_enf_pref` on `pcEnfPref`. Default (`functioning`) — both names exposed.
2. APIC drops support → set `legacy_status: frozen`. Both names still exposed, but old no longer talks to the device.
3. Schema version bump → set `legacy_status: removed` AND bump the resource's `SchemaVersion`. Old name now gone from current schema; entry is migration-only.

---

## 7. `migration_source`

A class-level enum recording the lineage of the resource. Today the only non-zero value is `from_sdkv2`. The enum is extensible: additional sources can be added as plain new constants without changing the field type or any consumer.

Coherence rule: if `migration_source` is set, `state_upgrades` MUST contain at least one entry — that entry (or entries) describes the migration hop from the prior provider. The specific `prior_schema_version` depends on what the prior resource was at when migrated: commonly `0` for an SDKv2 resource that never had a schema bump, but it can be any value if the SDKv2 resource itself went through one or more upgrades before the framework migration.

The converse is allowed: `migration_source` omitted with a `prior_schema_version: 0` entry is a normal framework v0→v1 schema bump on a resource that was born in the framework.

---

## 8. Validation

`Class.setStateUpgrades()` runs the following checks after the YAML decode (typed enums already reject unknown values at parse time):

- Unique `prior_schema_version` across the list.
- `legacy_status: removed` requires `legacy_attribute`, `legacy_type`, AND `legacy_restriction` — cannot inherit from a non-existent current property.
- Duplicate `legacy_attribute` values within the same prior-version entry are rejected (collision detection).
- Scalar-wrap shape: a `children` entry without `legacy_attribute`, `legacy_type`, or any inner annotations is rejected as ambiguous.
- Keys under top-level `attributes` resolve to a current `Property.PropertyName` on the class (unless `legacy_status: removed`).
- Keys under top-level `children` resolve to a known meta child class name on the class (unless `legacy_status: removed`).
- `migration_source` set with an empty `state_upgrades` list is rejected (the migration hop from the prior provider must be described by at least one entry).

---

## 9. How the renderer consumes the data

- `Class.StateUpgrades` carries the full validated tree, used by the upgrader template that emits one `resource.StateUpgrader` per entry.
- `Property.StateUpgradeValues` is a flattened convenience map of `prior_schema_version -> StateUpgradeValue`, populated by `Class.setPropertyStateUpgradeValues()` for each top-level attribute on the owning class. Templates that iterate properties can ask "what was my legacy name at v0?" without traversing the upgrade tree.
- Each `StateUpgradeValue` carries the `legacy_status` lifecycle stage from §6 in its `Status` field. The generator uses this to populate `Property.TestValues.Legacy` for each property with a `Functioning` or `Frozen` renamed alias — see `test_configuration.md` §2.4 "Legacy bucket auto-derivation" for the full rules. `Removed` entries never produce a Legacy bucket because they describe migration-only attributes that no longer exist in the current schema.
- Inner `children[X].attributes` entries are **not** distributed into child-class Property maps in this phase. Templates that need inner-child state-upgrade data walk `Class.StateUpgrades` directly.

---

## 10. Out of scope (in this phase)

- Template wiring: how upgraders, prior schemas, and legacy alias rendering are emitted is a follow-up effort.
- Bulk migration of existing `migration_blocks` / `migration_version` / `type_changes` keys from the older `gen/definitions/classes/*.yaml` files into the new model — those files are not loaded by the current generator.
- Orphan-frozen edge case (legacy attribute with no current replacement, frozen alias only) — deferred until a real case appears.

---

## 11. Worked example: lifecycle in YAML

A fictional `fvBD` definition exercising all three `legacy_status` values in a single `state_upgrades` entry. Illustrates how the same input drives both the `UpgradeResourceState` upgrader and the auto-derived `TestValues.Legacy` bucket described in `test_configuration.md` §2.4.

```yaml
# gen/definitions/fvBD.yaml (excerpt)
migration_source: from_sdkv2

state_upgrades:
  - prior_schema_version: 0
    attributes:
      # 1. Functioning rename: both names exposed in the current schema.
      #    APIC still accepts the value via the new attribute name.
      arpFlood:
        legacy_attribute: arp_flooding

      # 2. Frozen rename: legacy name still accepted by Terraform configs but
      #    no longer round-trips to the device. State-preserving alias only.
      mac:
        legacy_attribute: mac_address
        legacy_status: frozen

      # 3. Removed: prior schema had the attribute, current schema does not.
      #    Migration-only; the legacy fields are required because there is
      #    no current property to inherit them from. Bump the resource
      #    SchemaVersion to 1 in lockstep with this change.
      legacyUnicastRoute:
        legacy_attribute: enable_unicast_routing
        legacy_type: string_attribute
        legacy_restriction: optional
        legacy_status: removed
```

What the generator produces from this single entry:

| Meta property | `Status` | Current schema | Auto-derived `TestValues.Legacy` |
|---|---|---|---|
| `arpFlood` | `Functioning` | Both `arp_flood` (current) and `arp_flooding` (deprecated alias) exposed; `ConflictsWith` enforced. | Clone of `arpFlood`'s `Create` bucket. Template renders one scenario configuring the resource under `arp_flooding` and asserts the value lands in state at the current name. |
| `mac` | `Frozen` | Both names exposed; `mac_address` carries a state-preserving plan modifier so existing configs don't diff. | Clone of `mac`'s `Create` bucket. Template renders a scenario configuring the resource under `mac_address` and asserts no device write occurs. |
| `legacyUnicastRoute` | `Removed` | Not exposed. Resource `SchemaVersion` bumped to 1. | Nil. Coverage comes from the `UpgradeResourceState` migration scenario, not from a current-schema scenario. |

If the `mac` rename had also changed type (for example a `string_attribute` becoming a `set_attribute`), `hasDivergentLegacyType()` would skip auto-derivation, emit a generator `Warn`, and require an explicit `test_config.legacy` block on the `mac` property carrying the prior-shape HCL values.
