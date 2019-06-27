---
layout: "aci"
page_title: "ACI: aci_lldp_interface_policy"
sidebar_current: "docs-aci-data-source-lldp_interface_policy"
description: |-
  Data source for ACI LLDP Interface Policy
---

# aci_lldp_interface_policy #
Data source for ACI LLDP Interface Policy

## Example Usage ##

```hcl
data "aci_lldp_interface_policy" "dev_lldp_pol" {
  name  = "foo_lldp_pol"
}
```
## Argument Reference ##
* `name` - (Required) name of Object lldp_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the LLDP Interface Policy.
* `admin_rx_st` - (Optional) admin receive state.
* `admin_tx_st` - (Optional) admin transmit state.
* `annotation` - (Optional) annotation for object lldp_interface_policy.
* `name_alias` - (Optional) name_alias for object lldp_interface_policy.
