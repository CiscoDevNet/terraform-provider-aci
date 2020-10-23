---
layout: "aci"
page_title: "ACI: aci_access_port_block"
sidebar_current: "docs-aci-data-source-access_port_block"
description: |-
  Data source for ACI Access Port Block
---

# aci_access_port_block #
Data source for ACI Access Port Block

## Example Usage ##

```hcl

data "aci_access_port_block" "dev_port_blk" {
  access_port_selector_dn  = "${aci_access_port_selector.example.id}"
  name                     = "foo_port_blk"
}

```


## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) name of Object access_port_block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Port Block.
* `annotation` - (Optional) annotation for object access_port_block.
* `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf access port block.
* `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf access port block.
* `name_alias` - (Optional) name_alias for object access_port_block.
* `to_card` - (Optional) The end (to-range) of the card range block for the leaf access port block.
* `to_port` - (Optional) The end (to-range) of the port range block for the leaf access port block.
