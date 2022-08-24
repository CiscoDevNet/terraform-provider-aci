---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_multicast_pool_block"
sidebar_current: "docs-aci-data-source-multicast_pool_block"
description: |-
  Data source for ACI the Multicast Address Pool Block
---

# aci_multicast_pool_block #

Data source for ACI the Multicast Address Pool Block

## API Information ##

* `Class` - fvnsMcastAddrBlk
* `Distinguished Name` - uni/infra/maddrns-{name}/fromaddr-[{from}]-toaddr-[{to}]

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> Multicast Address -> Address Blocks

## Example Usage ##

```hcl
data "aci_multicast_pool" "example-pool" {
  name  = "example-pool"
}

data "aci_multicast_pool_block" "example" {
  multicast_pool_dn  = data.aci_multicast_pool.example-pool.id
  from  = "224.0.0.30"
  to  = "224.0.0.40"
}
```

## Argument Reference ##

* `multicast_pool_dn` - (Required) Distinguished name of the parent Multicast Pool object.
* `from` - (Required) First multicast ip of the Multicast Address Pool Block.
* `to` - (Required) Last multicast ip of the Multicast Address Pool Block.
* `annotation` - (Optional) Annotation of the Multicast Address Pool Block.
* `description` - (Optional) Description of the Multicast Address Pool Block.
* `name_alias` - (Optional) Name Alias of the Multicast Address Pool Block.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the of Multicast Address Pool Block.
