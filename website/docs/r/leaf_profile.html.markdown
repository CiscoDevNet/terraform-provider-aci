---
layout: "aci"
page_title: "ACI: aci_leaf_profile"
sidebar_current: "docs-aci-resource-leaf_profile"
description: |-
  Manages ACI Leaf Profile
---

# aci_leaf_profile

Manages ACI Leaf Profile

## Example Usage

```hcl

resource "aci_leaf_profile" "example" {
  name        = "leaf1"
  description  = "From Terraform"
  annotation  = "example"
  name_alias  = "example"
  leaf_selector {
    name                    = "one"
    switch_association_type = "range"
    node_block {
      name  = "blk1"
      from_ = "101"
      to_   = "102"
    }
    node_block {
      name  = "blk2"
      from_ = "103"
      to_   = "104"
    }
  }
  leaf_selector {
    name                    = "two"
    switch_association_type = "range"
    node_block {
      name  = "blk3"
      from_ = "105"
      to_   = "107"
    }
  }
}

```

## Argument Reference

- `name` - (Required) Name of Object leaf profile.
- `description` - (Optional) Description for object leaf profile.
- `annotation` - (Optional) Annotation for object leaf profile.
- `name_alias` - (Optional) Name alias for object leaf profile.

- `leaf_selector` - (Optional) Leaf Selector block to attach with the leaf profile.
- `leaf_selector.name` - (Required) Name of the leaf selector.
- `leaf_selector.switch_association_type` - (Required) Type of switch association.
  Allowed values: "ALL", "range", "ALL_IN_POD"

- `leaf_selector.node_block` - (Optional) Node block to attach with leaf selector.
- `leaf_selector.node_block.name` - (Required) Name of the node block.
- `leaf_selector.node_block.from_` - (Optional) Start of Node Block range. Range from 1 to 16000. Default value is "1".
- `leaf_selector.node_block.to_` - (Optional) End of Node Block range. Range from 1 to 16000. Default value is "1".

- `relation_infra_rs_acc_card_p` - (Optional) Relation to class infraAccCardP. Cardinality - N_TO_M. Type - Set of String.
- `relation_infra_rs_acc_port_p` - (Optional) Relation to class infraAccPortP. Cardinality - N_TO_M. Type - Set of String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Profile.

## Importing

An existing Leaf Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_leaf_profile.example <Dn>
```
