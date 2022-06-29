---
layout: "aci"
page_title: "ACI: aci_multicast_pool"
sidebar_current: "docs-aci-data-source-multicast_pool"
description: |-
  Data source for ACI Multicast Address Pool
---

# aci_multicast_pool #

Data source for ACI Multicast Address Pool

## API Information ##

* `Class` - fvnsMcastAddrInstP
* `Distinguished Name` - uni/infra/maddrns-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
data "aci_multicast_pool" "example-pool" {
  name  = "example-pool"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Multicast Address Pool.
* `annotation` - (Optional) Annotation of the Multicast Address Pool.
* `description` - (Optional) Description of the Multicast Address Pool.
* `name_alias` - (Optional) Name Alias of the Multicast Address Pool.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Multicast Address Pool.
