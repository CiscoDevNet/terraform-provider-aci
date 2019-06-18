---
layout: "aci"
page_title: "ACI: aci_access_port_block"
sidebar_current: "docs-aci-resource-access_port_block"
description: |-
  Manages ACI Access Port Block
---

# aci_access_port_block #
Manages ACI Access Port Block

## Example Usage ##

```hcl
	resource "aci_access_port_block" "fooaccess_port_block" {
		access_port_selector_dn = "${aci_access_port_selector.example.id}"
		description             = "%s"
		name                    = "demo_port_block"
		annotation              = "tag_port_block"
		from_card               = "1"
		from_port               = "1"
		name_alias              = "alias_port_block"
		to_card                 = "3"
		to_port                 = "3"
	}
```
## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) name of Object access_port_block.
* `annotation` - (Optional) annotation for object access_port_block.
* `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf access port block.
* `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf access port block.
* `name_alias` - (Optional) name_alias for object access_port_block.
* `to_card` - (Optional) The end (to-range) of the card range block for the leaf access port block.
* `to_port` - (Optional) The end (to-range) of the port range block for the leaf access port block.

* `relation_infra_rs_acc_bndl_subgrp` - (Optional) Relation to class infraAccBndlSubgrp. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Port Block.

## Importing ##

An existing Access Port Block can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_port_block.example <Dn>
```