---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_interface_profile"
sidebar_current: "docs-aci-resource-aci_leaf_interface_profile"
description: |-
  Manages ACI Leaf Interface Profile
---

# aci_leaf_interface_profile #
Manages ACI Leaf Interface Profile

## Example Usage ##

```hcl
resource "aci_leaf_interface_profile" "fooleaf_interface_profile" {
	description = "From Terraform"
	name        = "demo_leaf_profile"
	annotation  = "tag_leaf"
	name_alias  = "name_alias"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object leaf interface profile.
* `description` - (Optional) Description for object leaf interface profile.
* `annotation` - (Optional) Annotation for object leaf interface profile.
* `name_alias` - (Optional) Name alias for object leaf interface profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Interface Profile.

## Importing ##

An existing Leaf Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_interface_profile.example <Dn>
```
