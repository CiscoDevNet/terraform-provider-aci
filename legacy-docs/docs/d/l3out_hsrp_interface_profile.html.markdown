---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_interface_profile"
sidebar_current: "docs-aci-data-source-l3out_hsrp_interface_profile"
description: |-
  Data source for ACI L3-out HSRP interface profile
---

# aci_l3out_hsrp_interface_profile

Data source for ACI L3-out HSRP interface profile

## Example Usage

```hcl
data "aci_l3out_hsrp_interface_profile" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3-out HSRP interface profile.
- `annotation` - (Optional) Annotation for object L3-out HSRP interface profile.
- `name_alias` - (Optional) Name alias for object L3-out HSRP interface profile.
- `version` - (Optional) Compatibility catalog version.
- `description` - (Optional) Description for object L3-out HSRP interface profile.
