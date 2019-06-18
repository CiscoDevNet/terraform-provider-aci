---
layout: "aci"
page_title: "ACI: aci_external_network_instance_profile"
sidebar_current: "docs-aci-resource-external_network_instance_profile"
description: |-
  Manages ACI External Network Instance Profile
---

# aci_external_network_instance_profile #
Manages ACI External Network Instance Profile

## Example Usage ##

```hcl
	resource "aci_external_network_instance_profile" "fooexternal_network_instance_profile" {
		l3_outside_dn  = "${aci_l3_outside.example.id}"
		description    = "%s"
		name           = "demo_inst_prof"
		annotation     = "tag_network_profile"
		exception_tag  = "2"
		flood_on_encap = "disabled"
		match_t        = "%s"
		name_alias     = "alias_profile"
		pref_gr_memb   = "exclude"
		prio           = "level1"
		target_dscp    = "exclude"
	}
```
## Argument Reference ##
* `l3_outside_dn` - (Required) Distinguished name of parent L3Outside object.
* `name` - (Required) name of Object external_network_instance_profile.
* `annotation` - (Optional) annotation for object external_network_instance_profile.
* `exception_tag` - (Optional) exception_tag for object external_network_instance_profile. Allowed value range is "0" - "512".
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default value is "disabled".
* `match_t` - (Optional) The provider label match criteria. Allowed values are "All", "AtleastOne", "AtmostOne" and "None". Default is "AtleastOne".
* `name_alias` - (Optional) name_alias for object external_network_instance_profile.
* `pref_gr_memb` - (Optional)  Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "include" and "exclude". Default is "exclude".

* `prio` - (Optional) The QoS priority class identifier. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified".
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".	


* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_l3ext_rs_l3_inst_p_to_dom_p` - (Optional) Relation to class extnwDomP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_l3ext_rs_inst_p_to_nat_mapping_e_pg` - (Optional) Relation to class fvAEPg. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_l3ext_rs_inst_p_to_profile` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_M. Type - Set of Map.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the External Network Instance Profile.

## Importing ##

An existing External Network Instance Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_external_network_instance_profile.example <Dn>
```