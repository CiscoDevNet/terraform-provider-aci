---
subcategory: "L4-L7 Services"
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

### Creating a Virtual Concrete Device ###

```hcl
resource "aci_concrete_device" "example1" {
  l4_l7_device_dn   = aci_l4_l7_device.example.id
  name              = "example1"
  vmm_controller_dn = "uni/vmmp-VMware/dom-ACI-vDS/ctrlr-vcenter"
  vm_name           = "tenant1-ASA1"
}
```
### Creating a Physical Concrete Device ###

```hcl
resource "aci_concrete_device" "example2" {
  l4_l7_device_dn   = aci_l4_l7_device.example.id
  name              = "example2"
}
```

## Argument Reference ##

* `l4-l7_device_dn` - (Required) Distinguished name of the parent L4-L7 Device object.
* `name` - (Required) Name of the Concrete Device object.
* `annotation` - (Optional) Annotation of the Concrete Device object.
* `name_alias` - (Optional) Name Alias of the Concrete Device object.
* `vmm_controller_dn` - (Optional) Distinguished name of the VMM controller object. This can only be used for Virtual L4-L7 Devices. Type: String.
* `vm_name` - (Optional) The name of the Virtual Machine (VM) in the vCenter on which the device in the L4-L7 device cluster is hosted. It uniquely identifies the VM. This can only be used for Virtual L4-L7 Devices. Type: String.

## Importing ##

An existing Concrete Device can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_concrete_device.example <Dn>
```