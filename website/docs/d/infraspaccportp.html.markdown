---
layout: "aci"
page_title: "ACI: aci_spine_interface_profile"
sidebar_current: "docs-aci-data-source-spine_interface_profile"
description: |-
  Data source for ACI Spine Interface Profile
---

# aci_spine_interface_profile #
Data source for ACI Spine Interface Profile

## Example Usage ##

```hcl

data "aci_spine_interface_profile" "example" {
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object spine_interface_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Interface Profile.
* `annotation` - (Optional) annotation for object spine_interface_profile.
* `name_alias` - (Optional) name_alias for object spine_interface_profile.
