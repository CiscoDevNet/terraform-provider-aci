---
layout: "aci"
page_title: "ACI: aci_l3out_static_route_next_hop"
sidebar_current: "docs-aci-data-source-l3out_static_route_next_hop"
description: |-
  Data source for ACI L3out Static Route Next Hop
---

# aci_l3out_static_route_next_hop

Data source for ACI L3out Static Route Next Hop

## Example Usage

```hcl
data "aci_l3out_static_route_next_hop" "example" {
  static_route_dn  = "${aci_l3out_static_route.example.id}"
  nh_addr  = "example"
}
```

## Argument Reference

- `static_route_dn` - (Required) Distinguished name of parent static route object.
- `nh_addr` - (Required) The nexthop IP address for the static route to the outside network.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out Static Route Next Hop.
- `annotation` - (Optional) Annotation for object l3out static route next hop.
- `description` - (Optional) Description for object l3out static route next hop.
- `name_alias` - (Optional) Name alias for object l3out static route next hop.
- `pref` - (Optional) Administrative preference value for this route.
- `nexthop_profile_type` - (Optional) Component type.
