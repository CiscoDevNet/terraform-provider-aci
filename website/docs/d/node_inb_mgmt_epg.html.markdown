---
layout: "aci"
page_title: "ACI: aci_node_inb_mgmt_epg"
sidebar_current: "docs-aci-data-source-node_inb_mgmt_epg"
description: |-
  Data source for ACI In-Band Management EPg
---

# aci_node_inb_mgmt_epg #
Data source for ACI In-Band Management EPg

## Example Usage ##

```hcl
data "aci_node_inb_mgmt_epg" "example" {

  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"

}
```
## Argument Reference ##
* `management_profile_dn` - (Required) Distinguished name of parent management profile object.
* `name` - (Required) name of Object in band management EPg.



## Attribute Reference

* `id` - Attribute id set to the Dn of the In-Band Management EPg.
* `annotation` - (Optional) annotation for object in band management EPg.
* `encap` - (Optional) inband access encap
* `exception_tag` - (Optional) exception_tag for object in band management EPg.
* `flood_on_encap` - (Optional) flood_on_encap for object in band management EPg.
* `match_t` - (Optional) match criteria
* `name_alias` - (Optional) name_alias for object in band management EPg.
* `pref_gr_memb` - (Optional) pref_gr_memb for object in band management EPg.
* `prio` - (Optional) in-band qos priority
