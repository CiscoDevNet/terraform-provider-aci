---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_bgp_route_control_profile"
sidebar_current: "docs-aci-data-source-aci_bgp_route_control_profile"
description: |-
  Data source for ACI BGP Route Control Profile
---

# aci_bgp_route_control_profile
!> **WARNING:** This data source is deprecated and will be removed in the next major version use aci_route_control_profile instead.
Data source for ACI BGP Route Control Profile

## Example Usage

```hcl
data "aci_bgp_route_control_profile" "check" {
  parent_dn = aci_tenant.tenentcheck.id
  name      = "one"
}

data "aci_bgp_route_control_profile" "check" {
  parent_dn = aci_l3_outside.example.id
  name      = "bgp_route_control_profile_1"
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of the parent object.
- `name` - (Required) Name of router control profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the router control profile object.
- `annotation` - Annotation for router control profile object.
- `description` - Description for router control profile object.
- `name_alias` - Name alias for router control profile object.
- `route_control_profile_type` - Component type for router control profile object.
