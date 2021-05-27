---
layout: "aci"
page_title: "ACI: aci_access_port_block"
sidebar_current: "docs-aci-data-source-access_port_block"
description: |-
  Data source for ACI Access Port Block
---

# aci_access_port_block

Data source for ACI Access Port Block

## Example Usage

```hcl

data "aci_access_port_block" "dev_port_blk" {
  access_port_selector_dn  = aci_access_port_selector.example.id
  name                     = "foo_port_blk"
}

```

## Argument Reference

- `access_port_selector_dn` - (Required) Distinguished name of parent Access Port Selector object.
- `name` - (Required) Name of Object Access Port Block.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Access Port Block.
- `description` - (Optional) Description for object Access Port Block.
- `annotation` - (Optional) Annotation for object Access Port Block.
- `from_card` - (Optional) The beginning (from-range) of the card range block for the leaf Access Port Block.
- `from_port` - (Optional) The beginning (from-range) of the port range block for the leaf Access Port Block.
- `name_alias` - (Optional) Name alias for object Access Port Block.
- `to_card` - (Optional) The end (to-range) of the card range block for the leaf Access Port Block.
- `to_port` - (Optional) The end (to-range) of the port range block for the leaf Access Port Block.
