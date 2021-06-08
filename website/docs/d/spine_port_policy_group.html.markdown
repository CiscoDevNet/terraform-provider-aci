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
  name  = "spine_port_policy_group_1"
}

```


## Argument Reference ##
* `name` - (Required) Name of Object Spine Port Policy Group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Port Policy Group.
* `description` - (Optional) Description for object Spine Port Policy Group.
* `annotation` - (Optional) Annotation for object Spine Port Policy Group.
* `name_alias` - (Optional) Name alias for object Spine Port Policy Group.
