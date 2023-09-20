---
subcategory: "Firmware"
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
* `name` - (Required) Name of Object firmware_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Firmware Group.
* `annotation` - (Optional) Annotation for object Firmware Group.
* `description` - (Optional) Description for object Firmware Group.
* `name_alias` - (Optional) Name alias for object Firmware Group.
* `firmware_group_type` - (Optional) Component type.
