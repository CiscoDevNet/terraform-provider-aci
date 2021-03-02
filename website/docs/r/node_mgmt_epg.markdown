---
layout: "aci"
page_title: "ACI: aci_node_mgmt_epg"
sidebar_current: "docs-aci-resource-node_mgmt_epg"
description: |-
  Manages ACI Node Management EPg
---

# aci_node_mgmt_epg

Manages ACI Node Management EPg

## Example Usage

```hcl

resource "aci_node_mgmt_epg" "in_band_example" {
  type = "in_band"
  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"
  annotation  = "example"
  encap  = "vlan-1"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
}

resource "aci_node_mgmt_epg" "out_of_band_example" {
  type = "out_of_band"
  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  prio = "level1"
}

```

## Argument Reference

- `type` - (Required) Type of node management EPg to be configured.  
  Allowed values: "in_band", "out_of_band".
- `management_profile_dn` - (Required) Distinguished name of parent management profile object.

### `type = "in_band"`

- `name` - (Required) The in-band management endpoint group name. This name can be up to 64 alphanumeric characters.
- `annotation` - (Optional) Annotation for object in-band management EPg.
- `encap` - (Optional) The in-band access encapsulation.
- `exception_tag` - (Optional) Exception tag for object in-band management EPg.
- `flood_on_encap` - (Optional) Control at EPg level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
  Allowed values: "disabled", "enabled". Default value: "disabled".
- `match_t` - (Optional) The provider label match criteria.
  Allowed values: "All", "AtleastOne", "AtmostOne", "None". Default value: "AtleastOne".
- `name_alias` - (Optional) Name alias for object in-band management EPg.

- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
  Allowed values: "exclude", "include". Default value: "exclude".
- `prio` - (Optional) The QoS priority class identifier.
  Allowed values: "level1", "level2", "level3", "level4", "level5", "level6", "unspecified". Default value: "unspecified".

- `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
- `relation_mgmt_rs_mgmt_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
- `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
- `relation_mgmt_rs_in_b_st_node` - (Optional) Relation to class fabricNode. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.

### `type = "out_of_band"`

- `name` - (Required) The out-of-band management endpoint group name. This name can be up to 64 alphanumeric characters.
- `annotation` - (Optional) Annotation for object out-of-band management EPg.

- `name_alias` - (Optional) Name alias for object out-of-band management EPg.

- `prio` - (Optional) The QoS priority class identifier.
  Allowed values: "level1", "level2", "level3", "level4", "level5", "level6", "unspecified". Default value: "unspecified".

- `relation_mgmt_rs_oo_b_prov` - (Optional) Relation to class vzOOBBrCP. Cardinality - N_TO_M. Type - Set of String.
- `relation_mgmt_rs_oo_b_st_node` - (Optional) Relation to class fabricNode. Cardinality - N_TO_M. Type - Set of String.
- `relation_mgmt_rs_oo_b_ctx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.

**It is recommended to use arguments corresponding to the `type` value. Any invalid argument for applied node management EPg type will be discarded.**

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Node Management EPg.

## Importing

An existing Node Management EPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_node_mgmt_epg.example <Dn>
```
