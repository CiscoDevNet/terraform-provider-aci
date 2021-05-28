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

  firmware_group_dn  = aci_firmware_group.example.id

  description = "from terraform"
  name  = "example"
}
```
## Argument Reference ##
* `firmware_group_dn` - (Required) Distinguished name of parent Firmware Group object.
* `name` - (Required) Name of Object Node Block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Node Block.
* `description` - (Optional) Description for Object Node Block.
* `annotation` - (Optional) Annotation for Object Node Block.
* `from_` - (Optional) From value for Object Node Block.
* `name_alias` - (Optional) Name alias for Object Node Block..
* `to_` - (Optional) To value for Object Node Block.
