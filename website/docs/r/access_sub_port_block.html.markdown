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
  access_port_selector_dn  = aci_access_port_selector.example.id
  description             = "From Terraform"
  name                    = "example"
  annotation              = "example"
  from_card               = "1"
  from_port               = "1"
  from_sub_port           = "1"
  name_alias              = "example"
  to_card                 = "1"
  to_port                 = "1"
  to_sub_port             = "1"
}
```

## Argument Reference ##

* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) Name of Object access sub port block.
* `annotation` - (Optional) Annotation for object access sub port block.
* `description` - (Optional) Description for object access sub port block.
* `from_card` - (Optional) From card.
  Allowed Values are between 1 to 100. Default Value is "1".
* `from_port` - (Optional) Port block from port
  Allowed Values are between 1 to 127. Default Value is "1".
* `from_sub_port` - (Optional) From sub port for object access sub port block.
  Allowed Values are between 1 to 64. Default Value is "1".
* `name_alias` - (Optional) Name alias for object access sub port block.
* `to_card` - (Optional) To card.
  Allowed Values are between 1 to 100. Default Value is "1".
* `to_port` - (Optional) To port.
 Allowed Values are between 1 to 127. Default Value is "1".
* `to_sub_port` - (Optional) To sub port for object access sub port block.
  Allowed Values are between 1 to 64. Default Value is "1".


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Sub Port Block.

## Importing ##

An existing Access Sub Port Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_access_sub_port_block.example <Dn>
```
