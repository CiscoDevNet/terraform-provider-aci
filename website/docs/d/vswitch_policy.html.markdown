---
layout: "aci"
page_title: "ACI: aci_vswitch_policy"
sidebar_current: "docs-aci-data-source-vswitch_policy"
description: |-
  Data source for ACI VSwitch Policy Group
---

# aci_vswitch_policy #

Data source for ACI VSwitch Policy Group

## API Information ##

* `Class` - vmmVSwitchPolicyCont
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/vswitchpolcont


## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> VSwitch Policy

## Example Usage ##

```hcl
data "aci_vswitch_policy" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the VSwitch Policy Group.
* `annotation` - (Optional) Annotation of object VSwitch Policy Group.
* `description` - (Optional) Description of object VSwitch Policy Group.
* `name_alias` - (Optional) Name Alias of object VSwitch Policy Group.
