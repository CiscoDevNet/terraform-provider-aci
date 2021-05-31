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
  access_port_selector_dn  = aci_access_port_selector.example.id
  name  = "example"
}
```
## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `name` - (Required) Name of Object access sub port block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Sub Port Block.
* `annotation` - (Optional) Annotation for object access sub port block.
* `from_card` - (Optional) From card
* `from_port` - (Optional) Port block from port
* `from_sub_port` - (Optional) From sub port for object access sub port block.
* `name_alias` - (Optional) Name alias for object access sub port block.
* `to_card` - (Optional) To card
* `to_port` - (Optional) To port
* `to_sub_port` - (Optional) To sub port for object access sub port block.
