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
* `name` - (Required) name of Object vxlan_pool.



## Attribute Reference

* `id` - Attribute id set to the Dn of the VXLAN Pool.
* `annotation` - (Optional) annotation for object vxlan_pool.
* `name_alias` - (Optional) name_alias for object vxlan_pool.
