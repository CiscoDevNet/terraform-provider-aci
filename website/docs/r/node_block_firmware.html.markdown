---
layout: "aci"
page_title: "ACI: aci_node_block_firmware"
sidebar_current: "docs-aci-resource-node_block_firmware"
description: |-
  Manages ACI Node Block
---

# aci_node_block_firmware #
Manages ACI Node Block

## Example Usage ##

```hcl
resource "aci_node_block_firmware" "example" {

  firmware_group_dn  = aci_firmware_group.example.id

  name  = "example"
  description = "from terraform"
  annotation  = "example"
  from_  = "example"
  name_alias  = "example"
  to_  = "example"
}
```
## Argument Reference ##
* `firmware_group_dn` - (Required) Distinguished name of parent Firmware Group Object.
* `name` - (Required) Name of Object Node Block.
* `description` - (Optional) Description for Object Node Block.
* `annotation` - (Optional) Annotation for Object Node Block.
* `from_` - (Optional) From value for Object Node Block.
* `name_alias` - (Optional) Name alias for Object Node Block.
* `to_` - (Optional) To value for Object Node Block.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Node Block.

## Importing ##

An existing Node Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_node_block_firmware.example <Dn>
```
