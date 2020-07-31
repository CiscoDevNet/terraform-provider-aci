---
layout: "aci"
page_title: "ACI: aci_access_group"
sidebar_current: "docs-aci-resource-access_group"
description: |-
  Manages ACI Access Group
---

# aci_access_group #
Manages ACI Access Group

## Example Usage ##

```hcl

resource "aci_access_group" "example" {
  access_port_selector_dn   = "${aci_access_port_selector.example.id}"
  annotation                = "one"
  fex_id                    = "101"
  tdn                       = "${aci_fex_bundle_group.example.id}"
}

```


## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.
* `annotation` - (Optional) annotation for object access_access_group.
* `fex_id` - (Optional) interface policy group fex id
* `tdn` - (Optional) interface policy group's target rn



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Access Group.

## Importing ##

An existing Access Access Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_group.example <Dn>
```