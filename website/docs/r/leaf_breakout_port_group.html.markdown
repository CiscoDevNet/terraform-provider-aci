---
layout: "aci"
page_title: "ACI: aci_leaf_breakout_port_group"
sidebar_current: "docs-aci-resource-leaf_breakout_port_group"
description: |-
  Manages ACI Leaf Breakout Port Group
---

# aci_leaf_breakout_port_group #
Manages ACI Leaf Breakout Port Group

## Example Usage ##

```hcl
resource "aci_leaf_breakout_port_group" "example" {
  name        = "first"
  annotation  = "example"
  brkout_map  = "100g-4x"
  name_alias  = "aliasing"
}
```


## Argument Reference ##

* `name` - (Required) Name of leaf breakout port group object.
* `annotation` - (Optional) Annotation for leaf breakout port group object.
* `brkout_map` - (Optional) Breakout map for leaf breakout port group object. Allowed values are "100g-2x", "100g-4x", "10g-4x", "25g-4x", "50g-8x" and "none". Default value is "none".
* `name_alias` - (Optional) Name alias for leaf breakout port group object.

* `relation_infra_rs_mon_brkout_infra_pol` - (Optional) Relation to class monInfraPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Breakout Port Group.

## Importing ##

An existing Leaf Breakout Port Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_breakout_port_group.example <Dn>
```