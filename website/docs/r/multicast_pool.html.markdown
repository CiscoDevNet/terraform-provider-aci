---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_multicast_pool"
sidebar_current: "docs-aci-resource-multicast_pool"
description: |-
  Manages ACI Multicast Address Pool
---

# aci_multicast_pool #

Manages ACI Multicast Address Pool

## API Information ##

* `Class` - fvnsMcastAddrInstP
* `Distinguished Name` - uni/infra/maddrns-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Pools -> Multicast Address

## Example Usage ##

```hcl
resource "aci_multicast_pool" "example-pool" {
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

## Importing ##

An existing MulticastAddressPool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_multicast_pool.example <Dn>
```