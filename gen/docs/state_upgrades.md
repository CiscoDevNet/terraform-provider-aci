# State Upgrades

`state_upgrades` records historical schema differences and controls transitional
legacy attributes. The active model and resource templates generate compatibility
models, schemas, conversions, and registered upgraders from these definitions.
The latest metadata and normal property overrides lead the current class model.
Normal generation does not read the preserved SDKv2 schema snapshot.

## 1. Overview

`schema_version` pins the current resource schema version. `state_upgrades` is a
map keyed by saved schema version, not a list of sequential migration steps.
Terraform selects the upgrader matching the saved version; every entry converts
directly to the current model. Version numbers may have gaps and are independent
of provider releases, APIC versions, and other resources.

Record only historical differences. Unchanged root fields such as `name`,
`description`, `annotation`, and `id` need no entries. Historical flags alone do
not create another configurable alias. Datasources expose applicable legacy
outputs but have no resource schema versions or registered state upgraders.

## 2. Single derivation rule

An `attribute_name` records the historical Terraform name. The same declaration
drives prior-state decoding and, when it differs from the current representation,
the current compatibility alias. `status` determines whether that alias remains
deprecated, is frozen, or has been removed. Type and restriction overrides can
describe historical state without adding an alias to the current schema.

## 3. YAML shape

This is the current Tenant history:

```yaml
schema_version: 2
state_upgrades:
  1:
    source: sdkv2
    attributes:
      relation_fv_rs_tn_deny_rule:
        type: set_attribute
        element_type: string_attribute
        restriction: optional
        status: frozen
    children:
      fvRsTenantMonPol:
        attributes:
          tnMonEPGPolName:
            attribute_name: relation_fv_rs_tenant_mon_pol
            status: deprecated
```

Omitting `schema_version` means version zero. The migration script imports the
published version explicitly, including resources that have no state upgrader.
Do not infer the current version from the number or maximum key of the history.

## 4. Attribute nodes

`StateUpgradeDefinition` contains `Source`, `Attributes`, and `Children`.
`AttributeUpgradeDefinition` supplies these fields:

| YAML field | Meaning |
|---|---|
| `attribute_name` | Terraform name in the historical schema; omit to inherit the current name. |
| `type` | Historical framework attribute type; omit to inherit the current type. |
| `value_type` | Explicit historical string semantics, using the existing value-type vocabulary such as `ip_address`. Renamed aliases otherwise remain plain strings. BD's `llAddr` records `ip_address` to preserve its deprecated `ll_addr` custom type without changing the Terraform wire type. |
| `element_type` | Historical element type for a collection. |
| `restriction` | Historical `required`, `optional`, `optional_computed`, or `read_only` flags. |
| `status` | Current exposure: `deprecated`, `frozen`, or `removed`. |
| `attributes` | Historical scalar members of a child block. |
| `children` | Historical child-class differences. |

Keys under `attributes` use metadata property names, or a known synthetic name
such as `parentDn`/`tDn`. Keys under `children` use APIC class names. An orphaned
frozen or removed field needs enough explicit information to reconstruct its
historical type; it cannot inherit from a property that no longer exists.

## 5. Historical child shapes

A `children` entry with an outer `attribute_name` describes a historical block.
List its complete historical membership with explicit names, including members
whose names did not change. This prevents newer child properties from entering
the historical schema. See the static-leaf relationship example in
[model.md](model.md#legacy-attributes-and-state-upgrades).

An entry without an outer historical name can map an old flat scalar into a new
child attribute, as in Tenant's monitoring-policy relationship above. A type-only
entry can describe a historical set-to-singleton shape without exposing a new
legacy block. Generation reports unsupported historical shapes rather than
guessing a conversion. The recursive input type does not imply that every
possible nested restructuring is implemented.

## 6. Status lifecycle

| Status | Current behavior |
|---|---|
| `deprecated` | Expose the old input, map it to the modern value, and reject simultaneous explicit configuration of both forms. |
| `frozen` | Retain deprecated configuration/state without configuring the value on APIC. |
| `removed` | Omit the old input from the current schema/model; retain historical decoding information for state upgrades. |

Versions are resolved numerically. The latest explicit status for a legacy path
wins; omission inherits the previous status. A new alias defaults to
`deprecated`. Omission must not restore a removed field. There is no
`functioning` status and no `legacy_status` input.

When removing an exposed attribute, advance `schema_version` and record its
removal in the history for the saved schema that still contained it. Immediate
removal is also supported. Preserve earlier history instead of erasing the
information needed to decode existing state.

## 7. Source lineage

Each version entry may declare `source: sdkv2` or `source: framework`. This is
version-specific history, not a class-level `migration_source` setting.
SDKv2 entries also supply the resource documentation's migration warning.
A resource requiring only a warning uses its existing documentation warnings;
it does not need an artificial state-upgrade entry.

## 8. Validation

YAML parsing rejects unknown fields and enum values. Normalization checks saved
versions against the declared current version, resolves property/class names,
and reports unsupported mappings through generator diagnostics. Removed scalar
nodes require `attribute_name`, `type`, and `restriction`; collection element
types must also describe a supported historical value. Block membership and
scalar-to-child mappings must contain enough information for decoding.

Do not add unchanged root attributes merely to build a complete-schema
inventory. Prior schemas start from the normalized current view and apply the
historical differences. Missing optional JSON fields decode to null.

## 9. Generated behavior

The model template keeps APIC values separate from Terraform compatibility:

- The root model embeds the class model and the current version's legacy model.
- Historical root models and resource-schema functions retain prior names/types.
- `SetFromLegacyConfig` converts configured legacy inputs during planning.
- `SetLegacyFromResponse` populates compatibility outputs independently of modern
  response decoding; unreadable secrets retain only their own prior/configured value.
- Typed `UpgradeFromV<n>` methods convert each supported saved version directly
  into the current root model. Resource adapters register those methods with
  their corresponding prior schemas.

Named relationship aliases accept historical DNs and extract names for modern
attributes. Legacy reads use APIC's resolved `tDn`. Explicit-DN relationships
retain their full DN. Existing null/empty behavior is part of each conversion.

## 10. Coverage and deferred work

`resource_migration_test.go.tmpl` generates APIC-independent unit tests for the
declared resource histories. They invoke registered upgraders with historical
JSON and compare complete typed state. Current fixtures cover values, null
relations, empty relations, and missing optional fields. Scalar fields and
collection entries use distinct markers; historical child blocks contain two
distinct objects so swapped fields and dropped children are observable. These
synthetic markers exercise state conversion, not APIC configuration validation.

The generated upgrader and expected fixtures share normalized definitions, so
passing generated tests alone cannot establish that those definitions match the
published contract. `TestResourceStateUpgradeHistoricalFixtures` in
`internal/provider/provider_migration_test.go` supplements them with handwritten
input and expected JSON independent of the generator definitions. Its examples
cover scalar/parent aliases, named and explicit-DN relationships, child blocks,
frozen values, list-to-set conversion, NetFlow set-to-singleton conversion, and
the accepted VMM null-preservation rule. Omitted fields decode to typed nulls;
the harness still compares the complete current state. These fixtures are
synthetic historical examples, not captured customer state.

Generated resource acceptance tests separately exercise legacy input forms,
unchanged plans, switching back to current inputs, conflicting forms, and
frozen-field retention/clearing where applicable. Preserved resource/datasource
tests remain unchanged. Running an older provider binary inside an acceptance
step and resource identity imports remain outside this coverage.

Migration-source corrections belong in
`gen/scripts/migrate_class_definitions.go` and must survive rerunning that script.
The preserved legacy definition files remain migration inputs, not active
generator inputs.

## 11. Worked example: removing a frozen attribute

The following hypothetical next Tenant schema removes the frozen deny-rule
input while preserving decoding for the versions that exposed it. The original
version-1 entry remains; the new entry records the version-2 shape saved before
removal:

```yaml
schema_version: 3
state_upgrades:
  1:
    source: sdkv2
    attributes:
      relation_fv_rs_tn_deny_rule:
        type: set_attribute
        element_type: string_attribute
        restriction: optional
        status: frozen
    children:
      fvRsTenantMonPol:
        attributes:
          tnMonEPGPolName:
            attribute_name: relation_fv_rs_tenant_mon_pol
            status: deprecated
  2:
    source: framework
    attributes:
      relation_fv_rs_tn_deny_rule:
        attribute_name: relation_fv_rs_tn_deny_rule
        type: set_attribute
        element_type: string_attribute
        restriction: optional
        status: removed
```

The current schema omits `relation_fv_rs_tn_deny_rule`. Both old versions still
have decoding views, and each upgrader converts directly to version 3. The
monitoring-policy alias retains its earlier deprecated status because no later
entry changes it. This example does not change the repository's Tenant schema.
