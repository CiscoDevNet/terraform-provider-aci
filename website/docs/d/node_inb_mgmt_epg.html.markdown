---
layout: "aci"
page_title: "ACI: aci_node_inb_mgmt_epg"
sidebar_current: "docs-aci-data-source-node_inb_mgmt_epg"
description: |-
  Data source for ACI Node In-Band Management EPg
---

# aci_node_inb_mgmt_epg

Data source for ACI Node In-Band Management EPg

## Example Usage

```hcl
data "aci_node_inb_mgmt_epg" "example" {

  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"

}
```

## Argument Reference

- `management_profile_dn` - (Required) Distinguished name of parent management profile object.
- `name` - (Required) The in-band management endpoint group name. This name can be up to 64 alphanumeric characters.

## Attribute Reference

- `id` - Attribute id set to the Dn of the node in-band management EPg.
- `annotation` - (Optional) annotation for object node in-band management EPg.
- `encap` - (Optional) The in-band access encapsulation.
- `exception_tag` - (Optional) Exception-tag for object node in-band management EPg.
- `flood_on_encap` - (Optional) Control at EPg level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
- `match_t` - (Optional) The provider label match criteria.
- `name_alias` - (Optional) Name alias for object node in-band management EPg.
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
- `prio` - (Optional) The in-band QoS priority class identifier.
