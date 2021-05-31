---
layout: "aci"
page_title: "ACI: aci_firmware_policy"
sidebar_current: "docs-aci-resource-firmware_policy"
description: |-
  Manages ACI Firmware Policy
---

# aci_firmware_policy

Manages ACI Firmware Policy

## Example Usage

```hcl
resource "aci_firmware_policy" "example" {
  name  = "example"
  description = "from terraform"
  annotation  = "example"
  effective_on_reboot  = "no"
  ignore_compat  = "no"
  internal_label  = "example_policy"
  name_alias  = "example"
  version  = "n9000-14.2(3q)"
  version_check_override  = "untriggered"
}
```

## Argument Reference

- `name` - (Required) The firmware policy name.
- `annotation` - (Optional) Specifies a annotation of the policy definition.
- `description` - (Optional) Specifies a description of the policy definition.
- `effective_on_reboot` - (Optional) A property that indicates if the selected firmware version will be active after reboot. The firmware must be effective on an unplanned reboot before the scheduled maintenance operation. Allowed values: "no", "yes". Default value: "no".
- `ignore_compat` - (Optional) A property for specifying whether compatibility checks should be ignored when applying the firmware policy. Allowed values: "no", "yes". Default value: "no".
- `internal_label` - (Optional) Specifies a firmware of the policy definition. 
- `name_alias` - (Optional) Name alias of the policy definition. 
- `version` - (Optional) The firmware version.
- `version_check_override` - (Optional) The version check override.
  Allowed values: "trigger-immediate", "trigger", "triggered", "untriggered". Default value: "untriggered".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Firmware Policy.

## Importing

An existing Firmware Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_firmware_policy.example <Dn>
```
