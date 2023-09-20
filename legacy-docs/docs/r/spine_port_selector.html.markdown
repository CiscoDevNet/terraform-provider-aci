---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_port_selector"
sidebar_current: "docs-aci-resource-spine_port_selector"
description: |-
  Manages ACI Spine Port selector
---

# aci_spine_port_selector #
!> **WARNING:** This resource is deprecated and will be removed in the next major version use aci_spine_interface_profile_selector instead.
Manages ACI Spine Port selector

## Example Usage ##

```hcl

resource "aci_spine_port_selector" "example" {
  spine_profile_dn  = aci_spine_profile.example.id
  tdn               = aci_spine_interface_profile.example.id
  annotation        = "example"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent Spine Profile.
* `tdn` - (Required) tDn of the Spine Interface Profile.
* `annotation` - (Optional) Annotation for Spine Port selector.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Port selector.

## Importing ##

An existing Spine port selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_interface_profile_selector.example <Dn>
```