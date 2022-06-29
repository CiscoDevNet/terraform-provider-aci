---
layout: "aci"
page_title: "ACI: aci_abstractionof_ipaddress_block"
sidebar_current: "docs-aci-data-source-abstractionof_ipaddress_block"
description: |-
  Data source for ACI Abstraction of IP Address Block
---

# aci_abstractionof_ipaddress_block #

Data source for ACI Abstraction of IP Address Block


## API Information ##

* `Class` - fvnsMcastAddrBlk
* `Distinguished Name` - uni/infra/maddrns-{name}/fromaddr-[{from}]-toaddr-[{to}]

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_abstractionof_ipaddress_block" "example" {
  multicast_address_pool_dn  = aci_multicast_address_pool.example.id
  _from  = "example"
  to  = "example"
}
```

## Argument Reference ##

* `multicast_address_pool_dn` - (Required) Distinguished name of parent MulticastAddressPool object.
* `_from` - (Required) _from of object Abstraction of IP Address Block.
* `to` - (Required) To of object Abstraction of IP Address Block.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Abstraction of IP Address Block.
* `annotation` - (Optional) Annotation of object Abstraction of IP Address Block.
* `name_alias` - (Optional) Name Alias of object Abstraction of IP Address Block.
* `from` - (Optional) Starting IP Address. Start of the multicast address block.
* `name_alias` - (Optional) Name alias. 
* `to` - (Optional) Ending IP Address. End of the multicast address block.
