---
layout: "aci"
page_title: "ACI: aci_concrete_device"
sidebar_current: "docs-aci-resource-concrete_device"
description: |-
  Manages ACI Concrete Device
---

# aci_concrete_device #

Manages ACI Concrete Device

## API Information ##

* `Class` - vnsCDev
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/cDev-{concrete_device_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Concrete Device

## Example Usage ##

```hcl
resource "aci_concrete_device" "example" {
  l4_l7_devices_dn                 = aci_l4_l7_devices.example.id
  name                             = "tenant1-ASA1"
  vcenter_name                     = "vcenter"
  vm_name                          = "tenant1-ASA1"
  relation_vns_rs_c_dev_to_ctrlr_p = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
}
```

## Argument Reference ##

* `l4-l7_devices_dn` - (Required) Distinguished name of the parent L4-L7 Devices object.
* `name` - (Required) Name of the object Concrete Device.
* `annotation` - (Optional) Annotation of the object Concrete Device.
* `name_alias` - (Optional) Name Alias of object Concrete Device.
* `vcenter_name` - (Optional) The virtual center name on which the device is hosted in the L4-L7 device cluster. It uniquely identifies the virtual center. Type: String.
* `vm_name` - (Optional) The virtual center VM name on which the device is hosted in the L4-L7 device cluster. It uniquely identifies the VM. Type: String.
* `relation_vns_rs_c_dev_to_ctrlr_p` - (Optional) Represents the relation from a Concrete Device to a VMM Controller (class vmmCtrlrP). It is an implicit relation to validate the controller profile. Type: String.

## Importing ##

An existing Concrete Device can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_concrete_device.example <Dn>
```