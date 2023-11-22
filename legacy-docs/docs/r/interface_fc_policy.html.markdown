---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_interface_fc_policy"
sidebar_current: "docs-aci-resource-interface_fc_policy"
description: |-
  Manages ACI Interface FC Policy
---

# aci_interface_fc_policy

Manages ACI Interface FC Policy

## Example Usage

```hcl
resource "aci_interface_fc_policy" "example" {
  name         = "demo_policy"
  annotation   = "tag_if_policy"
  description  = "from terraform"
  automaxspeed = "32G"
  fill_pattern = "IDLE"
  name_alias   = "demo_alias"
  port_mode    = "f"
  rx_bb_credit = "64"
  speed        = "auto"
  trunk_mode   = "trunk-off"
}

```

## Argument Reference

- `name` - (Required) Name of Object interface FC policy.
- `annotation` - (Optional) Annotation for object interface FC policy.
- `description` - (Optional) Description for object interface FC policy.
- `automaxspeed` - (Optional) Auto-max-speed for object interface FC policy. Allowed values are "2G", "4G", "8G", "16G" and "32G". Default value is "32G".
- `fill_pattern` - (Optional) Fill Pattern for native FC ports. Allowed values are "ARBFF" and "IDLE". Default is "IDLE".
- `name_alias` - (Optional) Name alias for object Interface FC policy.
- `port_mode` - (Optional) In which mode Ports should be used. Allowed values are "f" and "np". Default is "f".
- `rx_bb_credit` - (Optional) Receive buffer credits for native FC ports Range:(16 - 64). Default value is "64".
- `speed` - (Optional) CPU or port speed. All the supported values are "unknown", "auto", "4G", "8G", "16G", "32G". Default value is "auto".
- `trunk_mode` - (Optional) Trunking on/off for native FC ports. Allowed values are "un-init", "trunk-off", "trunk-on" and "auto". Default value is "trunk-off".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Interface FC Policy.

## Importing

An existing Interface FC Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_interface_fc_policy.example <Dn>
```
