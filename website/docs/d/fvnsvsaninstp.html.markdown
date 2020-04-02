---
layout: "aci"
page_title: "ACI: aci_vsan_pool"
sidebar_current: "docs-aci-data-source-vsan_pool"
description: |-
  Data source for ACI VSAN Pool
---

# aci_vsan_pool #
Data source for ACI VSAN Pool

## Example Usage ##

```hcl
data "aci_vsan_pool" "example" {


  name  = "example"

  allocMode  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object vsan_pool.
* `allocMode` - (Required) allocMode of Object vsan_pool.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VSAN Pool.
* `alloc_mode` - (Optional) alloc_mode for object vsan_pool.
* `annotation` - (Optional) annotation for object vsan_pool.
* `name_alias` - (Optional) name_alias for object vsan_pool.
