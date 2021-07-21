---
layout: "aci"
page_title: "ACI: aci_ranges"
sidebar_current: "docs-aci-resource-ranges"
description: |-
  Manages ACI Ranges
---

# aci_ranges #

Manages ACI Ranges

## API Information ##

* `Class` - fvnsEncapBlk
* `Distinguished Named` - uni/infra/vlanns-[{name}]-dynamic/from-[{vlan-num}]-to-[{vlan-num}]

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> VLAN -> Encap Blocks

## Example Usage ##

```hcl
resource "aci_ranges" "example" {
  vlan_pool_dn  = aci_vlan_pool.example.id
  description   = "pool range description"
  from          = "vlan-1"
  to            = "vlan-10"
  alloc_mode    = "example"
  annotation    = "example"
  name_alias    = "example"
  role          = "external"
}
```

## Argument Reference ##

* `vlan_pool_dn` - (Required) Distinguished name of parent VLANPool object.
* `from` - (Required) _from of Object ranges.
* `to` - (Required) to of Object ranges.
* `alloc_mode` - (Optional) allocation modes for object ranges.  Allowed values: "dynamic", "static", "inherit".
* `annotation` - (Optional) annotation for object ranges.
* `name_alias` - (Optional) name_alias for object ranges.
* `role` - (Optional) system role type.  Allowed values: "external", "internal".  Default is "external".
* `description` - (Optional) description for object ranges.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Ranges.

## Importing ##

An existing Ranges can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_ranges.example <Dn>
```
