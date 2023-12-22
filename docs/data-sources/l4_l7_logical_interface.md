---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_logical_interface"
sidebar_current: "docs-aci-data-source-aci_l4_l7_logical_interface"
description: |-
  Data source for ACI L4-L7 Logical Interface
---

# aci_l4_l7_logical_interface #

Data source for ACI L4-L7 Logical Interface

## API Information ##

* `Class` - vnsLIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/lIf-{logical_interface_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Cluster Interfaces

## Example Usage ##

```hcl
data "aci_l4_l7_logical_interface" "example" {
  l4_l7_devices_dn  = aci_l4_l7_device.example.id
  name              = "example"
}
```

## Argument Reference ##

* `l4-l7_devices_dn` - (Required) Distinguished name of the parent L4-L7 Device object.
* `name` - (Required) Name of the Logical Interface object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Logical Interface.
* `annotation` - (Optional) Annotation of the Logical Interface object.
* `name_alias` - (Optional) Name Alias of the Logical Interface object.
* `encap` - (Optional) The port encapsulation to be used with the device. Type: String.
* `lag_policy_name` - (Optional) Name of the enhanced Lag policy. Type: String.
* `relation_vns_rs_c_if_att_n` - (Optional) Represents the relation between a set of Concrete Interfaces and the Device Cluster (class vnsCIf). Type: List.