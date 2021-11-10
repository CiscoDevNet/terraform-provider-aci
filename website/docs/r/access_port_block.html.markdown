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

## Example Usage

```hcl
	resource "aci_access_port_block" "fooaccess_port_block" {
		access_port_selector_dn = aci_access_port_selector.example.id
		description             = "from terraform"
		name                    = "demo_port_block"
		annotation              = "tag_port_block"
		from_card               = "1"
		from_port               = "1"
		name_alias              = "alias_port_block"
		to_card                 = "3"
		to_port                 = "3"
	}
```

## Argument Reference

- `access_port_selector_dn` - (Required) Distinguished name of parent Access Port Selector object.
- `name` - (Optional) name of Object Access Port Block.
- `annotation` - (Optional) Annotation for object Access Port Block.
- `description` - (Optional) Description for object Access Port Block.
- `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf access port block. Allowed value range is 1-100. Default value is "1".
- `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf access port block. Allowed value range is 1-127. Default value is "1".
- `name_alias` - (Optional) Name alias for object Access Port Block.
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
