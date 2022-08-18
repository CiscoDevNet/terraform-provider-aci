---
subcategory: "Access Policies"
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

* `Location` - Fabric -> Access Policies -> Pools -> Multicast Address

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
* `multicast_address_block` - (Optional) Multicast Address Pool Block Configuration. 
* `multicast_address_block.from` - (Required) First multicast ip of the Multicast Address Pool Block.
* `multicast_address_block.to` - (Required) Last multicast ip of the Multicast Address Pool Block.
* `multicast_address_block.name` - (Optional) Name Alias of the Multicast Address Pool Block. 
* `multicast_address_block.annotation` - (Optional) Annotation of Multicast Address Pool Block.
* `multicast_address_block.description` - (Optional) Description of Multicast Address Pool Block.
* `multicast_address_block.name_alias` - (Optional) Name Alias of Multicast Address Pool Block.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Multicast Address Pool.
* `multicast_address_block.dn` - (Optional) Multicast Address Pool Block Dn.