---
layout: "aci"
page_title: "ACI: aci_access_sub_port_block"
sidebar_current: "docs-aci-resource-access_sub_port_block"
description: |-
  Manages ACI Access Sub Port Block
---

# aci_access_sub_port_block #

Manages ACI Access Sub Port Block

## Example Usage ##

```hcl
resource "aci_access_sub_port_block" "example" {
  access_port_selector_dn  = "${aci_access_port_selector.example.id}"
  name                     = "example"
  annotation               = "example"
  from_card                = "example"
  from_port                = "example"
  from_sub_port            = "example"
  name_alias               = "example"
  to_card                  = "example"
  to_port                  = "example"
  to_sub_port              = "example"
}
```

## Argument Reference ##

* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) name of Object access_sub_port_block.
* `annotation` - (Optional) annotation for object access_sub_port_block.
* `from_card` - (Optional) from card
* `from_port` - (Optional) port block from port
* `from_sub_port` - (Optional) from_sub_port for object access_sub_port_block.
* `name_alias` - (Optional) name_alias for object access_sub_port_block.
* `to_card` - (Optional) to card
* `to_port` - (Optional) to port
* `to_sub_port` - (Optional) to_sub_port for object access_sub_port_block.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Sub Port Block.

## Importing ##

An existing Access Sub Port Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_access_sub_port_block.example <Dn>
```
