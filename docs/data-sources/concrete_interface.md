---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_concrete_interface"
sidebar_current: "docs-aci-data-source-concrete_interface"
description: |-
  Data source for ACI Concrete Interface
---

# aci_concrete_interface #

Data source for ACI Concrete Interface

## API Information ##

* `Class` - vnsCIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/cDev-{concrete_device_name}/cIf-[{concrete_interface_name}]

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Concrete Device -> Concrete Interface

## Example Usage ##

```hcl
data "aci_concrete_interface" "example" {
  concrete_device_dn  = aci_concrete_device.example.id
  name                = "example"
}
```

## Argument Reference ##

* `concrete_device_dn` - (Required) Distinguished name of the parent Concrete Device object.
* `name` - (Required) Name of the Concrete Interface object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Concrete Interface.
* `name_alias` - (Optional) Name Alias of the Concrete Interface object.
* `encap` - (Optional) The port encapsulation. Type: String.
* `vnic_name` - (Optional) The virtual NIC (vNIC) name of the L4-L7 Device VM represented by the concrete interface. Type: String.
* `relation_vns_rs_c_if_path_att` - (Optional) Represents a relation from the Concrete Interface to the Physical Port on the Leaf (class fabricPathEp). Type: String.
