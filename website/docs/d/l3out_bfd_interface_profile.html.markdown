---
layout: "aci"
page_title: "ACI: aci_l3out_bfd_interface_profile"
sidebar_current: "docs-aci-data-source-l3out_bfd_interface_profile"
description: |-
  Data source for ACI L3out BFD Interface Profile
---

# aci_l3out_bfd_interface_profile

Data source for ACI L3out BFD Interface Profile

## Example Usage

```hcl
data "aci_l3out_bfd_interface_profile" "check" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out BFD interface profile object.
- `annotation` - Annotation for L3out BFD interface profile object.
- `description` - Description for L3out BFD interface profile object.
- `key_id` - Authentication key id for L3out BFD interface profile object.
- `name_alias` - Name alias for L3out BFD interface profile object.
- `interface_profile_type` - Component type for L3out BFD interface profile object.
- `userdom` - Userdom for L3out BFD interface profile object.
