---
layout: "aci"
page_title: "ACI: aci_leaf_breakout_port_group"
sidebar_current: "docs-aci-data-source-leaf_breakout_port_group"
description: |-
  Data source for ACI Leaf Breakout Port Group
---

# aci_leaf_breakout_port_group #
Data source for ACI Leaf Breakout Port Group

## Example Usage ##

```hcl
data "aci_leaf_breakout_port_group" "example" {
  name  = "example"
}
```


## Argument Reference ##
* `name` - (Required) name of leaf breakout port group object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the leaf breakout port group object.
* `annotation` - (Optional) annotation for leaf breakout port group object.
* `brkout_map` - (Optional) breakout map for leaf breakout port group object.
* `name_alias` - (Optional) name alias for leaf breakout port group object.
