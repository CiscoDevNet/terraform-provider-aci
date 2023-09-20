---
subcategory: "L2Out"
layout: "aci"
page_title: "ACI: aci_l2out_extepg"
sidebar_current: "docs-aci-resource-l2out_extepg"
description: |-
  Manages ACI ACI L2-Out External EPG
---

# aci_l2out_extepg

Manages ACI L2-Out External EPG

## Example Usage

```hcl
resource "aci_l2out_extepg" "example" {
  l2_outside_dn  = aci_l2_outside.example.id
  description = "from terraform"
  name  = "l2out_extepg_1"
  annotation  = "l2out_extepg_tag"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
  target_dscp = "AF11"
}
```

## Argument Reference

- `l2_outside_dn` - (Required) Distinguished name of parent L2 Outside object.
- `name` - (Required) The name of the layer 2 L2 Out External EPG. This name can be up to 64 alphanumeric characters. Note that you cannot change this name after the object has been saved.
- `annotation` - (Optional) Annotation for object L2 Out External EPG.
- `description` - (Optional) Description for object L2 Out External EPG.
- `exception_tag` - (Optional) Exception tag for object L2 Out External EPG.
- `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link-Local Layer should be flooded only on ENCAP or based on bridge-domain settings.  
  Allowed values: "disabled", "enabled". Default value is "disabled".
- `match_t` - (Optional) The provider label match criteria.  
  Allowed values: "All", "AtleastOne", "AtmostOne", "None". Default value is "AtleastOne".
- `name_alias` - (Optional) Name alias for object L2 Out External EPG. 
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPG is part of a group that does not a contract for communication. Allowed values: "exclude", "include". Default value is "exclude".
- `prio` - (Optional) The QoS priority class identifier.  
  Allowed values: "level1", "level2", "level3", "level4", "level5", "level6", "unspecified". Default value is "unspecified".
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.  
  Allowed values: "AF11", "AF12", "AF13", "AF21", "AF22", "AF23", "AF31", "AF32", "AF33", "AF41", "AF42", "AF43", "CS0", "CS1", "CS2", "CS3", "CS4", "CS5", "CS6", "CS7", "EF", "VA", "unspecified". Default value is "unspecified".
- `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
- `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
- `relation_l2ext_rs_l2_inst_p_to_dom_p` - (Optional) Relation to class l2extDomP. Cardinality - N_TO_ONE. Type - String.
- `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the L2-Out External EPG.

## Importing

An existing L2-Out External EPG can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l2out_extepg.example <Dn>
```
