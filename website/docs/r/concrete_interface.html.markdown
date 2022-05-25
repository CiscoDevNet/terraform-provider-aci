---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_concrete_interface"
sidebar_current: "docs-aci-resource-concrete_interface"
description: |-
  Manages ACI Concrete Interface
---

# aci_concrete_interface #

Manages ACI Concrete Interface

## API Information ##

* `Class` - vnsCIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/cDev-{concrete_device_name}/cIf-[{concrete_interface_name}]

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Concrete Device -> Concrete Interface

## Example Usage ##

### Creating an interface for a Virtual Concrete Device ###

```hcl
resource "aci_concrete_interface" "example1" {
  concrete_device_dn            = aci_concrete_device.concrete.id
  name                          = "g0/4"
  encap                         = "unknown"
  vnic_name                     = "Network adapter 5"
  relation_vns_rs_c_if_path_att = "topology/pod-1/paths-101/pathep-[eth1/1]"
}
```

### Creating an interface for a Physical Concrete Device ###

```hcl
resource "aci_concrete_interface" "example2" {
  concrete_device_dn            = aci_concrete_device.concrete.id
  name                          = "g0/4"
}
```

## Argument Reference ##

* `concrete_device_dn` - (Required) Distinguished name of the parent Concrete Device object.
* `name` - (Required) Name of the Concrete Interface object.
* `encap` - (Optional) The port encapsulation. Type: String.
* `vnic_name` - (Optional) The virtual NIC (vNIC) name of the L4-L7 Device VM represented by the concrete interface. Type: String.
* `relation_vns_rs_c_if_path_att` - (Optional) Represents a relation from the Concrete Interface to the Physical Port on the Leaf (class fabricPathEp). Type: String.

## Importing ##

An existing Concrete Interface can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_concrete_interface.example <Dn>
```