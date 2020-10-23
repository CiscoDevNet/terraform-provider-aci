---
layout: "aci"
page_title: "ACI: aci_l2_interface_policy"
sidebar_current: "docs-aci-data-source-l2_interface_policy"
description: |-
  Data source for ACI L2 Interface Policy
---

# aci_l2_interface_policy #
Data source for ACI L2 Interface Policy

## Example Usage ##

```hcl
data "aci_l2_interface_policy" "dev_l2_int_pol" {
  name  = "foo_l2_int_pol"
}
```
## Argument Reference ##
* `name` - (Required) name of Object l2_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L2 Interface Policy.
* `annotation` - (Optional) annotation for object l2_interface_policy.
* `name_alias` - (Optional) name_alias for object l2_interface_policy.
* `qinq` - (Optional) Determines if QinQ is disabled or if the port should be considered a core or edge port.
* `vepa` - (Optional) Determines if Virtual Ethernet Port Aggregator is disabled or enabled.
* `vlan_scope` - (Optional) The scope of the VLAN.
