---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_access_port_block"
sidebar_current: "docs-aci-resource-access_port_block"
description: |-
  Manages ACI Access Port Block
---

# aci_access_port_block

Manages ACI Access Port Block

## API Information ##

* `Class` - infraPortBlk
* `Distinguished Name` - uni/infra/accportprof-{leaf_interface_profile_name}/hports-{leaf_interface_profile_dn}-typ-{type}/portblk-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Interfaces -> Profiles -> Create Leaf Interface Profile -> Interface Selectors -> Port Blocks

## Example Usage

```hcl
resource "aci_access_port_block" "test_port_block" {
  access_port_selector_dn           = aci_access_port_selector.fooaccess_port_selector.id
  name                              = "tf_test_block"
  description                       = "From Terraform"
  annotation                        = "tag_port_block"
  from_card                         = "1"
  from_port                         = "1"
  name_alias                        = "alias_port_block"
  to_card                           = "3"
  to_port                           = "3"
  relation_infra_rs_acc_bndl_subgrp = aci_leaf_access_bundle_policy_sub_group.test_access_bundle_policy_sub_group.id
}
```

## Argument Reference

- `access_port_selector_dn` - (Required) Distinguished name of the parent Access Port Selector or Spine Access Port Selector object.
- `name` - (Optional) Name of the Access Port Block object.
- `annotation` - (Optional) Annotation of the Access Port Block object.
- `description` - (Optional) Description of the Access Port Block object.
- `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf access port block. Allowed value range is 1-100. Default value is "1".
- `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf access port block. Allowed value range is 1-127. Default value is "1".
- `name_alias` - (Optional) Name alias of the Access Port Block object.
- `to_card` - (Optional) The end (to-range) of the card range block for the leaf access port block. Allowed value range is 1-100. Default value is "1".
- `to_port` - (Optional) The end (to-range) of the port range block for the leaf access port block. Allowed value range is 1-127. Default value is "1".

- `relation_infra_rs_acc_bndl_subgrp` - (Optional) Relation to class infraAccBndlSubgrp. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Port Block.

## Importing

An existing Access Port Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_access_port_block.example <Dn>
```
