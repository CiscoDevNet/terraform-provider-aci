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
  vlan_pool_dn  = aci_vlan_pool.example.id
  from          = "vlan-1"
  description   = "From Terraform"
  to            = "vlan-2"
  alloc_mode    = "inherit"
  annotation    = "example"
  name_alias    = "name_alias"
  role          = "external"
}
```

## Argument Reference ##

* `vlan_pool_dn` - (Required) Distinguished name of parent VLANPool object.
* `from` - (Required) From of Object ranges. Allowed value min: vlan-1, max: vlan-4094.
* `to` - (Required) To of Object ranges. Allowed value min: vlan-1, max: vlan-4094.
* `alloc_mode` - (Optional) Alloc mode for object ranges.  Allowed values: "dynamic", "static", "inherit". Default is "inherit".
* `annotation` - (Optional) Annotation for object ranges.
* `name_alias` - (Optional) Name alias for object ranges.
* `description` - (Optional) Description for object ranges.
* `role` - (Optional) System role type.  Allowed values: "external", "internal".  Default is "external".

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Ranges.

## Importing ##

An existing Ranges can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_ranges.example <Dn>
```
