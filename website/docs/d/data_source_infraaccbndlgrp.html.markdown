---
layout: "aci"
page_title: "ACI: aci_pcvpc_interface_policy_group"
sidebar_current: "docs-aci-data-source-pcvpc_interface_policy_group"
description: |-
  Data source for ACI PC/VPC Interface Policy Group
---

# aci_pcvpc_interface_policy_group #
Data source for ACI PC/VPC Interface Policy Group

## Example Usage ##

```hcl
data "aci_pcvpc_interface_policy_group" "dev_pol_grp" {
  name  = "foo_pol_grp"
}
```
## Argument Reference ##
* `name` - (Required) name of Object pcvpc_interface_policy_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the PC/VPC Interface Policy Group.
* `annotation` - (Optional) annotation for object pcvpc_interface_policy_group.
* `lag_t` - (Optional) The bundled ports group link aggregation type: port channel vs virtual port channel.
* `name_alias` - (Optional) name_alias for object pcvpc_interface_policy_group.
