---
layout: "aci"
page_title: "ACI: aci_node_block"
sidebar_current: "docs-aci-resource-node_block"
description: |-
  Manages ACI Node Block
---

# aci_node_block

Manages ACI Node Block

## Example Usage

```hcl
resource "aci_node_block" "check" {
  switch_association_dn   = aci_leaf_selector.example.id
  name                    = "block"
  annotation              = "aci_node_block"
  from_                   = "105"
  name_alias              = "node_block"
  to_                     = "106"
}
```

## Argument Reference

- `switch_association_dn` - (Required) Distinguished name of parent Leaf selector object.
- `name` - (Required) Name of Object node block.
- `annotation` - (Optional) Annotation for object node block.
- `from_` - (Optional) From Node ID. Range from 101 to 110. Default value is "1".
- `name_alias` - (Optional) Name alias for object node block.
- `to_` - (Optional) To node ID. Range from 101 to 110. Default value is "1".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Node Block.

## Importing

An existing Node Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_node_block.example <Dn>
```
