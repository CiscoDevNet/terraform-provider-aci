---
layout: "aci"
page_title: "ACI: aci_application_epg"
sidebar_current: "docs-aci-resource-application_epg"
description: |-
  Manages ACI Application EPG
---

# aci_application_epg #
Manages ACI Application EPG

## Example Usage ##

```hcl
  resource "aci_application_epg" "fooapplication_epg" {
    application_profile_dn  = "${aci_application_profile.app_profile_for_epg.id}"
    name  					        = "demo_epg"
    description 			      = "%s"
    annotation  			      = "tag_epg"
    exception_tag 		    	= "0"
    flood_on_encap  	      = "disabled"
    fwd_ctrl  			      	= "none"
    has_mcast_source     		= "no"
    is_attr_based_e_pg     	= "no"
    match_t  				        = "AtleastOne"
    name_alias  		      	= "alias_epg"
    pc_enf_pref  		      	= "unenforced"
    pref_gr_memb  	    		= "exclude"
    prio  				        	= "unspecified"
    shutdown  		      		= "no"
  }
```
## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) name of Object application_epg.
* `annotation` - (Optional) annotation for object application_epg.
* `exception_tag` - (Optional) exception_tag for object application_epg. Range: "0" - "512" .
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default is "disabled".
* `fwd_ctrl` - (Optional) Forwarding control at EPG level. Allowed values are "none" and "proxy-arp". Default is "none".
* `has_mcast_source` - (Optional) If the source for the EPG is multicast or not. Allowed values are "yes" and "no". Default values is "no".
* `is_attr_based_e_pg` - (Optional) If the EPG is attribute based or not. Allowed values are "yes" and "no". Default is "yes".
* `match_t` - (Optional) The provider label match criteria for EPG. Allowed values are "All", "AtleastOne", "AtmostOne", "None". Default is "AtleastOne".
* `name_alias` - (Optional) name_alias for object application_epg.
* `pc_enf_pref` - (Optional) The preferred policy control. Allowed values are "unenforced" and "enforced". Default is "unenforced".
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "exclude" and "include". Default is "exclude".
* `prio` - (Optional) qos priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
* `shutdown` - (Optional) shutdown for object application_epg. Allowed values are "yes" and "no". Default is "no".

* `relation_fv_rs_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_dom_att` - (Optional) Relation to class infraDomP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_fc_path_att` - (Optional) Relation to class fabricPathEp. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_graph_def` - (Optional) Relation to class vzGraphCont. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_node_att` - (Optional) Relation to class fabricNode. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_dpp_pol` - (Optional) Relation to class qosDppPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_prov_def` - (Optional) Relation to class vzCtrctEPgCont. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_trust_ctrl` - (Optional) Relation to class fhsTrustCtrlPol. Cardinality - N_TO_ONE. Type - String.
                                
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_ae_pg_mon_pol` - (Optional) Relation to class monEPGPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_intra_epg` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Application EPG.

## Importing ##

An existing Application EPG can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_application_epg.example <Dn>
```