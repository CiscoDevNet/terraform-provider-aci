---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_l2_interface_policy"
sidebar_current: "docs-aci-data-source-l2_interface_policy"
description: |-
  Data source for ACI L2 Interface Policy
---

# aci_l2_interface_policy

Data source for ACI L2 Interface Policy

## Example Usage

```hcl
data "aci_l2_interface_policy" "dev_l2_int_pol" {
  name  = "foo_l2_int_pol"
}
```

## Argument Reference

- `name` - (Required) Name of Object L2 interface policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L2 Interface Policy.
- `annotation` - (Optional) Annotation for object L2 interface policy.
- `description` - (Optional) Description for object L2 interface policy.
- `name_alias` - (Optional) Name alias for object L2 interface policy.
- `qinq` - (Optional) Determines if QinQ is disabled or if the port should be considered a core or edge port.
- `vepa` - (Optional) Determines if Virtual Ethernet Port Aggregator is disabled or enabled.
- `vlan_scope` - (Optional) The scope of the VLAN.
