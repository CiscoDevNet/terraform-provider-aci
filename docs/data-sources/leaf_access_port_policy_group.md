---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_access_port_policy_group"
sidebar_current: "docs-aci-data-source-leaf_access_port_policy_group"
description: |-
  Data source for ACI Leaf Access Port Policy Group
---

# aci_leaf_access_port_policy_group #
Data source for ACI Leaf Access Port Policy Group

## Example Usage ##

```hcl
data "aci_leaf_access_port_policy_group" "dev_leaf_port" {
  name  = "foo_leaf_port"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object leaf_access_port_policy_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Leaf Access Port Policy Group.
* `description` - (Optional) Description for object leaf access port policy group.
* `annotation` - (Optional) Annotation for object leaf access port policy group.
* `name_alias` - (Optional) Name alias for object leaf access port policy group.
