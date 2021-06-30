---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-resource-endpoint_security_group"
description: |-
  Manages ACI Endpoint Security Group
---

# aci_endpoint_security_group

Manages ACI Endpoint Security Group

## API Information

- `Class` - fvESg
- `Distinguished Named` - uni/tn-{name}/ap-{name}/esg-{name}

## GUI Information

- `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups

## Example Usage

```hcl
resource "aci_endpoint_security_group" "example" {
  application_profile_dn  = aci_application_profile.example.id
  name  = "example"
  description = "from terraform"
  annotation = "orchestrator:terraform"
  name_alias = "example"
  flood_on_encap = "disabled"
  match_t = "AtleastOne"
  pc_enf_pref = "unenforced"
  pref_gr_memb = "exclude"
  prio = "unspecified"

  relation_fv_rs_cons {
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  relation_fv_rs_cons_if {
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  relation_fv_rs_cust_qos_pol = aci_resource.example.id

  relation_fv_rs_intra_epg = [aci_resource.example.id]

  relation_fv_rs_prov {
    match_t = "AtleastOne"
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  relation_fv_rs_scope = aci_resource.example.id

  relation_fv_rs_sec_inherited = [aci_resource.example.id]
}
```

## Argument Reference

- `application_profile_dn` - (Required) Distinguished name of parent Application Profile object.
- `name` - (Required) Name of object Endpoint Security Group.
- `annotation` - (Optional) Annotation of object Endpoint Security Group.
- `description` - (Optional) Description of object Endpoint Security Group.
- `name_alias` - (Optional) Name alias of object Endpoint Security Group.
- `flood_on_encap` - (Optional) Handles L2 Multicast/Broadcast and Link-Layer traffic at EPG level. It represents Control at EPG level and decides if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP, or based on bridge-domain settings. Allowed values are "disabled", "enabled", and default value is "disabled". Type: String.
- `match_t` - (Optional) Provider Label Match Criteria. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.
- `pc_enf_pref` - (Optional) The preferred policy control enforcement. Allowed values are "enforced", "unenforced", and default value is "unenforced". Type: String.
- `pref_gr_memb` - (Optional) Preferred Group Member parameter is used to determine
  if EPg is part of a group that allows
  a contract for communication. Allowed values are "exclude", "include", and default value is "exclude". Type: String.
- `prio` - (Optional) QoS priority class identifier. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.

- `relation_fv_rs_cons` - (Optional) A block representing the relation to a Contract Consumer (class vzBrCP). The Consumer contract profile information. Type: Block.

  - `prio` - (Optional) The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  - `target_dn` - (Required) The distinguished name of the target contract. Type: String.

- `relation_fv_rs_cons_if` - (Optional) A block representing the relation to a Contract Interface (class vzCPIf). It is a contract for which the EPG will be a consumer. Type: Block.

  - `prio` - (Optional) The contract interface priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  - `target_dn` - (Required) The distinguished name of the target contract. Type: String.

- `relation_fv_rs_cust_qos_pol` - (Optional) Represents the relation to a Custom QOS Policy (class qosCustomPol). It is a source relation to a custom QoS policy that enables different levels of service to be assigned to network traffic, including specifications for the Differentiated Services Code Point (DSCP) value(s) and the 802.1p Dot1p priority. Type: String.

- `relation_fv_rs_intra_epg` - (Optional) Represents the relation to an Intra EPg Contract (class vzBrCP). Represents that the EPg is moving from "allow all within epg" mode to a "deny all within epg" mode. The only type of traffic allowed between EPs in this EPg is the one specified by contracts EPg associates to this relation. Type: List.

- `relation_fv_rs_prov` - (Optional) A block representing the relation to a Contract Provider (class vzBrCP). It is a contract for which the EPG will be a provider. Type: Block.

  - `match_t` - (Optional) The matched EPG type. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.
  - `prio` - (Optional) The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  - `target_dn` - (Required) The distinguished name of the target contract. Type: String.

- `relation_fv_rs_scope` - (Optional) Represents the relation to a Private Network (class fvCtx). Type: String.

- `relation_fv_rs_sec_inherited` - (Optional) Represents the relation to a Security inheritance (class fvEPg). It represents that the EPg is inheriting security configuration from another EPg. Type: List.

## Importing

An existing EndpointSecurityGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import endpoint_security_group.example <Dn>
```
