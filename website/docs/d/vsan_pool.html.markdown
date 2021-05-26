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

  name       = "example"
  alloc_mode  = "static"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object VSAN Pool.
* `alloc_mode` - (Required) Allocation Mode of Object VSAN Pool. Allowed values: "dynamic", "static". 


## Attribute Reference

* `id` - Attribute id set to the Dn of the VSAN Pool.
* `description` - (Optional) Description for object VSAN Pool.
* `annotation` - (Optional) Annotation for object VSAN Pool.
* `name_alias` - (Optional) Name alias for object VSAN Pool.
