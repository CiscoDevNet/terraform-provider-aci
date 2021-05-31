---
layout: "aci"
page_title: "ACI: aci_vxlan_pool"
sidebar_current: "docs-aci-resource-vxlan_pool"
description: |-
  Manages ACI VXLAN Pool
---

# aci_vxlan_pool #
Manages ACI VXLAN Pool

## Example Usage ##

```hcl
resource "aci_vxlan_pool" "example" {
  name  = "example"
  annotation  = "example"
  name_alias  = "name_alias_for_vxlan"
  description = "From Teraform"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object vxlan pool.
* `annotation` - (Optional) Annotation for object vxlan pool.
* `description` - (Optional) Description for object vxlan pool.
* `name_alias` - (Optional) Name alias for object vxlan pool.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VXLAN Pool.

## Importing ##

An existing VXLAN Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vxlan_pool.example <Dn>
```