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
* `name` - (Required) Name of Object LLDP Interface Policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the LLDP Interface Policy.
* `admin_rx_st` - (Optional) Admin receive state.
* `admin_tx_st` - (Optional) Admin transmit state.
* `annotation` - (Optional) Annotation for object LLDP Interface Policy.
* `name_alias` - (Optional) Name alias for object LLDP Interface Policy.
