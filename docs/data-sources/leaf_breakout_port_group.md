---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_breakout_port_group"
sidebar_current: "docs-aci-data-source-aci_leaf_breakout_port_group"
description: |-
  Data source for ACI Leaf Breakout Port Group
---

# aci_leaf_breakout_port_group

Data source for ACI Leaf Breakout Port Group

## Example Usage

```hcl
data "aci_leaf_breakout_port_group" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of leaf breakout port group object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the leaf breakout port group object.
- `description` - (Optional) Description for leaf breakout port group object.
- `annotation` - (Optional) Annotation for leaf breakout port group object.
- `brkout_map` - (Optional) Breakout map for leaf breakout port group object.
- `name_alias` - (Optional) Name alias for leaf breakout port group object.
