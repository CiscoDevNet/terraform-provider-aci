---
layout: "aci"
page_title: "ACI: aci_cloud_external_epg"
sidebar_current: "docs-aci-resource-cloud_external_e_pg"
description: |-
  Manages ACI Cloud External EPg
---

# aci_cloud_external_epg #
Manages ACI Cloud External EPg
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
	resource "aci_cloud_external_epg" "foocloud_external_e_pg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		description                   = "sample cloud external epg"
		name                          = "cloud_ext_epg"
		annotation                    = "tag_ext_epg"
		exception_tag                 = "0"
		flood_on_encap                = "disabled"
		match_t                       = "All"
		name_alias                    = "alias_ext"
		pref_gr_memb                  = "exclude"
		prio                          = "unspecified"
		route_reachability            = "inter-site"
	}
```
## Argument Reference ##
* `cloud_applicationcontainer_dn` - (Required) Distinguished name of parent CloudApplicationcontainer object.
* `name` - (Required) name of Object cloud_external_e_pg.
* `annotation` - (Optional) annotation for object cloud_external_e_pg.
* `exception_tag` - (Optional) exception_tag for object cloud_external_e_pg. Allowed value range is "0" to "512".
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default is "disabled".
* `match_t` - (Optional) The provider label match criteria. Allowed values are "All", "AtleastOne", "AtmostOne" and "None". Default values is "AtleastOne". 
* `name_alias` - (Optional) name_alias for object cloud_external_e_pg.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "include" and "exclude". Default is "exclude".
* `prio` - (Optional) qos priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
* `route_reachability` - (Optional) Route reachability for this EPG. Allowed values are "unspecified", "inter-site", "internet" and "inter-site-ext".Default is "inter-site".

* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_cloud_rs_cloud_e_pg_ctx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud External EPg.

## Importing ##

An existing Cloud External EPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_external_epg.example <Dn>
```