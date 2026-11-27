# Redesign change overview

Inventory of known differences and decisions. The final compatibility comparison
on 23 November 2026 used the agreed `develop` baseline at `97bf02b1f` and covered
126 framework resources and 127 datasources, including nested attributes, using
the same dependency versions. The audit is complete; the compatibility findings
below still require review. Change type and compatibility impact are separate:
a bug fix can still require resource replacement.

## User-visible differences

| Area | Type | Difference and impact |
| --- | --- | --- |
| Deprecated datasource lookup names | Enhancement | Datasources now accept deprecated identifier and parent-DN names as lookup inputs until removal. Configure one current or historical name per lookup field; it is normalized before constructing the DN. Other legacy attributes remain read-only. This extends the accepted datasource inputs compared with `develop`. |
| Application EPG `useg_epg` | Bug fix | Mode changes now require replacement after the tested APIC rejected an in-place change. Document the replacement impact. [Definition](../definitions/fvAEPg.yaml). |
| Legacy named-relation block state | Bug fix | Historical Bridge Domain NetFlow, VRF address-family policy, and ESG contract block values are DNs. State upgrading now extracts the target name for the modern child while retaining the original DN in the deprecated block, aligning it with existing legacy planning. The previous upgraders copied the DN into the name field. |
| NetFlow historical child values | Bug fix | The generated upgrader constructs exporter/record child values with the complete current object types. The baseline version-0 upgrader produces value-conversion diagnostics for historical objects lacking the newer nested annotation/tag fields; all four historical-shape probes now upgrade successfully. |
| VMM Domain historical null ARP learning | Migration state normalization; accepted | Upgrading version-1 state now preserves a missing/null `arp_learning` value as null. The old upgrader converted it to `""` through `ValueString()`. Null preservation is intentional; Terraform can distinguish these values before refresh. Normal schema defaults and APIC reads are unchanged. |
| STP legacy control configuration | Bug fix | The deprecated `ctrl` list now supplies `interface_controls` during planning. The previous legacy-input branch read the unset modern field instead, discarding the supplied controls. The modern set removes duplicate controls; the legacy list retains its configured order when it matches APIC, avoiding an inconsistent result for reordered input. |
| Nested child planning | Bug fix | Recursive reconciliation preserves computed child state by semantic identifiers when sets change or reorder. [Plan modifiers](../../internal/planmodifiers/). |
| Undocumented input restriction | Validation tightening; accepted | QoS policing no longer accepts hexadecimal `0xffff` / `0xffffffffffffffff` configuration aliases, and VMM Domain ARP learning no longer accepts `0x0` / `0x1`. Documented readable values and numeric ranges remain accepted. Configurations using these undocumented inputs may fail validation; custom-type response comparisons remain available. |
| Unexpected APIC responses | Robustness enhancement | Defensive handling; no corresponding live APIC failure has been identified. Malformed response structures return descriptive errors. If APIC returns multiple instances of a child whose schema permits only one, warn and decode the first so the resource can continue. Valid responses are unaffected. [Response helpers](../../internal/provider/models/helpers/response.go). |
| Provider configuration diagnostics | Bug fix | Corrects the error message for invalid boolean text supplied through `validate_relation_dn` or `ACI_VAL_REL_DN`: it now names `validate_relation_dn` instead of incorrectly naming `insecure`. Valid configuration, defaults, and accepted inputs are unchanged. |
| Legacy identifier replacement | Behavior correction; review | Fifteen deprecated identifier/parent aliases now have explicit `RequiresReplaceIfConfigured` modifiers. An offline Application Profile probe changing `tenant_dn` shows replacement and the new DN; the baseline requested no replacement and retained the old ID despite changing the parent. This aligns legacy inputs with the modern replacement rules, but its user-visible effect must be reviewed. |

The response policies have synthetic unit coverage; no corresponding live
failure has been identified.

Declared legacy aliases, frozen attributes, and state migrations now generate
compatibility code. Audit completion does not resolve the compatibility decisions
listed below.

## Internal mechanics

The future nested-migration corrections below were verified with temporary class
definitions. The current canonical definitions do not introduce those transitions.

| Area | Type | Summary |
| --- | --- | --- |
| Shared models, schemas, helpers, and registration | Refactor | Removes duplication; code relocation alone is not a user-visible change. |
| Recursive historical child shapes | Bug fix | A historical set containing a singleton child is unwrapped before recursive member conversion. This supports the same set-to-singleton conversion below other children and preserves old names/types until their mappings run. Live EPG/domain/uplink proof exposed the missing recursive path. |
| Deprecated child outputs during updates | Bug fix | Changing a current child property invalidates its omitted computed alias during planning. This prevents a stale deprecated value from causing an inconsistent-result error after APIC returns the updated value. Unchanged values and frozen/secret retention keep their existing handling. Live Annotation/Tenant proof exposed and verified the correction. |
| Frozen aliases with current fields | Bug fix | Frozen aliases may retain their own configured value alongside the current field. They no longer reject that current field or request replacement when a frozen identifier alias changes. This keeps frozen data isolated from APIC configuration. Live nested EPG proof exposed the conflict restriction. |
| Secret retention | Compatibility | A `computed: true` override preserves VMM Credential password retention when configuration later omits it. Initial creation through the modern password also seeds the deprecated `pwd` alias, matching the baseline; subsequent reads retain the alias's own configured/state value. Imports cannot recover the secret from APIC. Key Ring and peering passwords remain optional without computed behavior. Never-configured or unknown secrets remain null when APIC cannot return them. |
| Datasource sensitivity flags | Schema metadata alignment | Shared schemas now mark VMM Credential `password` and its deprecated `pwd` alias, Key Ring `key`, and peering `password` as sensitive in datasources too. These four datasource attributes remain null; their returned values are unchanged. |
| Replacement inference | Compatibility | Metadata `isNaming` or `createOnly` restores the four existing replacement rules for Imported Logical Device description, IP SLA Track List type, and the two domain delimiter attributes. Certificate Authority `certificate_usage` uses `requires_replace: false` to preserve its legacy set behavior pending separate validation. |
| HTTPS DH parameter | Compatibility | `commHttps.dhParam` uses `deprecated: false` to preserve the tested API/state behavior without a new provider deprecation warning. APIC's deprecation metadata remains intact; version-aware warnings are deferred. |
| Legacy deprecation messages | Wording | Existing deprecated attributes keep their deprecated status. Messages now identify the replacement directly without promising removal in the next major provider version; frozen attributes retain their separate guidance. |
| Templates and generated-file ownership | Generator workflow | Canonical definitions, manifest-managed outputs, and explicit unsupported-annotation refresh. |
| Acceptance execution | Test enhancement | Isolated workers, fixture prefixes, serial shared-setting tests, and individual logs/reports. |
| Migration regression coverage | Test enhancement | Generated fixtures distinguish fields and collection entries, and exercise two historical block objects. Handwritten historical input/expected JSON independently checks representative conversion rules through the registered upgraders. These checks require no APIC and do not change provider behavior. |
| Existing-object flag coverage | Test enhancement | After minimum creation, reject a duplicate Terraform address with `allow_existing_on_create = false`, adopt it with `true`, then forget only that address while preserving the original resource and APIC object. Direct coverage of the final `Create` rejection guard remains out of scope. |
| Assertion strategies | Test enhancement | Per-entry `value` or `presence_only`; Fabric's default site-ID check permits an MSO-assigned value. Generated datasource checks for configured named relationships' deprecated scalar DNs use native `NotNull()` assertions with APIC-version gates, accepting asynchronous target-DN resolution while retaining exact modern-name checks. Lifecycle operations are unchanged. |
| Fixtures and scenario ordering | Test fixes | Corrected inputs/dependencies, fixed create-only values, and import after update/reset. |

RN/DN/ID construction, certificate-wrapper POST routing, annotation defaults,
and empty-string-to-`none` handling remain compatibility checks, not assumed
behavior changes.

## Final compatibility audit

The registry, current Terraform attribute names/types, resource schema flags,
defaults, schema versions, and registered upgrader keys match the pinned baseline.
All 76 generated conversion expectations pass. For a direct comparison, the same
fixtures were restricted to fields present in the baseline historical schemas:
65 conversions match, six differ only through the documented named-relation
DN-to-name fixes, one differs for VMM ARP-learning null handling, and four NetFlow
conversions now succeed where the baseline returns diagnostics.

The two compatibility findings have been resolved as follows:

- **Bridge Domain `ll_addr` retains IP semantic equality.** An explicit historical
  `value_type: ip_address` preserves the custom type for the deprecated resource
  and datasource alias, its model, and state conversions. Valid equivalent IPv6
  spellings continue comparing equal. Other aliases keep their existing types.
- **VMM Domain version-1 upgrade changes a missing/null `arp_learning` value.**
  The baseline converted it to a known empty string; the generated upgrader
  preserves null. This affects the upgraded state before refresh and is accepted
  as an intentional normalization correction. No null-to-empty setting is added.

The audit also confirms that the generated EPG lifecycle keeps `useg_epg = "no"`
throughout its scenarios. Its replacement modifier is present, but a dedicated
mode-change check is still absent. This is a coverage gap, not an observed live
failure. The legacy replacement behavior described above also needs explicit
compatibility review. No provider changes were made as part of the audit.

## TODO

- [ ] **Reassess version-aware APIC deprecation handling, starting with HTTPS DH parameters.**
  - Controller-hosted `commHttps.dhParam` metadata on APIC **6.2(2e)** records
    `deprecatedSince: "6.0(2)"`; the older checked-in metadata lacks that field.
    It still marks the property configurable, read/write, and not hidden.
  - Keep `deprecated: false` in the definition for now. Resource tests configure
    `1024` and update to `2048`; datasource tests read `1024`. Both passed on
    APIC **6.1(5e)** and **6.2(2e)**. These are API/state checks, not TLS negotiation checks.
  - Refresh metadata and define how `deprecatedSince`, `isHidden`, and
    `hiddenSince` affect warnings, documentation, and test gates. APIC deprecation
    must not automatically promise removal from a future provider release.
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
  - Refresh the metadata, then reassess supported-property version gates using
    `versions`, `deprecatedSince`, `isHidden`, and `hiddenSince`.
  - Regenerate and rerun the resource and datasource tests across the relevant
    APIC versions. Keep these failures visible until the metadata and version
    handling have been reassessed.
- [x] **Run the compatibility audit after migration work.** Compared the pinned
  baseline and current implementation using identical dependencies, complete
  schemas, semantic-equivalence probes, historical-state fixtures, and focused
  offline planning probes. Findings and limits are recorded above.
- [x] **Resolve the two compatibility findings above.** Restore deprecated BD
  `ll_addr` IP semantics and retain/document VMM Domain null preservation.
- [ ] **Review legacy identifier replacement** and add dedicated coverage for
  changing EPG `useg_epg`, without treating unrelated identifier replacement as
  proof of this specific mode-change rule.
