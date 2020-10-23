---
layout: "aci"
page_title: "ACI: aci_spine_interface_profile"
sidebar_current: "docs-aci-resource-spine_interface_profile"
description: |-
  Manages ACI Spine Interface Profile
---

# aci_spine_interface_profile #
Manages ACI Spine Interface Profile

## Example Usage ##

```hcl

resource "aci_spine_interface_profile" "example" {
  name        = "example"
  annotation  = "example"
  name_alias  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object spine_interface_profile.
* `annotation` - (Optional) annotation for object spine_interface_profile.
* `name_alias` - (Optional) name_alias for object spine_interface_profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Interface Profile.

## Importing ##

An existing Spine Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_interface_profile.example <Dn>
```