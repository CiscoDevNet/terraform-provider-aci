---
layout: "aci"
page_title: "ACI: aci_multicast_address_pool"
sidebar_current: "docs-aci-data-source-multicast_address_pool"
description: |-
  Data source for ACI Multicast Address Pool
---

# aci_multicast_address_pool #

Data source for ACI Multicast Address Pool


## API Information ##

* `Class` - fvnsMcastAddrInstP
* `Distinguished Name` - uni/infra/maddrns-{name}

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_multicast_address_pool" "example" {

  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object Multicast Address Pool.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Multicast Address Pool.
* `annotation` - (Optional) Annotation of object Multicast Address Pool.
* `name_alias` - (Optional) Name Alias of object Multicast Address Pool.
* `name_alias` - (Optional) Name alias. 
