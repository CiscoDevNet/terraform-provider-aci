---
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_interface_profile"
sidebar_current: "docs-aci-data-source-l3out_hsrp_interface_profile"
description: |-
  Data source for ACI L3out HSRP Interface Profile
---

# aci_l3out_hsrp_interface_profile

Data source for ACI L3out HSRP Interface Profile

## Example Usage

```hcl
data "aci_l3out_hsrp_interface_profile" "example" {
  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out HSRP Interface Profile.
- `annotation` - (Optional) Annotation for object l3out hsrp interface profile.
- `name_alias` - (Optional) Name alias for object l3out hsrp interface profile.
- `version` - (Optional) Compatibility catalog version
