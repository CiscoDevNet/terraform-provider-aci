---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_route_control_profile"
sidebar_current: "docs-aci-resource-aci_route_control_profile"
description: |-
  Manages ACI Route Control Profile
---

# aci_route_control_profile

Manages ACI Route Control Profile

## Example Usage

```hcl
resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenentcheck.id
  name                       = "one"
  annotation                 = "route_control_profile_tag"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}

resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_l3_outside.example.id
  name                       = "route_control_profile_1"
  annotation                 = "route_control_profile_tag"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of the parent object.
- `name` - (Required) Name of router control profile object.
- `annotation` - (Optional) Annotation for router control profile object.
- `description` - (Optional) Description for router control profile object.
- `name_alias` - (Optional) Name alias for router control profile object.
- `route_control_profile_type` - (Optional) Component type for router control profile object. Allowed values are "combinable" and "global". Default value is "combinable".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the Route Control Profile.

## Importing

An existing Route Control Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_route_control_profile.example <Dn>
```
