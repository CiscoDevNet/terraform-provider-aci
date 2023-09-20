---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_profile"
sidebar_current: "docs-aci-data-source-spine_profile"
description: |-
  Data source for ACI Spine Profile
---

# aci_spine_profile #
Data source for ACI Spine Profile

## Example Usage ##

```hcl

data "aci_spine_profile" "example" {
  name  = "spine_profile_1"
}

```


## Argument Reference ##
* `name` - (Required) Name of Object Spine Profile.


## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Profile.
* `description` - (Optional) Description for object Spine Profile.
* `annotation` - (Optional) Annotation for object Spine Profile.
* `name_alias` - (Optional) Name alias for object Spine Profile.
