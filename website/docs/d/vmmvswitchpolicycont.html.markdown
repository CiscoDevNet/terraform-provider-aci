---
layout: "aci"
page_title: "ACI: aci_v_switch_policy_group"
sidebar_current: "docs-aci-data-source-v_switch_policy_group"
description: |-
  Data source for ACI VSwitch Policy Group
---

# vSwitch_Policy #

Data source for ACI VSwitch Policy Group


## API Information ##

* `Class` - vmmVSwitchPolicyCont
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/vswitchpolcont

## GUI Information ##

* `Location` - Virtual Networking -> VMM Domain -> VSwitchPolicy



## Example Usage ##

```hcl
data "aci_v_switch_policy_group" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the VSwitch Policy Group.
* `annotation` - (Optional) Annotation of object VSwitch Policy Group.
* `name_alias` - (Optional) Name Alias of object VSwitch Policy Group.
