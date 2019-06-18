---
layout: "aci"
page_title: "ACI: aci_autonomous_system_profile"
sidebar_current: "docs-aci-resource-autonomous_system_profile"
description: |-
  Manages ACI Autonomous System Profile
---

# aci_autonomous_system_profile #
Manages ACI Autonomous System Profile

## Example Usage ##

```hcl
resource "aci_autonomous_system_profile" "example" {

  annotation  = "example"
  asn  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `annotation` - (Optional) annotation for object autonomous_system_profile.
* `asn` - (Optional) autonomous system number
* `name_alias` - (Optional) name_alias for object autonomous_system_profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Autonomous System Profile.

## Importing ##

An existing Autonomous System Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_autonomous_system_profile.example <Dn>
```