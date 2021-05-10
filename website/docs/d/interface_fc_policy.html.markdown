---
layout: "aci"
page_title: "ACI: aci_interface_fc_policy"
sidebar_current: "docs-aci-data-source-interface_fc_policy"
description: |-
  Data source for ACI Interface FC Policy
---

# aci_interface_fc_policy

Data source for ACI Interface FC Policy

## Example Usage

```hcl
data "aci_interface_fc_policy" "test_pol" {
  name  = "demo_int_policy"
}
```

## Argument Reference

- `name` - (Required) Name of Object interface_fc_policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Interface FC Policy.
- `annotation` - (Optional) annotation for object interface_fc policy.
- `automaxspeed` - (Optional) automaxspeed for object interface fc policy.
- `fill_pattern` - (Optional) Fill Pattern for native FC ports.
- `name_alias` - (Optional) name_alias for object interface_fc_policy.
- `port_mode` - (Optional) In which mode Ports should be used.
- `rx_bb_credit` - (Optional) Receive buffer credits for native FC ports.
- `speed` - (Optional) cpu or port speed.
- `trunk_mode` - (Optional) Trunking on/off for native FC ports.Default value is OFF.
