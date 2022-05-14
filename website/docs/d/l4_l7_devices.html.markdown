---
layout: "aci"
page_title: "ACI: aci_l4_l7_devices"
sidebar_current: "docs-aci-data-source-l4_l7_devices"
description: |-
  Data source for ACI L4-L7 Devices
---

# aci_l4_l7_devices #

Data source for ACI L4-L7 Devices


## API Information ##

* `Class` - vnsLDevVip
* `Distinguished Name` - uni/tn-{tenant_name}/lDevVip-{device_name}

## GUI Information ##

* `Location` -  Tenant -> Services -> Devices



## Example Usage ##

```hcl
data "aci_l4_l7_devices" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of object L4-L7 Devices.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the L4-L7 Devices.
* `annotation` - (Optional) Annotation of object L4-L7 Devices.
* `name_alias` - (Optional) Name Alias of object L4-L7 Devices.
* `active` - (Optional) Active mode. Allowed values are "no", "yes", and default value is "no". Type: String.
* `context_aware` - (Optional) Determines if the L4-L7 device cluster supports multiple contexts (VRFs). Allowed values are "multi-Context", "single-Context", and default value is "single-Context". Type: String.
* `device_type` - (Optional) Device Type. Allowed values are "CLOUD", "PHYSICAL", "VIRTUAL", and default value is "PHYSICAL". Type: String.
* `function_type` - (Optional) Function Type of the L4-L7 device cluster. Allowed values are "GoThrough", "GoTo", "L1", "L2", "None", and default value is "GoTo". Type: String.
* `is_copy` - (Optional) Sets device as a copy. Allowed values are "no", "yes", and default value is "no". Type: String.
* `mode` - (Optional) Specifies if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS). Allowed value is "legacy-Mode", and default value is "legacy-Mode". Type: String.
* `promiscuous_mode` - (Optional) Promiscuous Mode supports port groups in an external VMM controller, such as a Vcenter. This needs to be turned on only for service devices in the cloud, not for Enterprise AVE service deployments. Allowed values are "no", "yes", and default value is "no". Type: String.
* `service_type` - (Optional) Service Type. Allowed values are "ADC", "COPY", "FW", "NATIVELB", "OTHERS", and default value is "OTHERS". Type: String.
* `trunking` - (Optional) Sets trunking port group for virtual devices. Allowed values are "no", "yes", and default value is "no". Type: String.
* `relation_vns_rs_al_dev_to_dom_p` - (Optional) Represents a relation from L4-L7 Device to a Vmm Domain Profile (class vmmDomP). Type: Block.
  * `domain_dn` - (Required) Distinguished name of the target. Type: String.
  * `switching_mode` - (Optional) Switching mode. Allowed values are "native", "AVE", and default value is "native". The value "AVE" is not supported with non-AVE VMM Domain. Type: String.
* `relation_vns_rs_al_dev_to_phys_dom_p` - (Optional) Represents a relation from L4-L7 Device to a Physical Domain Profile (class physDomP). Type: String.
