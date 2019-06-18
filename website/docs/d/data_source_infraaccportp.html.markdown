---
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
data "aci_leaf_interface_profile" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object leaf_interface_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Leaf Interface Profile.
* `annotation` - (Optional) annotation for object leaf_interface_profile.
* `name_alias` - (Optional) name_alias for object leaf_interface_profile.
