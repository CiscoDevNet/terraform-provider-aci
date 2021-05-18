---
layout: "aci"
page_title: "ACI: aci_l3out_bfd_interface_profile"
sidebar_current: "docs-aci-resource-l3out_bfd_interface_profile"
description: |-
  Manages ACI L3out BFD Interface Profile
---

# aci_l3out_bfd_interface_profile

Manages ACI L3out BFD Interface Profile

## Example Usage

```hcl
resource "aci_l3out_bfd_interface_profile" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  annotation                   = "example"
  description                  = "from terraform"
  key                          = "example"
  key_id                       = "25"
  name_alias                   = "example"
  interface_profile_type       = "sha1"
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `annotation` - (Optional) Annotation for L3out BFD interface profile object.
- `name_alias` - (Optional) Name alias for L3out BFD interface profile object.
- `description` - (Optional) Description for L3out BFD interface profile object.
- `key` - (Optional) Password to identify this L3out BFD interface profile object.
- `key_id` - (Optional) Authentication key id for L3out BFD interface profile object. Default value is "1".
- `interface_profile_type` - (Optional) Component type for L3out BFD interface profile object. Allowed values are "none" and "sha1". Default value is "none".

- `relation_bfd_rs_if_pol` - (Optional) Relation to class bfdIfPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Interface Profile.

## Importing

An existing Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_bfd_interface_profile.example <Dn>
```
