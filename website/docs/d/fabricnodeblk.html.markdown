---
layout: "aci"
page_title: "ACI: aci_node_block_firmware"
sidebar_current: "docs-aci-data-source-node_block_firmware"
description: |-
  Data source for ACI Node Block
---

# aci_node_block_firmware #
Data source for ACI Node Block

## Example Usage ##

```hcl
data "aci_node_block_firmware" "example" {

  firmware_group_dn  = "${aci_firmware_group.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `firmware_group_dn` - (Required) Distinguished name of parent FirmwareGroup object.
* `name` - (Required) name of Object node_block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Node Block.
* `annotation` - (Optional) annotation for object node_block.
* `from_` - (Optional) from
* `name_alias` - (Optional) name_alias for object node_block.
* `to_` - (Optional) to
