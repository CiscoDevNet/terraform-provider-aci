---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_profile"
sidebar_current: "docs-aci-data-source-leaf_profile"
description: |-
  Data source for ACI Leaf Profile
---

# aci_leaf_profile #
Data source for ACI Leaf Profile

## Example Usage ##

```hcl
data "aci_leaf_profile" "dev_leaf_prof" {
  name  = "foo_leaf_prof"
}
```
## Argument Reference ##
* `name` - (Required) name of Object leaf_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Leaf Profile.
* `description` - (Optional) Description for object leaf profile.
* `annotation` - (Optional) Annotation for object leaf profile.
* `name_alias` - (Optional) Name alias for object leaf profile.
