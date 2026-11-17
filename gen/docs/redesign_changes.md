# Redesign change overview

Inventory of known differences and decisions. The current-schema comparison
against local `develop` at `97bf02b1f` covered 126 framework resources and 127
datasources, including nested attributes, using the same dependency versions.
The final comparison after legacy migration remains open. Change type and
compatibility impact are separate: a bug fix can still require resource replacement.

## User-visible differences

| Area | Type | Difference and impact |
| --- | --- | --- |
| Application EPG `useg_epg` | Bug fix | Mode changes now require replacement after the tested APIC rejected an in-place change. Document the replacement impact. [Definition](../definitions/fvAEPg.yaml). |
| Nested child planning | Bug fix | Recursive reconciliation preserves computed child state by semantic identifiers when sets change or reorder. [Plan modifiers](../../internal/planmodifiers/). |
| Undocumented input restriction | Validation tightening; accepted | QoS policing no longer accepts hexadecimal `0xffff` / `0xffffffffffffffff` configuration aliases, and VMM Domain ARP learning no longer accepts `0x0` / `0x1`. Documented readable values and numeric ranges remain accepted. Configurations using these undocumented inputs may fail validation; custom-type response comparisons remain available. |
| Unexpected APIC responses | Robustness enhancement | Defensive handling; no corresponding live APIC failure has been identified. Malformed response structures return descriptive errors. If APIC returns multiple instances of a child whose schema permits only one, warn and decode the first so the resource can continue. Valid responses are unaffected. [Response helpers](../../internal/provider/models/helpers/response.go). |
| Provider configuration diagnostics | Bug fix | Corrects the error message for invalid boolean text supplied through `validate_relation_dn` or `ACI_VAL_REL_DN`: it now names `validate_relation_dn` instead of incorrectly naming `insecure`. Valid configuration, defaults, and accepted inputs are unchanged. |

The response policies have synthetic unit coverage; no corresponding live
failure has been identified.

Deprecated aliases, frozen attributes, and state migration are unfinished
compatibility work, not approved permanent removals.

## Internal mechanics

| Area | Type | Summary |
| --- | --- | --- |
| Shared models, schemas, helpers, and registration | Refactor | Removes duplication; code relocation alone is not a user-visible change. |
| Secret retention | Compatibility | A `computed: true` override preserves VMM Credential password retention when configuration later omits it. Key Ring and peering passwords remain optional without computed behavior. Never-configured or unknown secrets remain null when APIC cannot return them. |
| Datasource sensitivity flags | Schema metadata alignment | Shared schemas now mark VMM Credential password, Key Ring key, and peering password as sensitive in datasources too. These datasource attributes remain null; their returned values are unchanged. |
| Replacement inference | Compatibility | Metadata `isNaming` or `createOnly` restores the four existing replacement rules for Imported Logical Device description, IP SLA Track List type, and the two domain delimiter attributes. Certificate Authority `certificate_usage` uses `requires_replace: false` to preserve its legacy set behavior pending separate validation. |
| HTTPS DH parameter | Compatibility | `commHttps.dhParam` uses `deprecated: false` to preserve the tested API/state behavior without a new provider deprecation warning. APIC's deprecation metadata remains intact; version-aware warnings are deferred. |
| Templates and generated-file ownership | Generator workflow | Canonical definitions, manifest-managed outputs, and explicit unsupported-annotation refresh. |
| Acceptance execution | Test enhancement | Isolated workers, fixture prefixes, serial shared-setting tests, and individual logs/reports. |
| Existing-object flag coverage | Test enhancement | After minimum creation, reject a duplicate Terraform address with `allow_existing_on_create = false`, adopt it with `true`, then forget only that address while preserving the original resource and APIC object. Direct coverage of the final `Create` rejection guard remains out of scope. |
| Assertion strategies | Test enhancement | Per-entry `value` or `presence_only`; Fabric's default site-ID check permits an MSO-assigned value. Lifecycle operations are unchanged. |
| Fixtures and scenario ordering | Test fixes | Corrected inputs/dependencies, fixed create-only values, and import after update/reset. |

RN/DN/ID construction, certificate-wrapper POST routing, annotation defaults,
and empty-string-to-`none` handling remain compatibility checks, not assumed
behavior changes.

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
- [ ] **Complete the compatibility audit after migration work.** Compare the
  final provider against the agreed baseline, including schema flags, defaults,
  validators, sensitive omission/import behavior, deprecated/frozen attributes,
  and state upgrades. Confirm dedicated coverage for changing `useg_epg` and
  classify any remaining user-visible differences.
