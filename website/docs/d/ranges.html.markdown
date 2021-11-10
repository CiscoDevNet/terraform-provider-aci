---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_ranges"
sidebar_current: "docs-aci-data-source-ranges"
description: |-
  Data source for ACI VLAN Pool Ranges
---

# aci_ranges #
Data source for ACI Ranges

## Example Usage ##

```hcl

data "aci_ranges" "example" {
  vlan_pool_dn  = aci_vlan_pool.example.id
  from  = "vlan-1"
  to  = "vlan-2"
}

```

## Argument Reference ##

* `vlan_pool_dn` - (Required) Distinguished name of parent VLAN Pool object.
* `from` - (Required) from of Object ranges.
* `to` - (Required) to of Object ranges.


## Attribute Reference

* `id` - Attribute id set to the Dn of the Ranges.
* `alloc_mode` - (Optional) alloc_mode for object ranges.
* `annotation` - (Optional) Annotation for object ranges.
* `description` - (Optional) Description for object ranges.
* `name_alias` - (Optional) Name alias for object ranges.
* `role` - (Optional) System role type
