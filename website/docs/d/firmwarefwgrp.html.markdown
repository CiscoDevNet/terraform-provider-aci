---
layout: "aci"
page_title: "ACI: aci_firmware_group"
sidebar_current: "docs-aci-data-source-firmware_group"
description: |-
  Data source for ACI Firmware Group
---

# aci_firmware_group #
Data source for ACI Firmware Group

## Example Usage ##

```hcl
data "aci_firmware_group" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Firmware Group.
* `annotation` - (Optional) annotation for object firmware_group.
* `name_alias` - (Optional) name_alias for object firmware_group.
* `firmware_group_type` - (Optional) component type
