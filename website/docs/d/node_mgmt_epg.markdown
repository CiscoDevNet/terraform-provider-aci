---
layout: "aci"
page_title: "ACI: aci_node_mgmt_epg"
sidebar_current: "docs-aci-data-source-node_mgmt_epg"
description: |-
  Data source for ACI Node Management EPg
---

# aci_in_band_management_e_pg

Data source for ACI Node Management EPg

## Example Usage

```hcl

data "aci_node_mgmt_epg" "example" {
  type = "in_band"
  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"
}

data "aci_node_mgmt_epg" "example" {
  type = "out_of_band"
  management_profile_dn  = "${aci_management_profile.example.id}"
  name  = "example"
}

```

## Argument Reference

- `type` - (Required) Type of node management EPg to be configured.  
  Allowed values: "in_band", "out_of_band".
- `management_profile_dn` - (Required) Distinguished name of parent management profile object.
- `name` - (Required) Name of Object node management EPg.

## Attribute Reference

### `type = "in_band"`

- `id` - Attribute id set to the Dn of the Node Management EPg.
- `annotation` - (Optional) Annotation for object in-band management EPg.
- `encap` - (Optional) The in-band access encapsulation.
- `exception_tag` - (Optional) Exception tag for object in-band management EPg.
- `flood_on_encap` - (Optional) Control at EPg level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
- `match_t` - (Optional) The provider label match criteria.
- `name_alias` - (Optional) Name alias for object in-band management EPg.
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
- `prio` - (Optional) The QoS priority class identifier.

### `type = "out_of_band"`

- `id` - Attribute id set to the Dn of the Node Management EPg.
- `annotation` - (Optional) annotation for object out-of-band management EPg.
- `name_alias` - (Optional) Name alias for object out-of-band management EPg.
- `prio` - (Optional) The QoS priority class identifier.