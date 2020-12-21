---
layout: "aci"
page_title: "ACI: aci_ranges"
sidebar_current: "docs-aci-data-source-ranges"
description: |-
  Data source for ACI Ranges
---

# aci_ranges #
Data source for ACI Ranges

## Example Usage ##

```hcl

data "aci_ranges" "example" {
  vlan_pool_dn  = "${aci_vlan_pool.example.id}"
  from  = "example"
  to  = "example"
}

```

## Argument Reference ##

* `vlan_pool_dn` - (Required) Distinguished name of parent VLANPool object.
* `from` - (Required) from of Object ranges.
* `to` - (Required) to of Object ranges.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Ranges.
* `alloc_mode` - (Optional) alloc_mode for object ranges.
* `annotation` - (Optional) annotation for object ranges.
* `name_alias` - (Optional) name_alias for object ranges.
* `role` - (Optional) system role type
