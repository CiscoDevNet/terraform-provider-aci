---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_concrete_device"
sidebar_current: "docs-aci-data-source-aci_concrete_device"
description: |-
  Data source for ACI Concrete Device
---

# aci_concrete_device #

Data source for ACI Concrete Device

## API Information ##

* `Class` - vnsCDev
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}/cDev-{concrete_device_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices -> Concrete Device

## Example Usage ##

```hcl
data "aci_concrete_device" "example" {
  l4_l7_device_dn = aci_l4_l7_device.example.id
  name            = "example"
}
```

## Argument Reference ##

* `l4-l7_device_dn` - (Required) Distinguished name of the parent L4-L7 Device object.
* `name` - (Required) Name of the Concrete Device object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Concrete Device.
* `annotation` - (Optional) Annotation of the Concrete Device object.
* `name_alias` - (Optional) Name Alias of the Concrete Device object.
* `vmm_controller_dn` - (Optional) Distinguished name of the VMM controller object. Type: String.
* `vm_name` - (Optional) The name of the Virtual Machine (VM) in the vCenter on which the device is hosted in the L4-L7 device cluster. It uniquely identifies the VM. Type: String.
