---
layout: "aci"
page_title: "ACI: aci_spine_port_policy_group"
sidebar_current: "docs-aci-data-source-spine_port_policy_group"
description: |-
  Data source for ACI Spine Port Policy Group
---

# aci_spine_port_policy_group #
Data source for ACI Spine Port Policy Group

## Example Usage ##

```hcl

data "aci_spine_port_policy_group" "example" {
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object aci_spine_port_policy_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Port Policy Group.
* `annotation` - (Optional) annotation for object aci_spine_port_policy_group.
* `name_alias` - (Optional) name_alias for object aci_spine_port_policy_group.
