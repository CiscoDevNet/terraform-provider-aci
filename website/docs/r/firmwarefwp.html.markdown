---
layout: "aci"
page_title: "ACI: aci_firmware_policy"
sidebar_current: "docs-aci-resource-firmware_policy"
description: |-
  Manages ACI Firmware Policy
---

# aci_firmware_policy #
Manages ACI Firmware Policy

## Example Usage ##

```hcl
resource "aci_firmware_policy" "example" {


  name  = "example"
  annotation  = "example"
  effective_on_reboot  = "example"
  ignore_compat  = "example"
  internal_label  = "example"
  name_alias  = "example"
  version  = "example"
  version_check_override  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_policy.
* `annotation` - (Optional) annotation for object firmware_policy.
* `effective_on_reboot` - (Optional) firmware version effective on reboot selection
* `ignore_compat` - (Optional) whether compatibility check required
* `internal_label` - (Optional) firmware label
* `name_alias` - (Optional) name_alias for object firmware_policy.
* `version` - (Optional) firmware version
* `version_check_override` - (Optional) version check override



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Firmware Policy.

## Importing ##

An existing Firmware Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_firmware_policy.example <Dn>
```