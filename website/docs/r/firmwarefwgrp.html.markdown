---
layout: "aci"
page_title: "ACI: aci_firmware_group"
sidebar_current: "docs-aci-resource-firmware_group"
description: |-
  Manages ACI Firmware Group
---

# aci_firmware_group #
Manages ACI Firmware Group

## Example Usage ##

```hcl
resource "aci_firmware_group" "example" {


  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  firmware_group_type  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_group.
* `annotation` - (Optional) annotation for object firmware_group.
* `name_alias` - (Optional) name_alias for object firmware_group.
* `firmware_group_type` - (Optional) component type

* `relation_firmware_rs_fwgrpp` - (Optional) Relation to class firmwareFwP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Firmware Group.

## Importing ##

An existing Firmware Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_firmware_group.example <Dn>
```