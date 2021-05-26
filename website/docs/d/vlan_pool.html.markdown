---
layout: "aci"
page_title: "ACI: aci_vlan_pool"
sidebar_current: "docs-aci-data-source-vlan_pool"
description: |-
  Data source for ACI VLAN Pool
---

# aci_vlan_pool #
Data source for ACI VLAN Pool

## Example Usage ##

```hcl
data "aci_vlan_pool" "example" {
  name  = "example"
  allocMode  = "example"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object vlan pool.
* `allocMode` - (Required) AllocMode of Object vlan pool.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VLAN Pool.
* `alloc_mode` - (Optional) Allocation mode for object vlan_pool.
* `annotation` - (Optional) Annotation for object vlan_pool.
* `description` - (Optional) Description for  object vlan pool.
* `name_alias` - (Optional) Name alias for object vlan_pool.
