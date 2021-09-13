---
layout: "aci"
page_title: "ACI: aci_access_switch_policy_group"
sidebar_current: "docs-aci-data-source-access_switch_policy_group"
description: |-
  Data source for ACI Access Switch Policy Group
---

# aci_access_switch_policy_group #

Data source for ACI Access Switch Policy Group


## API Information ##

* `Class` - infraAccNodePGrp
* `Distinguished Named` - uni/infra/funcprof/accnodepgrp-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Switches -> Leaf Switches -> Policy Groups -> Create Access Switch Policy Group



## Example Usage ##

```hcl
data "aci_access_switch_policy_group" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object Access Switch Policy Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Access Switch Policy Group.
* `annotation` - (Optional) Annotation of object Access Switch Policy Group.
* `name_alias` - (Optional) Name Alias of object Access Switch Policy Group.
* `description` - (Optional) Description for object Access Switch Policy Group.
