---
layout: "aci"
page_title: "ACI: aci_node_block"
sidebar_current: "docs-aci-resource-node_block"
description: |-
  Manages ACI Node Block
---

# aci_node_block #
Manages ACI Node Block

## Example Usage ##

```hcl
resource "aci_node_block" "check" {
  switch_association_dn   = "${aci_switch_association.example.id}"
  name                    = "block"
  annotation              = "aci_node_block"
  from_                   = "105"
  name_alias              = "node_block"
  to_                     = "106"
}
```
## Argument Reference ##
* `switch_association_dn` - (Required) Distinguished name of parent SwitchAssociation object.
* `name` - (Required) name of Object node_block.
* `annotation` - (Optional) annotation for object node_block.
* `from_` - (Optional) from Node ID. Range from 101 to 110
* `name_alias` - (Optional) name_alias for object node_block.
* `to_` - (Optional) to node ID. Range from 101 to 110



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Node Block.

## Importing ##

An existing Node Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_node_block.example <Dn>
```