---
layout: "aci"
page_title: "ACI: aci_node_inb_mgmt_epg"
sidebar_current: "docs-aci-resource-node_inb_mgmt_epg"
description: |-
  Manages ACI In-Band Management EPg
---

# aci_node_inb_mgmt_epg #
Manages ACI In-Band Management EPg

## Example Usage ##

```hcl
resource "aci_node_inb_mgmt_epg" "example" {

  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"
  annotation  = "example"
  encap  = "example"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"

}
```
## Argument Reference ##
* `management_profile_dn` - (Required) Distinguished name of parent management profile object.
* `name` - (Required) The in-band management endpoint group name. This name can be up to 64 alphanumeric characters.
* `annotation` - (Optional) Annotation for object in-band management EPg.

* `encap` - (Optional) The in-band access encapsulation.  

* `exception_tag` - (Optional) Exception tag for object in-band management EPg.

* `flood_on_encap` - (Optional) Control at EPg level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
Allowed values: "disabled", "enabled". Default value: "disabled".
* `match_t` - (Optional) The provider label match criteria. 
Allowed values: "All", "AtleastOne", "AtmostOne", "None". Default value: "AtleastOne".
* `name_alias` - (Optional) Name alias for object in-band management EPg.

* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
Allowed values: "exclude", "include". Default value: "exclude".
* `prio` - (Optional) The in-band QoS priority class identifier. 
Allowed values: "level1", "level2", "level3", "level4", "level5", "level6", "unspecified".

* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_mgmt_rs_mgmt_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_mgmt_rs_in_b_st_node` - (Optional) Relation to class fabricNode. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the In-Band Management EPg.

## Importing ##

An existing In-Band Management EPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_node_inb_mgmt_epg.example <Dn>
```