---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_logical_interface"
sidebar_current: "docs-aci-resource-l4_l7_logical_interface"
description: |-
  Manages ACI L4-L7 Logical Interface
---

# aci_l4_l7_logical_interface #

Manages ACI L4-L7 Logical Interface

## API Information ##

* `Class` - vnsLIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/lIf-{logical_interface_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Cluster Interfaces

## Example Usage ##

### Creating a logical interface for a virtual device ###

```hcl
resource "aci_l4_l7_logical_interface" "example1" {
  l4_l7_device_dn            = aci_l4_l7_device.virtual_device.id
  name                       = "example1"
  enhanced_lag_policy_name   = "Lacp"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.virtual_interface1.id, aci_concrete_interface.virtual_interface2.id]
}
```

### Creating a logical interface for a physical device ###

```hcl
resource "aci_l4_l7_logical_interface" "example2" {
  l4_l7_device_dn            = aci_l4_l7_device.physical_device.id
  name                       = "example2"
  encap                      = "vlan-1"
  relation_vns_rs_c_if_att_n = [aci_concrete_interface.physical_interface.id]
}
```

## Argument Reference ##

* `l4_l7_device_dn` - (Required) Distinguished name of the parent L4-L7 Device object.
* `name` - (Required) Name of the Logical Interface object.
* `annotation` - (Optional) Annotation of the Logical Interface object.
* `encap` - (Optional) The port encapsulation to be used with the device. Type: String.
* `enhanced_lag_policy_name` - (Optional) Name of the enhanced Lag policy. Type: String.
* `relation_vns_rs_c_if_att_n` - (Optional) Represents the relation between a set of Concrete Interfaces and the Device Cluster (class vnsCIf). Type: List.

## Importing ##

An existing L4-L7 Logical Interface can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l4_l7_logical_interface.example <Dn>
```