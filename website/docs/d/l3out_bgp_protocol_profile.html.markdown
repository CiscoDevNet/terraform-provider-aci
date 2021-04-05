---
layout: "aci"
page_title: "ACI: aci_l3out_bgp_protocol_profile"
sidebar_current: "docs-aci-data-source-l3out_bgp_protocol_profile"
description: |-
  Data source for ACI L3out BGP Protocol Profile
---

# aci_l3out_bgp_protocol_profile

Data source for ACI L3out BGP Protocol Profile

## Example Usage

```hcl
data "aci_l3out_bgp_protocol_profile" "example" {
  logical_node_profile_dn  = "${aci_logical_node_profile.example.id}"
}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of parent logical node profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out BGP Protocol Profile.
- `annotation` - (Optional) Annotation for object L3out BGP Protocol Profile.
- `name_alias` - (Optional) Name alias for object L3out BGP Protocol Profile.
