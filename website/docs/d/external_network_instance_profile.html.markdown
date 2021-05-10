---
layout: "aci"
page_title: "ACI: aci_external_network_instance_profile"
sidebar_current: "docs-aci-data-source-external_network_instance_profile"
description: |-
  Data source for ACI External Network Instance Profile
---

# aci_external_network_instance_profile

Data source for ACI External Network Instance Profile

## Example Usage

```hcl
data "aci_external_network_instance_profile" "dev_ext_net_prof" {
  l3_outside_dn  = "${aci_l3_outside.example.id}"
  name           = "foo_ext_net_prof"
}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of parent L3Outside object.
- `name` - (Required) Name of Object external network instance profile.

## Attribute Reference

- `id` - Attribute id set to the Dn of the External Network Instance Profile.
- `annotation` - (Optional) Annotation for object external network instance profile.
- `exception_tag` - (Optional) Exception tag for object external network instance profile.
- `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings.
- `match_t` - (Optional) The provider label match criteria.
- `name_alias` - (Optional) Name alias for object external network instance profile.
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication.
- `prio` - (Optional) The QoS priority class identifier.
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
