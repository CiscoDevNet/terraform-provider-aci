---
layout: "aci"
page_title: "ACI: aci_leaf_interface_profile"
sidebar_current: "docs-aci-resource-leaf_interface_profile"
description: |-
  Manages ACI Leaf Interface Profile
---

# aci_leaf_interface_profile #
Manages ACI Leaf Interface Profile

## Example Usage ##

```hcl
	resource "aci_leaf_interface_profile" "fooleaf_interface_profile" {
		description = "%s"
		name        = "demo_leaf_profile"
		annotation  = "tag_leaf"
		name_alias  = "%s"
	}
```
## Argument Reference ##
* `name` - (Required) Name of Object leaf_interface_profile.
* `description` - (Optional) Description for object leaf_interface_profile.
* `annotation` - (Optional) Annotation for object leaf_interface_profile.
* `name_alias` - (Optional) Name alias for object leaf_interface_profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Interface Profile.

## Importing ##

An existing Leaf Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_interface_profile.example <Dn>
```