# Redesign change overview

Inventory of known differences and decisions. The final compatibility comparison
on 23 November 2026 used the agreed `develop` baseline at `97bf02b1f` and covered
126 framework resources and 127 datasources, including nested attributes, using
the same dependency versions. Its results and the resolved findings are recorded
below; the remaining follow-ups are listed under TODO. Change type and
compatibility impact are separate: a bug fix can still require resource replacement.

## User-visible differences

| Area | Type | Difference and impact |
| --- | --- | --- |
| Deprecated datasource lookup names | Enhancement | Datasources now accept deprecated identifier and parent-DN names as lookup inputs until removal. Configure one current or historical name per lookup field; it is normalized before constructing the DN. Other legacy attributes remain read-only. This extends the accepted datasource inputs compared with `develop`. |
| Application EPG `useg_epg` | Bug fix | Mode changes now require replacement after the tested APIC rejected an in-place change. Changing between regular and microsegmented mode therefore recreates the EPG. Dedicated acceptance coverage of the `"no"` → `"yes"` replacement is deferred; it does not block this refactor. [Definition](../definitions/fvAEPg.yaml). |
| Legacy named-relation block state | Bug fix | Historical Bridge Domain NetFlow, VRF address-family policy, and ESG contract block values are DNs. State upgrading now extracts the target name for the modern child while retaining the original DN in the deprecated block, aligning it with existing legacy planning. The previous upgraders copied the DN into the name field. |
| NetFlow historical child values | Bug fix | The generated upgrader constructs exporter/record child values with the complete current object types. The baseline version-0 upgrader produces value-conversion diagnostics for historical objects lacking the newer nested annotation/tag fields; all four historical-shape probes now upgrade successfully. |
| VMM Domain historical null ARP learning | Migration state normalization; accepted | Upgrading version-1 state now preserves a missing/null `arp_learning` value as null. The old upgrader converted it to `""` through `ValueString()`. Null preservation is intentional; Terraform can distinguish these values before refresh. Normal schema defaults and APIC reads are unchanged. |
| STP legacy control configuration | Bug fix | The deprecated `ctrl` list now supplies `interface_controls` during planning. The previous legacy-input branch read the unset modern field instead, discarding the supplied controls. The modern set removes duplicate controls; the legacy list retains its configured order when it matches APIC, avoiding an inconsistent result for reordered input. |
| Nested child planning | Bug fix | Recursive reconciliation preserves computed child state by semantic identifiers when sets change or reorder. [Plan modifiers](../../internal/planmodifiers/). |
| Undocumented input restriction | Validation tightening; accepted | QoS policing no longer accepts hexadecimal `0xffff` / `0xffffffffffffffff` configuration aliases, and VMM Domain ARP learning no longer accepts `0x0` / `0x1`. Documented readable values and numeric ranges remain accepted. Configurations using these undocumented inputs may fail validation; custom-type response comparisons remain available. |
| Unexpected APIC responses | Robustness enhancement | Defensive handling; no corresponding live APIC failure has been identified. Malformed response structures return descriptive errors. If APIC returns multiple instances of a child whose schema permits only one, warn and decode the first so the resource can continue. Valid responses are unaffected. [Response helpers](../../internal/provider/models/helpers/response.go). |
| Provider configuration diagnostics | Bug fix | Corrects the error message for invalid boolean text supplied through `validate_relation_dn` or `ACI_VAL_REL_DN`: it now names `validate_relation_dn` instead of incorrectly naming `insecure`. Valid configuration, defaults, and accepted inputs are unchanged. |
| Legacy identifier replacement | Bug fix | Changing one of fifteen deprecated identifier/parent aliases now requires replacement, matching the corresponding current attribute. An offline Application Profile probe changing `tenant_dn` shows replacement and the new DN; the baseline requested no replacement and retained the old ID despite changing the parent. Existing configurations using these aliases can therefore show a replacement when their identifier or parent changes. |

The response policies have synthetic unit coverage; no corresponding live
failure has been identified.

Declared legacy aliases, frozen attributes, and state migrations now generate
compatibility code.

## Internal mechanics

The future nested-migration corrections below were verified with temporary class
definitions. The current canonical definitions do not introduce those transitions.

| Area | Type | Summary |
| --- | --- | --- |
| Shared models, schemas, helpers, and registration | Refactor | Removes duplication; code relocation alone is not a user-visible change. |
| Custom types and APIC versions | Refactor and robustness | Replaces per-property semantic string implementations with configured types in `internal/customtypes`, while retaining distinct IP-address and VMM ARP-learning behavior. IP/CIDR parsing checks errors before accessing the network mask, avoiding a nil-pointer panic on malformed text. Shared child schemas also align datasource custom types with their resource counterparts. Class models register their support ranges; shared `internal/apicversion` comparisons drive subtree filtering and test gates. |
| Definitions and input validation | Generator workflow | Canonical per-class definitions replace the legacy classes/properties layout for normal generation. The migration script reproduces those definitions; strict loading and normalization reject conflicting aliases, invalid historical shapes, and incompatible nested histories before rendering. Reference-only definitions preserve naming/link metadata without creating artifacts when class metadata is absent. |
| Versioned state upgrades | Generator capability | Explicit resource schema versions register direct conversions from each historical version to the current model. Definitions record historical differences and deprecated, frozen, or removed fields; current metadata remains the source for unchanged attributes. Recursive child histories reuse class-owned mappings and keep historical shapes separate from current compatibility fields. Synthetic fixtures cover future renames, type changes, removal, and multiple versions without introducing those transitions into production definitions. |
| Examples and documentation | Generator workflow | Templates consume normalized class/property data for minimum/full resource examples, datasource examples, and resource/datasource pages. Examples use consistent labels, typed references, shared supporting files, and no duplicate full example when it matches the minimum. Documentation uses current schema names/types, parent/child resource links, defaults, version support, and compatibility guidance. |
| Recursive historical child shapes | Bug fix | A historical set containing a singleton child is unwrapped before recursive member conversion. This supports the same set-to-singleton conversion below other children and preserves old names/types until their mappings run. Live EPG/domain/uplink proof exposed the missing recursive path. |
| Deprecated child outputs during updates | Bug fix | Changing a current child property invalidates its omitted computed alias during planning. This prevents a stale deprecated value from causing an inconsistent-result error after APIC returns the updated value. Unchanged values and frozen/secret retention keep their existing handling. Live Annotation/Tenant proof exposed and verified the correction. |
| Frozen aliases with current fields | Bug fix | Frozen aliases may retain their own configured value alongside the current field. They no longer reject that current field or request replacement when a frozen identifier alias changes. This keeps frozen data isolated from APIC configuration. Live nested EPG proof exposed the conflict restriction. |
| Secret retention | Compatibility | A `computed: true` override preserves VMM Credential password retention when configuration later omits it. Initial creation through the modern password also seeds the deprecated `pwd` alias, matching the baseline; subsequent reads retain the alias's own configured/state value. Imports cannot recover the secret from APIC. Key Ring and peering passwords remain optional without computed behavior. Never-configured or unknown secrets remain null when APIC cannot return them. |
| Datasource sensitivity flags | Schema metadata alignment | Shared schemas now mark VMM Credential `password` and its deprecated `pwd` alias, Key Ring `key`, and peering `password` as sensitive in datasources too. These four datasource attributes remain null; their returned values are unchanged. |
| Replacement inference | Compatibility | Metadata `isNaming` or `createOnly` restores the four existing replacement rules for Imported Logical Device description, IP SLA Track List type, and the two domain delimiter attributes. Certificate Authority `certificate_usage` uses `requires_replace: false` to preserve its legacy set behavior pending separate validation. |
| Historical IP semantic equality | Compatibility | An explicit historical `value_type: ip_address` preserves IP semantic equality for Bridge Domain's deprecated `ll_addr` resource/datasource alias and its state conversions. Equivalent IPv6 spellings continue comparing equal. |
| Property lifecycle metadata | Generator behavior | APIC deprecation feeds documentation without enabling Terraform schema deprecation. Known hidden-version ranges narrow generated test configuration and assertions, including exact cutoff releases and backported intervals. Provider deprecation stays explicit in definitions or compatibility history and is documented for current properties, including computed outputs; the temporary HTTPS DH parameter suppression is unnecessary. |
| Property lifecycle validation | Verification | Temporary `fvEpIpTag` metadata validation on 27 November 2026 against APIC 6.2(2e) at `10.15.39.246` passed both generated resource and datasource tests. Documentation limited `name` and `nameAlias` to releases before `6.2(1)`; none of the 22 IP Tag POST payloads contained either property. Original metadata was restored afterward. MAC Tag was not included; permanent refresh remains open. |
| Legacy deprecation messages | Wording | Existing deprecated attributes keep their deprecated status. Messages now identify the replacement directly without promising removal in the next major provider version; frozen attributes retain their separate guidance. |
| Templates and generated-file ownership | Generator workflow | Canonical definitions, manifest-managed outputs, and explicit unsupported-annotation refresh. |
| Generation formatting and dependencies | Generator workflow | Go and Terraform outputs are formatted during rendering. The redundant repository-wide `goimports` pass is removed; active templates declare their imports. Module/vendor changes support the new generator and test helpers. |
| Acceptance execution | Test enhancement | Isolated workers, fixture prefixes, serial shared-setting tests, and individual logs/reports. |
| Migration regression coverage | Test enhancement | Generated fixtures distinguish fields and collection entries, and exercise two historical block objects. Handwritten historical input/expected JSON independently checks representative conversion rules through the registered upgraders. These checks require no APIC and do not change provider behavior. |
| Compatibility validation | Verification | Compared the pinned `develop` baseline and redesigned provider with identical dependencies, schema inventories, semantic-equivalence probes, historical-state fixtures, and focused offline planning probes. See the [validation summary](#compatibility-validation); intentional differences are recorded above. |
| Existing-object flag coverage | Test enhancement | After minimum creation, reject a duplicate Terraform address with `allow_existing_on_create = false`, adopt it with `true`, then forget only that address while preserving the original resource and APIC object. Direct coverage of the final `Create` rejection guard remains out of scope. |
| Assertion strategies | Test enhancement | Per-entry `value` or `presence_only`; Fabric's default site-ID check permits an MSO-assigned value. Generated datasource checks for configured named relationships' deprecated scalar DNs use native `NotNull()` assertions with APIC-version gates, accepting asynchronous target-DN resolution while retaining exact modern-name checks. Lifecycle operations are unchanged. |
| Fixtures and scenario ordering | Test fixes | Corrected inputs/dependencies, fixed create-only values, and import after update/reset. |

RN/DN/ID construction, certificate-wrapper POST routing, annotation defaults,
and empty-string-to-`none` handling remain compatibility checks, not assumed
behavior changes.

## Compatibility validation

Compatibility was checked against `develop` at `97bf02b1f` using schema
comparisons, historical-state conversion fixtures, and offline planning checks.
Intentional behavior differences and deferred test coverage are documented in
the change tables above.

## TODO

- [ ] **Standardize secret omission behavior as an explicit compatibility decision.**
  VMM Credential retains its legacy Optional+Computed password behavior through
  a definition override. Aligning it with optional-only secrets would discard
  retained state values when configuration omits the password and can break
  dependent references. Preserve the distinction until a breaking change is agreed.
- [ ] **Endpoint IP/MAC Tag failures: refresh metadata and reassess property versions.**
  Applies to `fvEpIpTag` / `aci_endpoint_tag_ip` and `fvEpMacTag` /
  `aci_endpoint_tag_mac`.
  - APIC **6.2(2e)**: resource minimum-state checks fail because reads omit
    `name` and `nameAlias` (`name_alias` in Terraform). Datasource tests fail
    earlier, while creating their resource fixtures: APIC rejects writing
    `name` with HTTP 400, error 170,
    `Admin implicit props cannot be modified, prop=name.` These failures occur
    before the datasource lookup assertions.
  - Permanent metadata refresh remains separate. Property `versions` and known
    `hiddenSince` ranges now drive generated availability; `deprecatedSince`
    describes deprecation without ending support.
  - Refresh metadata permanently and rerun across the relevant APIC versions.
    With the original checked-in metadata, the known failures remain visible.
- [ ] **Clean up legacy code and generator artifacts after migration validation.**
  Review and remove the old class/property definitions under
  `gen/definitions/classes/` and `gen/definitions/properties/`, obsolete
  `gen/testvars/`, reference templates under `gen/templates/legacy_templates/`,
  the definition migration script and its migration-only documentation, obsolete
  legacy documentation, and superseded legacy tests. Keep active canonical
  definitions, current provider documentation, and generated regression tests.
- [ ] **Replace `gen/docs/` with a consolidated development and code guide.**
  Cover the generator/data model, definitions and overrides, templates and output
  ownership, provider/model structure, state upgrades, testing, and normal
  development workflows. Update links and retire superseded generator documents
  when the replacement guide is ready.
- [ ] **Investigate generated-test compile time, memory usage, and CI changes.**
  Measure compilation costs and assess reducing repeated generated test code
  without reducing coverage. Evaluate suitable compiler settings, compiling once
  for multiple APIC runs, bounded test parallelism, and CI publication of summaries
  and log artifacts before changing the CI workflows.
- [ ] **Investigate authentication failures above two concurrent acceptance tests.**
  Higher parallelism produced intermittent APIC login failures, including `EOF`,
  in the shared test environment. Investigate `aci-go-client` client/session reuse
  and APIC authentication limits; URL-pointer cache comparisons can recreate
  clients and cause extra logins, but the complete failure cause is not confirmed.
  Retain the tested parallelism of two while investigating. Consider moving client
  ownership into the provider as part of removing the client dependency. See the
  [authentication limitation](test_configuration.md#provider-authentication-limitation).
