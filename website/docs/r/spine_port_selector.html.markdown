---
layout: "aci"
page_title: "ACI: aci_spine_port_selector"
sidebar_current: "docs-aci-resource-spine_port_selector"
description: |-
  Manages ACI Spine port selector
---

# aci_spine_port_selector #
Manages ACI Spine port selector

## Example Usage ##

```hcl

resource "aci_spine_port_selector" "example" {
  spine_profile_dn  = "${aci_spine_profile.example.id}"
  tdn               = "${aci_spine_interface_profile.example.id}"
  annotation        = "example"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent SpineProfile object.
* `tdn` - (Required) tdn of Object Interface profile.
* `annotation` - (Optional) Annotation for object Spine port selector.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine port selector.

## Importing ##

An existing Spine port selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_port_selector.example <Dn>
```