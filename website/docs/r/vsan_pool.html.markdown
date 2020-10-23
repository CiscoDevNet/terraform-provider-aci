---
layout: "aci"
page_title: "ACI: aci_vsan_pool"
sidebar_current: "docs-aci-resource-vsan_pool"
description: |-
  Manages ACI VSAN Pool
---

# aci_vsan_pool #
Manages ACI VSAN Pool

## Example Usage ##

```hcl
resource "aci_vsan_pool" "example" {


  name  = "example"

  alloc_mode  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object vsan_pool.
* `alloc_mode` - (Optional) alloc_mode for object vsan_pool.
Allowed values: "dynamic", "static"
* `annotation` - (Optional) annotation for object vsan_pool.
* `name_alias` - (Optional) name_alias for object vsan_pool.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VSAN Pool.

## Importing ##

An existing VSAN Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vsan_pool.example <Dn>
```