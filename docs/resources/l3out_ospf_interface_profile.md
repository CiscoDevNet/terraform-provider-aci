---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_ospf_interface_profile"
sidebar_current: "docs-aci-resource-aci_l3out_ospf_interface_profile"
description: |-
  Manages ACI L3out OSPF Interface Profile
---

# aci_l3out_ospf_interface_profile

Manages ACI L3out OSPF Interface Profile

## Example Usage

```hcl
resource "aci_l3out_ospf_interface_profile" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  description                  = "from terraform"
  annotation                   = "example"
  auth_key                     = "key"
  auth_key_id                  = "255"
  auth_type                    = "simple"
  name_alias                   = "example"
  relation_ospf_rs_if_pol      = aci_ospf_interface_policy.example.id
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of the parent logical interface profile object.
- `auth_key` - (Optional) OSPF authentication key for L3out OSPF interface profile object.
- `annotation` - (Optional) Annotation for L3out OSPF interface profile object.
- `description` - (Optional) Description for L3out OSPF interface profile object.
- `auth_key_id` - (Optional) Authentication key id for L3out OSPF interface profile object. Allowed ranges are from "1" to "255". The default value is "1".
- `auth_type` - (Optional) OSPF authentication type for L3out OSPF interface profile object. Allowed values are "none", "md5" and "simple". Default value is "none".
- `name_alias` - (Optional) Name alias for L3out OSPF interface profile object.

- `relation_ospf_rs_if_pol` - (Optional) Relation to class ospfIfPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out OSPF Interface Profile.

## Importing

An existing L3out OSPF Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_ospf_interface_profile.example <Dn>
```
