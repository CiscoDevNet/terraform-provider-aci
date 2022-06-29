---
layout: "aci"
page_title: "ACI: aci_multicast_address_pool"
sidebar_current: "docs-aci-resource-multicast_address_pool"
description: |-
  Manages ACI Multicast Address Pool
---

# aci_multicast_address_pool #

Manages ACI Multicast Address Pool

## API Information ##

* `Class` - fvnsMcastAddrInstP
* `Distinguished Name` - uni/infra/maddrns-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_multicast_address_pool" "example" {

  name  = "example"
  annotation = "orchestrator:terraform"

  name_alias = 
}
```

## Argument Reference ##


* `name` - (Required) Name of the object Multicast Address Pool.
* `annotation` - (Optional) Annotation of the object Multicast Address Pool.

* `name_alias` - (Optional) Name alias.


## Importing ##

An existing MulticastAddressPool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_multicast_address_pool.example <Dn>
```