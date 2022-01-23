---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_l3_interface_policy"
sidebar_current: "docs-aci-resource-l3_interface_policy"
description: |-
  Manages ACI L3 Interface Policy
---

# aci_l3_interface_policy #
Manages ACI L3 Interface Policy

## API Information ##
* `Class` - l3IfPol
* `Distinguished Name` - uni/fabric/l3IfP-{name}

## GUI Information ##
* `Location` - Fabric -> Fabric Policies -> Policies -> Interface -> L3 Interface -> Create L3 Interface Policy

## Example Usage ##

```hcl
resource "aci_l3_interface_policy" "example" {
  name  = "example"
  annotation  = "example"
  bfd_isis = "disabled"
  name_alias  = "example"
  description = "example"
}
```

## Argument Reference ##
* `name` - (Required) Name of object L3 Interface Policy.
* `annotation` - (Optional) Annotation for object L3 Interface Policy.
* `bfd_isis` - (Optional) BFD ISIS Configuration for object L3 Interface Policy. Allowed values are "disabled" and "enabled". Default value is "disabled".
* `name_alias` - (Optional) Name alias for object L3 Interface Policy.
* `description` - (Optional) Description for object L3 Interface Policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3 Interface Policy.

## Importing ##

An existing L3 Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l3_interface_policy.example <Dn>
```