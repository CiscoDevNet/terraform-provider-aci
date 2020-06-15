---
layout: "aci"
page_title: "ACI: aci_node_block"
sidebar_current: "docs-aci-data-source-node_block"
description: |-
  Data source for ACI Node Block
---

# aci_node_block #
Data source for ACI Node Block

## Example Usage ##

```hcl
data "aci_node_block" "example" {
  switch_association_dn   = "${aci_switch_association.example.id}"
  name                    = "example"
}
```

## Argument Reference ##
* `switch_association_dn` - (Required) Distinguished name of parent SwitchAssociation object.
* `name` - (Required) name of Object node_block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Node Block.
* `annotation` - (Optional) annotation for object node_block.
* `from_` - (Optional) from Node ID. Range from 101 to 110.
* `name_alias` - (Optional) name_alias for object node_block.
* `to_` - (Optional) to node ID. Range from 101 to 110.
