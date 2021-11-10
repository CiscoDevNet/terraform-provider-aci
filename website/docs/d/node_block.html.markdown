---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_node_block"
sidebar_current: "docs-aci-data-source-node_block"
description: |-
  Data source for ACI Node Block
---

# aci_node_block

Data source for ACI Node Block

## Example Usage

```hcl
data "aci_node_block" "example" {
  switch_association_dn   = aci_leaf_selector.example.id
  name                    = "example"
}
```

## Argument Reference

- `switch_association_dn` - (Required) Distinguished name of parent Leaf selector object.
- `name` - (Required) Name of Object node block.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Node Block.
- `annotation` - (Optional) Annotation for object node block.
- `description` - (Optional) Description for object node block.
- `from_` - (Optional) From Node ID.
- `name_alias` - (Optional) Name alias for object node block.
- `to_` - (Optional) To node ID.
