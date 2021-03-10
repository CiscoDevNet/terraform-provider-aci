---
layout: "aci"
page_title: "ACI: aci_v_switch_policy_group"
sidebar_current: "docs-aci-data-source-v_switch_policy_group"
description: |-
  Data source for ACI VSwitch Policy Group
---

# aci_v_switch_policy_group #
Data source for ACI VSwitch Policy Group

## Example Usage ##

```hcl
data "aci_v_switch_policy_group" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VSwitch Policy Group.
* `annotation` - (Optional) annotation for object v_switch_policy_group.
* `name_alias` - (Optional) name_alias for object v_switch_policy_group.
