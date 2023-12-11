---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_ranges"
sidebar_current: "docs-aci-resource-aci_ranges"
description: |-
  Manages ACI VLAN Pool Ranges
---

# aci_ranges #

Manages ACI Ranges

## API Information ##

* `Class` - fvnsEncapBlk
* `Distinguished Name` - uni/infra/vlanns-[{name}]-dynamic/from-[{vlan-num}]-to-[{vlan-num}]

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> VLAN -> Encap Blocks

## Example Usage ##

```hcl
resource "aci_vlan_pool" "vlan_pool_1" {
  name       = "VLANPool1"
  alloc_mode = "static"
}

resource "aci_ranges" "range_1" {
  vlan_pool_dn  = aci_vlan_pool.vlan_pool_1.id
  description   = "From Terraform"
  from          = "vlan-1"
  to            = "vlan-2"
  alloc_mode    = "inherit"
  annotation    = "example"
  name_alias    = "name_alias"
  role          = "external"
}
```

## Argument Reference ##

* `vlan_pool_dn` - (Required) Distinguished name of parent VLAN Pool object (from aci_vlan_pool resource/data source).
* `from` - (Required) From of Object VLAN Pool ranges. Allowed value min: vlan-1, max: vlan-4094.
* `to` - (Required) To of Object VLAN Pool ranges. Allowed value min: vlan-1, max: vlan-4094.
* `alloc_mode` - (Optional) Alloc mode for object VLAN Pool ranges.  Allowed values: "dynamic", "static", "inherit". Default is "inherit".
* `annotation` - (Optional) Annotation for object VLAN Pool ranges.
* `name_alias` - (Optional) Name alias for object VLAN Pool ranges.
* `role` - (Optional) System role type.  Allowed values: "external", "internal".  Default is "external".
* `description` - (Optional) Description for object VLAN Pool ranges.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Ranges.

## Importing ##

An existing Ranges can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_ranges.example <Dn>
```
