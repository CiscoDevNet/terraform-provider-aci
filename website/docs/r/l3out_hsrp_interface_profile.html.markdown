---
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_interface_profile"
sidebar_current: "docs-aci-resource-l3out_hsrp_interface_profile"
description: |-
  Manages ACI L3-out HSRP interface profile
---

# aci_l3out_hsrp_interface_profile

Manages ACI L3-out HSRP interface profile

## Example Usage

```hcl
resource "aci_l3out_hsrp_interface_profile" "example" {

  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
  annotation  = "example"
  name_alias  = "example"
  description = "from terraform"
  version = "v1"

}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `annotation` - (Optional) Annotation for object L3-out HSRP interface profile.
- `description` - (Optional) Description for object L3-out HSRP interface profile.
- `name_alias` - (Optional) Name alias for object L3-out HSRP interface profile.
- `version` - (Optional) Compatibility catalog version.  
  Allowed values: "v1", "v2". Default value: "v1".
- `relation_hsrp_rs_if_pol` - (Optional) Relation to class hsrpIfPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3-out HSRP interface profile.

## Importing

An existing L3-out HSRP interface profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_hsrp_interface_profile.example <Dn>
```
