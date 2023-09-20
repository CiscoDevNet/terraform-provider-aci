---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_vlan_pool"
sidebar_current: "docs-aci-data-source-vlan_pool"
description: |-
  Data source for ACI VLAN Pool
---

# aci_vlan_pool

Data source for ACI VLAN Pool

## Example Usage

```hcl
data "aci_vlan_pool" "example" {
  name  = "example"
  alloc_mode  = "dynamic"
}
```

## Argument Reference

- `name` - (Required) name of Object vlan_pool.
- `alloc_mode` - (Required) allocation mode of Object vlan_pool.

## Attribute Reference

- `id` - Attribute id set to the Dn of the VLAN Pool.
- `alloc_mode` - (Optional) Allocation mode for object vlan pool.
- `annotation` - (Optional) Annotation for object vlan pool.
- `description` - (Optional) Description for object vlan pool.
- `name_alias` - (Optional) Name alias for object vlan pool.
