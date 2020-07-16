---
layout: "aci"
page_title: "ACI: aci_application_epg"
sidebar_current: "docs-aci-data-source-application_epg"
description: |-
  Data source for ACI Application EPG
---

# aci_application_epg #
Data source for ACI Application EPG

## Example Usage ##

```hcl
data "aci_application_epg" "foo_epg" {

  application_profile_dn  = "${aci_application_profile.foo_app.id}"
  name                    = "dev_app_epg"
}
```
## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) name of Object application_epg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Application EPG.
* `annotation` - (Optional) annotation for object application_epg.
* `exception_tag` - (Optional) exception_tag for object application_epg.
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
* `fwd_ctrl` - (Optional) Forwarding control at EPG level.
* `has_mcast_source` - (Optional) If the source for the EPG is multicast or not.
* `is_attr_based_epg` - (Optional) If the EPG is attribute based or not.
* `match_t` - (Optional) The provider label match criteria for EPG.
* `name_alias` - (Optional) name_alias for object application_epg.
* `pc_enf_pref` - (Optional) The preferred policy control. 
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
* `prio` - (Optional) qos priority class id
* `shutdown` - (Optional) shutdown for object application_epg.
