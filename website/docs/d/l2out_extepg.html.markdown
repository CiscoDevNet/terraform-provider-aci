---
layout: "aci"
page_title: "ACI: aci_l2out_extepg"
sidebar_current: "docs-aci-data-source-l2out_extepg"
description: |-
  Data source for ACI L2-Out External EPg
---

# aci_l2out_extepg #

Data source for ACI L2-Out External EPg

## Example Usage ##

```hcl
data "aci_l2out_extepg" "example" {
  l2_outside_dn  = aci_l2_outside.example.id
  name  = "l2out_extepg_1"
}
```

## Argument Reference ##

- `l2_outside_dn` - (Required) Distinguished name of parent L2 Outside object.
- `name` - (Required) The name of the layer 2 external network instance profile. This name can be up to 64 alphanumeric characters. Note that you cannot change this name after the object has been saved.

## Attribute Reference ##

- `id` - Attribute id set to the Dn of the L2 Out External EPG.
- `annotation` - (Optional) Annotation for object L2 Out External EPG.
- `description` - (Optional) Description for object L2 Out External EPG.
- `exception_tag` - (Optional) Exception tag for object L2 Out External EPG.
- `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link-Local Layer should be flooded only on ENCAP or based on bridge-domain settings.
- `match_t` - (Optional) The provider label match criteria.
- `name_alias` - (Optional) Name alias for object L2 Out External EPG.
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPG is part of a group that does not a contract for communication.
- `prio` - (Optional) The QoS priority class identifier.
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
