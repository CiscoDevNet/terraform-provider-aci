---
layout: "aci"
page_title: "ACI: aci_access_sub_port_block"
sidebar_current: "docs-aci-data-source-access_sub_port_block"
description: |-
  Data source for ACI Access Sub Port Block
---

# aci_access_sub_port_block #
Data source for ACI Access Sub Port Block

## Example Usage ##

```hcl
data "aci_access_sub_port_block" "example" {

  access_port_selector_dn  = "${aci_access_port_selector.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) name of Object access_sub_port_block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Sub Port Block.
* `annotation` - (Optional) annotation for object access_sub_port_block.
* `from_card` - (Optional) from card
* `from_port` - (Optional) port block from port
* `from_sub_port` - (Optional) from_sub_port for object access_sub_port_block.
* `name_alias` - (Optional) name_alias for object access_sub_port_block.
* `to_card` - (Optional) to card
* `to_port` - (Optional) to port
* `to_sub_port` - (Optional) to_sub_port for object access_sub_port_block.
