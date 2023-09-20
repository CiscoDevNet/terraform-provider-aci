---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_firmware_policy"
sidebar_current: "docs-aci-data-source-firmware_policy"
description: |-
  Data source for ACI Firmware Policy
---

# aci_firmware_policy

Data source for ACI Firmware Policy

## Example Usage

```hcl
data "aci_firmware_policy" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) The firmware policy name.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Firmware Policy.
- `annotation` - (Optional) Specifies the annotation of the policy definition.
- `description` - (Optional) Specifies the description of the policy definition.
- `effective_on_reboot` - (Optional) A property that indicates if the selected firmware version will be active after reboot. The firmware must be effective on an unplanned reboot before the scheduled maintenance operation.
- `ignore_compat` - (Optional) A property for specifying whether compatibility checks should be ignored when applying the firmware policy.
- `internal_label` - (Optional) Specifies a firmware of the policy definition.
- `name_alias` - (Optional) Name alias of the policy definition.
- `version` - (Optional) The firmware version.
- `version_check_override` - (Optional) The version check override.
