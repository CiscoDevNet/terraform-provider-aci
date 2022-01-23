---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_interface_profile_selector"
sidebar_current: "docs-aci-resource-spine_interface_profile_selector"
description: |-
  Manages ACI Spine Interface Profile selector
---

# aci_spine_interface_profile_selector #
Manages ACI Spine Interface Profile selector

## Example Usage ##

```hcl

resource "aci_spine_interface_profile_selector" "example" {
  spine_profile_dn  = aci_spine_profile.example.id
  tdn               = aci_spine_interface_profile.example.id
  annotation        = "example"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent Spine Profile.
* `tdn` - (Required) tDn of the Spine Interface Profile.
* `annotation` - (Optional) Annotation for Spine Interface Profile selector.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Interface Profile selector.

## Importing ##

An existing Spine Interface Profile selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_interface_profile_selector.example <Dn>
```