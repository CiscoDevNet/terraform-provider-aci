---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_external_network_instance_profile"
sidebar_current: "docs-aci-data-source-aci_external_network_instance_profile"
description: |-
  Data source for ACI External Network Instance Profile
---

# aci_external_network_instance_profile

Data source for ACI External Network Instance Profile

## API Information ##

* `Class` - l3extOut
* `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/instP-{external_epg}

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs -> External EPGs

## Example Usage

```hcl
data "aci_external_network_instance_profile" "external_epg" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  name          = "external_epg"
}
```

## Argument Reference

* `l3_outside_dn` - (Required) Distinguished name of the parent L3Outside object.
* `name` - (Required) Name of the External Network Instance Profile object.

## Attribute Reference

* `id` - Attribute id set to the Dn of the External Network Instance Profile.
* `annotation` - (Optional) Annotation of the External Network Instance Profile object.
* `exception_tag` - (Optional) Exception tag of the External Network Instance Profile object.
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
* `match_t` - (Optional) The provider label match criteria of the External Network Instance Profile object.
* `name_alias` - (Optional) Name alias of the External Network Instance Profile object.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if an External EPG is part of a group that does not require a contract for communication.
* `prio` - (Optional) The QoS priority class identifier of the External Network Instance Profile object.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
* `relation_fv_rs_sec_inherited` - (Optional) Relation to EPGs to be used as Contract Masters (class fvEPg). Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_prov` - (Optional) Relation to Provided Contracts (class vzBrCP). Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_cons_if` - (Optional) Relation to Provided Contract Interfaces (class vzCPIf). Cardinality - N_TO_M. Type - Set of String.
* `relation_l3ext_rs_inst_p_to_profile` - (Optional) Relation to Route Control Profiles (class rtctrlProfile). Cardinality - N_TO_M. Type: Block.
  * tn_rtctrl_profile_dn - (Optional) Distinguished name of the Route map for import and export route control.
  * direction - (Optional) Direction of the Route Control Profile.
* `relation_fv_rs_cons` - (Optional) Relation to Consumed Contracts (class vzBrCP). Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_prot_by` - (Optional) Relation to Taboo Contracts (vzTaboo). Cardinality - N_TO_M. Type - Set of String.
