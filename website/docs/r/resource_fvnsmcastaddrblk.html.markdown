---
layout: "aci"
page_title: "ACI: aci_abstractionof_ipaddress_block"
sidebar_current: "docs-aci-resource-abstractionof_ipaddress_block"
description: |-
  Manages ACI Abstraction of IP Address Block
---

# aci_abstractionof_ipaddress_block #

Manages ACI Abstraction of IP Address Block

## API Information ##

* `Class` - fvnsMcastAddrBlk
* `Distinguished Name` - uni/infra/maddrns-{name}/fromaddr-[{from}]-toaddr-[{to}]

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_abstractionof_ipaddress_block" "example" {
  multicast_address_pool_dn  = aci_multicast_address_pool.example.id
  _from  = "example"
  to  = "example"
  annotation = "orchestrator:terraform"
  from = 

  name_alias = 
  to = 
}
```

## Argument Reference ##

* `multicast_address_pool_dn` - (Required) Distinguished name of the parent MulticastAddressPool object.
* `_from` - (Required) _from of the object Abstraction of IP Address Block.* `to` - (Required) To of the object Abstraction of IP Address Block.
* `annotation` - (Optional) Annotation of the object Abstraction of IP Address Block.

* `from` - (Optional) Starting IP Address.Start of the multicast address block.
* `name_alias` - (Optional) Name alias.
* `to` - (Optional) Ending IP Address.End of the multicast address block.


## Importing ##

An existing AbstractionofIPAddressBlock can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_abstractionof_ipaddress_block.example <Dn>
```