---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_application_epg"
sidebar_current: "docs-aci-resource-application_epg"
description: |-
  Manages ACI Application EPG
---

# aci_application_epg

Manages ACI Application EPG

## Example Usage

```hcl
resource "aci_application_epg" "fooapplication_epg" {
    application_profile_dn  = aci_application_profile.app_profile_for_epg.id
    name  					        = "demo_epg"
    description 			      = "from terraform"
    annotation  			      = "tag_epg"
    exception_tag 		    	= "0"
    flood_on_encap  	      = "disabled"
    fwd_ctrl  			      	= "none"
    has_mcast_source     		= "no"
    is_attr_based_epg     	= "no"
    match_t  				        = "AtleastOne"
    name_alias  		      	= "alias_epg"
    pc_enf_pref  		      	= "unenforced"
    pref_gr_memb  	    		= "exclude"
    prio  				        	= "unspecified"
    shutdown  		      		= "no"
}
```

## Argument Reference ##
* `name` - (Required) Name of Object application epg.
* `annotation` - (Optional) Annotation for object application epg.
* `description` - (Optional) Description for object application epg.
* `exception_tag` - (Optional) Exception tag for object application epg. Range: "0" - "512" .
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default is "disabled".
* `fwd_ctrl` - (Optional) Forwarding control at EPG level. Allowed values are "none" and "proxy-arp". Default is "none".
* `has_mcast_source` - (Optional) If the source for the EPG is multicast or not. Allowed values are "yes" and "no". Default values is "no".
* `is_attr_based_epg` - (Optional) If the EPG is attribute based or not. Allowed values are "yes" and "no". Default is "no".
* `match_t` - (Optional) The provider label match criteria for EPG. Allowed values are "All", "AtleastOne", "AtmostOne", "None". Default is "AtleastOne".
* `name_alias` - (Optional) Name alias for object application epg.
* `pc_enf_pref` - (Optional) The preferred policy control. Allowed values are "unenforced" and "enforced". Default is "unenforced".
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "exclude" and "include". Default is "exclude".
* `prio` - (Optional) QoS priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4","level5" and "level6". By default the value is inherited from the parent application profile.
* `shutdown` - (Optional) Shutdown for object application epg. Allowed values are "yes" and "no". Default is "no".

* `relation_fv_rs_bd` - (Required) Relation to Bridge domain associated with EPG (Point to class fvBD). Cardinality - N_TO_ONE. Type - String.

* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to Custom Quality of Service traffic policy name (Point to class qosCustomPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> Custom QoS -->

* `relation_fv_rs_fc_path_att` - (Optional) Relation to Fibre Channel (Paths) (Point to class fabricPathEp). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_prov` - (Optional) Relation to Provided Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_cons_if` - (Optional) Relation to Imported Contract (Point to class vzCPIf). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_sec_inherited` - (Optional) Relation represents that the EPG is inheriting security configuration from other EPGs (Point to class fvEPg). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_node_att` - (Optional) Relation used to define a Static Leaf binding (Point to class fabricNode). Cardinality - N_TO_M. Type - Set of String.
<!-- tenant -> Application Profile -> EPG ->Static Leaf -->

* `relation_fv_rs_dpp_pol` - (Optional) Relation to define a Data Plane Policing policy (Point to class qosDppPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> Data Plane Policing -->

* `relation_fv_rs_cons` - (Optional) Relation to Consumed Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_trust_ctrl` - (Optional) Relation to First Hop Security trust control (Point to class fhsTrustCtrlPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> First Hop Security -->

* `relation_fv_rs_prot_by` - (Optional) Relation to Taboo Contract (Point to class vzTaboo). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_aepg_mon_pol` - (Optional) Relation to create a container for monitoring policies associated with the tenant. This allows you to apply tenant-specific policies (Point to class monEPGPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> Monitoring -->

* `relation_fv_rs_intra_epg` - (Optional) Relation to Intra EPG Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Application EPG.

## Importing

An existing Application EPG can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_application_epg.example <Dn>
```
