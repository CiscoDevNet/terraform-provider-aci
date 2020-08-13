---
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
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object spine_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Profile.
* `annotation` - (Optional) annotation for object spine_profile.
* `name_alias` - (Optional) name_alias for object spine_profile.
