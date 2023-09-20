---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_multicast_pool"
sidebar_current: "docs-aci-resource-multicast_pool"
description: |-
  Manages the ACI Multicast Address Pool
---

# aci_multicast_pool #

Manages the ACI Multicast Address Pool

## API Information ##

* `Class` - fvnsMcastAddrInstP
* `Distinguished Name` - uni/infra/maddrns-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> Multicast Address

## Example Usage ##

Do not use the `multicast_address_block` from this resource in combination with the `aci_multicast_pool_block` resource.

```hcl
resource "aci_multicast_pool" "example-pool" {
  name  = "example-pool"
  multicast_address_block {
    from = "224.0.0.40"
    to = "224.0.0.44"
    name = "testing-1"
  }
}
```

## Argument Reference ##

* `name` - (Required) Name of the Multicast Address Pool.
* `annotation` - (Optional) Annotation of the Multicast Address Pool.
* `description` - (Optional) Description of the Multicast Address Pool.
* `name_alias` - (Optional) Name Alias of the Multicast Address Pool.
* `multicast_address_block` - (Optional) Multicast Address Pool Block Configuration. Type: block.
 * `from` - (Required) First multicast IP of the Multicast Address Pool Block.
 * `to` - (Required) Last multicast IP of the Multicast Address Pool Block.
 * `name` - (Optional) Name Alias of the Multicast Address Pool Block. 
 * `annotation` - (Optional) Annotation of the Multicast Address Pool Block.
 * `description` - (Optional) Description of the Multicast Address Pool Block.
 * `name_alias` - (Optional) Name Alias of the Multicast Address Pool Block.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Multicast Address Pool.

## Importing ##

An existing MulticastAddressPool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_multicast_pool.example <Dn>
```