---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_switch_policy_group"
sidebar_current: "docs-aci-data-source-spine_switch_policy_group"
description: |-
  Data source for ACI Spine Switch Policy Group
---

# aci_spine_switch_policy_group #
Data source for ACI Spine Switch Policy Group


## API Information ##
* `Class` - infraSpineAccNodePGrp
* `Distinguished Name` - uni/infra/funcprof/spaccnodepgrp-{name}

## GUI Information ##
* `Location` - Fabric -> Access Policies -> Switches -> Spine Switches -> Policy Groups -> Create Spine Switch Policy Group

## Example Usage ##
```hcl
data "aci_spine_switch_policy_group" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object Spine Switch Policy Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Spine Switch Policy Group.
* `annotation` - (Optional) Annotation of object Spine Switch Policy Group.
* `name_alias` - (Optional) Name Alias of object Spine Switch Policy Group.
* `description` - (Optional) Description for object Spine Switch Policy Group.
