---
layout: "aci"
page_title: "ACI: aci_firmware_policy"
sidebar_current: "docs-aci-data-source-firmware_policy"
description: |-
  Data source for ACI Firmware Policy
---

# aci_firmware_policy #
Data source for ACI Firmware Policy

## Example Usage ##

```hcl
data "aci_firmware_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Firmware Policy.
* `annotation` - (Optional) annotation for object firmware_policy.
* `effective_on_reboot` - (Optional) firmware version effective on reboot selection
* `ignore_compat` - (Optional) whether compatibility check required
* `internal_label` - (Optional) firmware label
* `name_alias` - (Optional) name_alias for object firmware_policy.
* `version` - (Optional) firmware version
* `version_check_override` - (Optional) version check override
