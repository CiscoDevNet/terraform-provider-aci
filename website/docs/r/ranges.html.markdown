---
layout: "aci"
page_title: "ACI: aci_ranges"
sidebar_current: "docs-aci-resource-ranges"
description: |-
  Manages ACI Ranges
---

# aci_ranges #
Manages ACI Ranges

## Example Usage ##

```hcl
resource "aci_ranges" "example" {

  vlan_pool_dn  = "${aci_vlan_pool.example.id}"

  _from  = "example"

  to  = "example"
  alloc_mode  = "example"
  annotation  = "example"
  from  = "example"
  name_alias  = "example"
  role  = "example"
}
```
## Argument Reference ##
* `vlan_pool_dn` - (Required) Distinguished name of parent VLANPool object.
* `_from` - (Required) _from of Object ranges.
* `to` - (Required) to of Object ranges.
* `alloc_mode` - (Optional) alloc_mode for object ranges.
Allowed values: "dynamic", "static", "inherit"
* `annotation` - (Optional) annotation for object ranges.
* `from` - (Optional) encapsulation block start
* `name_alias` - (Optional) name_alias for object ranges.
* `role` - (Optional) system role type.
Allowed values: "external", "internal"
* `to` - (Optional) encapsulation block end



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Ranges.

## Importing ##

An existing Ranges can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ranges.example <Dn>
```