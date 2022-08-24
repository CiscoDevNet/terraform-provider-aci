---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_multicast_pool_block"
sidebar_current: "docs-aci-resource-multicast_pool_block"
description: |-
  Manages ACI Abstraction of IP Address Block
---

# aci_multicast_pool_block #

Manages ACI the Multicast Address Pool Block

## API Information ##

* `Class` - fvnsMcastAddrBlk
* `Distinguished Name` - uni/infra/maddrns-{name}/fromaddr-[{from}]-toaddr-[{to}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> Multicast Address -> Address Blocks

## Example Usage ##

Do not use the `multicast_address_block` from the `aci_multicast_pool` resource in combination with this resource.

```hcl
resource "aci_multicast_pool" "example-pool" {
  name  = "example-pool"
}

resource "aci_multicast_pool_block" "example" {
  multicast_pool_dn  = aci_multicast_pool.example-pool.id
  from  = "224.0.0.10"
  to  = "224.0.0.20"
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

## Importing ##

An existing MulticastAddressBlock can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_multicast_pool_block.example <Dn>
```