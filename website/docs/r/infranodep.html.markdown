---
layout: "aci"
page_title: "ACI: aci_leaf_profile"
sidebar_current: "docs-aci-resource-leaf_profile"
description: |-
  Manages ACI Leaf Profile
---

# aci_leaf_profile #
Manages ACI Leaf Profile

## Example Usage ##

```hcl

resource "aci_leaf_profile" "example" {
  name        = "leaf1"
  annotation  = "example"
  name_alias  = "example"
  leaf_selector {
    name                    = "one"
    switch_association_type = "range"
    node_block {
      name  = "blk1"
      from_ = "105"
      to_   = "106"
    }
    node_block {
      name  = "blk2"
      from_ = "102"
      to_   = "104"
    }
  }
  leaf_selector {
    name                    = "two"
    switch_association_type = "range"
    node_block {
      name  = "blk3"
      from_ = "105"
      to_   = "106"
    }
  }
}

```

## Argument Reference ##
* `name` - (Required) name of Object leaf_profile.
* `annotation` - (Optional) annotation for object leaf_profile.
* `name_alias` - (Optional) name_alias for object leaf_profile.

* `leaf_selector` - (Optional) leaf Selector block to attach with the leaf profile.
* `leaf_selector.name` - (Required) name of the leaf selector.
* `leaf_selector.switch_association_type` - (Required) type of switch association. 
Allowed values: "ALL", "range", "ALL_IN_POD"

* `leaf_selector.node_block` - (Optional) Node block to attach with leaf selector.
* `leaf_selector.node_block.name` - (Required) Name of the node block.
* `leaf_selector.node_block.from_` - (Optional) from Node ID. Range from 101 to 110.
* `leaf_selector.node_block.to_` - (Optional) to node ID. Range from 101 to 110.

* `relation_infra_rs_acc_card_p` - (Optional) Relation to class infraAccCardP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_acc_port_p` - (Optional) Relation to class infraAccPortP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Profile.

## Importing ##

An existing Leaf Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_profile.example <Dn>
```