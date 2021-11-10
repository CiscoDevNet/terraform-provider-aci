---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_interface_profile"
sidebar_current: "docs-aci-data-source-leaf_interface_profile"
description: |-
  Data source for ACI Leaf Interface Profile
---

# aci_leaf_interface_profile #
Data source for ACI Leaf Interface Profile

## Example Usage ##

```hcl
data "aci_leaf_interface_profile" "dev_leaf_int_prof" {
  name  = "foo_leaf_int_prof"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object Leaf Interface Profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Leaf Interface Profile.
* `description` - (Optional) Description for object Leaf Interface Profile.
* `annotation` - (Optional) Annotation for object Leaf Interface Profile.
* `name_alias` - (Optional) Name alias for object Leaf Interface Profile.
