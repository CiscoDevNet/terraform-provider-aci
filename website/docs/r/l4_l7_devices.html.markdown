---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_device"
sidebar_current: "docs-aci-resource-l4_l7_device"
description: |-
  Manages ACI L4-L7 Device
---

# aci_l4_l7_device #

Manages ACI L4-L7 Device

## API Information ##

* `Class` - vnsLDevVip
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}

## GUI Information ##

* `Location` - Tenant -> Services -> Devices


## Example Usage ##

### Creating a PHYSICAL Device ###

```hcl
resource "aci_l4_l7_device" "example1" {
  tenant_dn                            = aci_tenant.terraform_tenant.id
  name                                 = "example1"
  active                               = "no"
  context_aware                        = "single-Context"
  device_type                          = "PHYSICAL"
  function_type                        = "GoTo"
  is_copy                              = "no"
  mode                                 = "legacy-Mode"
  promiscuous_mode                     = "no"
  service_type                         = "OTHERS"
  relation_vns_rs_al_dev_to_phys_dom_p = "uni/phys-test_dom"
}
```

### Creating a VIRTUAL Device ###

```hcl
resource "aci_l4_l7_device" "example2" {
  tenant_dn        = aci_tenant.terraform_tenant.id
  name             = "example2"
  active           = "no"
  context_aware    = "single-Context"
  device_type      = "VIRTUAL"
  function_type    = "GoTo"
  is_copy          = "no"
  mode             = "legacy-Mode"
  promiscuous_mode = "no"
  service_type     = "OTHERS"
  trunking         = "no"
  relation_vns_rs_al_dev_to_dom_p {
    target_dn      = "uni/vmmp-VMware/dom-ESX0-leaf102-vds"
    switching_mode = "AVE"
  }
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L4-L7 Device.
* `annotation` - (Optional) Annotation of the L4-L7 Device object.
* `active` - (Optional) Enables L4-L7 device cluster to operate in active/active mode. Allowed values are "no", "yes", and default value is "no". Type: String.
* `context_aware` - (Optional) Determines if the L4-L7 device cluster supports multiple contexts (VRFs). Allowed values are "multi-Context", "single-Context", and default value is "single-Context". Type: String.
* `device_type` - (Optional) Device Type. Allowed values are "CLOUD", "PHYSICAL", "VIRTUAL", and default value is "PHYSICAL". Type: String.
* `function_type` - (Optional) Function Type of the L4-L7 device cluster. Allowed values are "GoThrough", "GoTo", "L1", "L2", "None", and default value is "GoTo". Type: String.
* `is_copy` - (Optional) Sets device as a copy device. Allowed values are "no", "yes", and default value is "no". Type: String.
* `promiscuous_mode` - (Optional) Enabling Promiscuous Mode to allow all the traffic in a port group to reach a VM attached to a promiscuous port. Allowed values are "no", "yes", and default value is "no". Type: String.
* `service_type` - (Optional) The type of service the L4-L7 device performs. Allowed values are "ADC", "COPY", "FW", "NATIVELB", "OTHERS", and default value is "OTHERS". Type: String.
* `trunking` - (Optional) Configures the device port group for trunking of virtual devices. Allowed values are "no", "yes", and default value is "no". Type: String.
* `relation_vns_rs_al_dev_to_dom_p` - (Optional) Represents a relation from L4-L7 Device to a VMM Domain Profile (class vmmDomP). Type: Block.
  * `domain_dn` - (Required) Distinguished name of the VMM Domain in which the VM is deployed. Type: String.
  * `switching_mode` - (Optional) Port group switching mode. Allowed values are "native", "AVE", and default value is "native". The value "AVE" is not supported with non-AVE VMM Domain. Type: String.
* `relation_vns_rs_al_dev_to_phys_dom_p` - (Optional) Represents a relation from L4-L7 Device to a Physical Domain Profile (class physDomP). Type: String.

## Importing ##

An existing L4-L7 Device can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l4-l7_device.example <Dn>
```