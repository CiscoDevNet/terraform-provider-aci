---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_device"
sidebar_current: "docs-aci-data-source-l4_l7_device"
description: |-
  Data source for ACI L4-L7 Device
---

# aci_l4_l7_device #

Data source for ACI L4-L7 Device


## API Information ##

* `Class` - vnsLDevVip
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}

## GUI Information ##

* `Location` -  Tenant -> Services -> Devices



## Example Usage ##

```hcl
data "aci_l4_l7_device" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L4-L7 Device.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the L4-L7 Device.
* `annotation` - (Optional) Annotation of the L4-L7 Device object.
* `name_alias` - (Optional) Name Alias of the L4-L7 Device object.
* `active` - (Optional) Enables L4-L7 device cluster to operate in active/active mode. Allowed values are "no", "yes", and default value is "no". Type: String.
* `context_aware` - (Optional) Determines if the L4-L7 device cluster supports multiple contexts (VRFs). Allowed values are "multi-Context", "single-Context", and default value is "single-Context". Type: String.
* `device_type` - (Optional) Device Type. Allowed values are "CLOUD", "PHYSICAL", "VIRTUAL", and default value is "PHYSICAL". Type: String.
* `function_type` - (Optional) Function Type of the L4-L7 device cluster. Allowed values are "GoThrough", "GoTo", "L1", "L2", "None", and default value is "GoTo". Type: String.
* `is_copy` - (Optional) Sets device as a copy. Allowed values are "no", "yes", and default value is "no". Type: String.
* `promiscuous_mode` - (Optional) Enabling Promiscuous Mode supports all the traffic in a port group to reach a VM attached to a promiscuous port. Allowed values are "no", "yes", and default value is "no". Type: String.
* `service_type` - (Optional) The type of service the L4-L7 device performs. Allowed values are "ADC", "COPY", "FW", "NATIVELB", "OTHERS", and default value is "OTHERS". Type: String.
* `trunking` - (Optional) Configures the device port group for trunking of virtual devices. Allowed values are "no", "yes", and default value is "no". Type: String.
* `relation_vns_rs_al_dev_to_dom_p` - (Optional) Represents a relation from L4-L7 Device to a VMM Domain Profile (class vmmDomP). Type: Block.
  * `domain_dn` - (Required) Distinguished name of the VMM Domain in which the VM is deployed. Type: String.
  * `switching_mode` - (Optional) Port group switching mode. Allowed values are "native", "AVE", and default value is "native". The value "AVE" is not supported with non-AVE VMM Domain. Type: String.
* `relation_vns_rs_al_dev_to_phys_dom_p` - (Optional) Represents a relation from L4-L7 Device to a Physical Domain Profile (class physDomP). Type: String.
