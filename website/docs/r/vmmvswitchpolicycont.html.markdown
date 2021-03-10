---
layout: "aci"
page_title: "ACI: aci_v_switch_policy_group"
sidebar_current: "docs-aci-resource-v_switch_policy_group"
description: |-
  Manages ACI VSwitch Policy Group
---

# aci_v_switch_policy_group #
Manages ACI VSwitch Policy Group

## Example Usage ##

```hcl
resource "aci_v_switch_policy_group" "example" {

  vmm_domain_dn  = "${aci_vmm_domain.example.id}"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.
* `annotation` - (Optional) annotation for object v_switch_policy_group.
* `name_alias` - (Optional) name_alias for object v_switch_policy_group.
