---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_access_port_block"
sidebar_current: "docs-aci-data-source-access_port_block"
description: |-
  Data source for ACI Access Port Block
---

# aci_access_port_block

Data source for ACI Access Port Block

## API Information ##

* `Class` - infraPortBlk
* `Distinguished Name` - uni/infra/accportprof-{leaf_interface_profile_name}/hports-{leaf_interface_profile_dn}-typ-{type}/portblk-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Interfaces -> Profiles -> Create Leaf Interface Profile -> Interface Selectors -> Port Blocks

## Example Usage

```hcl

data "aci_access_port_block" "dev_port_blk" {
  access_port_selector_dn  = aci_access_port_selector.example.id
  name                     = "foo_port_blk"
}

```

## Argument Reference

- `access_port_selector_dn` - (Required) Distinguished name of the parent Access Port Selector or Spine Access Port Selector object.
- `name` - (Required) Name of the Access Port Block object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Access Port Block.
- `description` - (Optional) Description of the Access Port Block object.
- `annotation` - (Optional) Annotation of the Access Port Block object.
- `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf Access Port Block.
- `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf Access Port Block.
- `name_alias` - (Optional) Name alias of the Access Port Block object.
- `to_card` - (Optional) The end (to-range) of the card range block for the leaf Access Port Block.
- `to_port` - (Optional) The end (to-range) of the port range block for the leaf Access Port Block.

- `relation_infra_rs_acc_bndl_subgrp` - (Optional) Relation to class infraAccBndlSubgrp. Cardinality - N_TO_ONE. Type - String.
