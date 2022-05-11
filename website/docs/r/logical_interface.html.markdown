---
layout: "aci"
page_title: "ACI: aci_logical_interface"
sidebar_current: "docs-aci-resource-logical_interface"
description: |-
  Manages ACI Logical Interface
---

# aci_logical_interface #

Manages ACI Logical Interface

## API Information ##

* `Class` - vnsLIf
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/lIf-{cluster_interface_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Cluster Interfaces


## Example Usage ##

```hcl
resource "aci_logical_interface" "example" {
  l4_l7_devices_dn  = aci_l4_l7_devices.example.id
  name  = "example"
  encap = "unknown"
  vns_rs_c_if_att_n = ["uni/tn-tenant1/lDevVip-ok/cDev-test/cIf-[g0/1]"]
}
```

## Argument Reference ##

* `l4_l7_devices_dn` - (Required) Distinguished name of the parent L4-L7Devices object.
* `name` - (Required) Name of the object Logical Interface.
* `annotation` - (Optional) Annotation of the object Logical Interface.
* `encap` - (Optional) The port encapsulation. Type: String.
* `lag_policy_name` - (Optional) Enhanced LAG Policy Name. Type: String.
* `relation_vns_rs_c_if_att_n` - (Optional) Represents a relation for Set of Concrete Interfaces from the Device in the Cluster (class vnsCIf). Type: List.

## Importing ##

An existing LogicalInterface can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_logical_interface.example <Dn>
```