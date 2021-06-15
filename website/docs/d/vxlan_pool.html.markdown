---
layout: "aci"
page_title: "ACI: aci_vxlan_pool"
sidebar_current: "docs-aci-data-source-vxlan_pool"
description: |-
  Data source for ACI VXLAN Pool
---

# aci_vxlan_pool #
Data source for ACI VXLAN Pool

## Example Usage ##

```hcl
data "aci_vxlan_pool" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object vxlan pool.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VXLAN Pool.
* `annotation` - (Optional) Annotation for object vxlan pool.
* `description` - (Optional) Description for object vxlan pool.
* `name_alias` - (Optional) Name alias for object vxlan pool.
