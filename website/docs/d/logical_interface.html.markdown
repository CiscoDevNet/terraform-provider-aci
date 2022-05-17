---
layout: "aci"
page_title: "ACI: aci_logical_interface"
sidebar_current: "docs-aci-data-source-logical_interface"
description: |-
  Data source for ACI Logical Interface
---

# aci_logical_interface #

Data source for ACI Logical Interface

## API Information ##

* `Class` - vnsLIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/lIf-{logical_interface_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Cluster Interfaces

## Example Usage ##

```hcl
data "aci_logical_interface" "example" {
  l4_l7_devices_dn  = aci_l4_l7_devices.example.id
  name              = "example"
}
```

## Argument Reference ##

* `l4-l7_devices_dn` - (Required) Distinguished name of parent L4-L7Devices object.
* `name` - (Required) Name of object Logical Interface.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Logical Interface.
* `annotation` - (Optional) Annotation of object Logical Interface.
* `name_alias` - (Optional) Name Alias of object Logical Interface.
* `encap` - (Optional) The port encapsulation.
* `lag_policy_name` - (Optional) Enhanced LAG Policy Name.
* `relation_vns_rs_c_if_att_n` - (Optional) Represents the relation between Set of Concrete Interfaces and the Device Cluster (class vnsCIf). Type: List.
