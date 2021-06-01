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
  name        = "example"
  description = "from terraform"
  alloc_mode  = "static"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object VSAN Pool.
* `description` - (Optional) Description for object VSAN Pool.
* `alloc_mode` - (Required) Allocation mode for object VSAN Pool.
* `annotation` - (Optional) Annotation for object VSAN Pool.
* `name_alias` - (Optional) Name alias for object VSAN Pool.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VSAN Pool.

## Importing ##

An existing VSAN Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vsan_pool.example <Dn>
```
