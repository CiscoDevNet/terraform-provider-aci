---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_application_epg"
sidebar_current: "docs-aci-data-source-aci_application_epg"
description: |-
  Data source for ACI Application EPG
---

# aci_application_epg #
Data source for ACI Application EPG

## Example Usage ##

```hcl
data "aci_application_epg" "foo_epg" {
  application_profile_dn  = aci_application_profile.foo_app.id
  name                    = "dev_app_epg"
}
```
## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) Name of Object application epg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Application EPG.
* `annotation` - (Optional) Annotation for object application epg.
* `description` - (Optional) Description for object application epg.
* `exception_tag` - (Optional) Exception tag for object application epg. Range: "0" - "512" .
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default is "disabled".
* `fwd_ctrl` - (Optional) Forwarding control at EPG level. Allowed values are "none" and "proxy-arp". Default is "none".
* `has_mcast_source` - (Optional) If the source for the EPG is multicast or not. 
* `is_attr_based_epg` - (Optional) If the EPG is attribute based or not.
* `match_t` - (Optional) The provider label match criteria for EPG. 
* `name_alias` - (Optional) Name alias for object application epg.
* `pc_enf_pref` - (Optional) The preferred policy control.
* `pc_tag` - (Read-Only) A numeric ID to represent a policy enforcement group.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
* `prio` - (Optional) QoS priority class id. 
* `shutdown` - (Optional) Shutdown for object application epg.
* `relation_fv_rs_sec_inherited` - (Optional) Relation to another epg (fvEPg class) to inherit contracts from. Named 'contract master' in GUI.
