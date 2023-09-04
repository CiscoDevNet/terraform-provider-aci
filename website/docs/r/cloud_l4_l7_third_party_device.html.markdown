---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_cloud_l4_l7_third_party_device"
sidebar_current: "docs-aci-resource-cloud_l4_l7_third_party_device"
description: |-
  Manages ACI Cloud L4-L7 Third Party Device
---

# aci_cloud_l4_l7_third_party_device #

Manages ACI Cloud L4-L7 Third Party Device

Note: This resource is supported in Azure Cloud Network Controller only.

## API Information ##

* `Class` - cloudLDev
* `Distinguished Name` - uni/tn-{tenant_name}/cld-{cld_name}

## GUI Information ##

* `Location` - Application Management -> Services -> Devices


## Example Usage ##

```hcl
resource "aci_cloud_l4_l7_third_party_device" "example" {
  tenant_dn        = aci_tenant.tf_tenant.id
  name             = "example"
  active_active    = "no"
  context_aware    = "single-Context"
  device_type      = "CLOUD"
  function_type    = "GoTo"
  instance_count   = "2"
  is_copy          = "no"
  is_instantiation = "no"
  managed          = "yes"
  prom_mode        = "no"
  service_type     = "FW"
  target_mode      = "unspecified"
  trunking         = "no"

  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id,
    aci_aaa_domain.aaa_domain_2.id
  ]

  relation_cloud_rs_ldev_to_ctx = aci_vrf.vrf1.id

  interface_selectors {
    allow_all = "no"
    name      = "Interface_1"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_1_ep_1"
    }
  }
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the Cloud L4-L7 Third Party Device object. Type: String.
* `annotation` - (Optional) Annotation of the Cloud L4-L7 Third Party Device object. Type: String.
* `name_alias` - (Optional) Name Alias of the Cloud L4-L7 Third Party Device object. Type: String.
* `version` - (Optional) Version of the Cloud L4-L7 Third Party Device object. Allowed values are and default value is "1.0". Type: String.
* `active_active` - (Optional) Active-Active mode of the Cloud L4-L7 Third Party Device object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `context_aware` - (Optional) A value to determine if the L4-L7 device cluster supports multiple contexts (VRFs). Allowed values are "multi-Context", "single-Context", and default value is "single-Context". Type: String.
* `custom_rg` - (Optional) Custom RG of the Cloud L4-L7 Third Party Device object. Type: String.
* `device_type` - (Optional) Device Type of the Cloud L4-L7 Third Party Device object. Allowed values are "CLOUD", "PHYSICAL", "VIRTUAL", and default value is "CLOUD". Type: String.
* `function_type` - (Optional) Function Type of the Cloud L4-L7 Third Party Device object. Allowed values are "GoThrough", "GoTo", "L1", "L2", "None", and default value is "GoTo". Type: String.
* `instance_count` - (Optional) Instance Count of the Cloud L4-L7 Third Party Device object. Type: String.
* `is_copy` - (Optional) Is the device is a copy device. Allowed values are "no", "yes", and default value is "no". Type: String.
* `is_instantiation` - (Optional) Is Instantiation of the Cloud L4-L7 Third Party Device object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `l4l7_device_application_security_group` - (Optional) Naming for the Third Party Device Application Security Group of the Cloud L4-L7 Third Party Device object. Type: String.
* `l4l7_third_party_device` - (Optional) Naming for the Third Party Device of the Cloud L4-L7 Third Party Device object. Type: String.
* `managed` - (Optional) Is the device is managed. Allowed values are "no", "yes", and default value is "yes". Type: String.
* `mode` - (Read-Only) Mode of the Cloud L4-L7 Third Party Device object. The value for specifying if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS). Type: String.
* `package_model` - (Optional) Package Model of the Cloud L4-L7 Third Party Device object. Type: String.
* `prom_mode` - (Optional) Promiscuous Mode of the Cloud L4-L7 Third Party Device object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `service_type` - (Optional) Service Type of the Cloud L4-L7 Third Party Device object. Allowed values are "ADC", "COPY", "FW", "NATIVELB", "OTHERS", and default value is "FW". Type: String.
* `target_mode` - (Optional) Target Mode of the Cloud L4-L7 Third Party Device object. Allowed values are "primary", "secondary", "unspecified", and default value is "unspecified". Type: String.
* `trunking` - (Optional) For virtual devices, if a trunking port group is to be used. Allowed values are "no", "yes", and default value is "no". Type: String.
* `interface_selectors` - (Optional) Interface Selectors of the Cloud L4-L7 Third Party Device object. Type: Block.
  * `name` - (Optional) Name of the Interface Selector object. Type: String.
  * `allow_all` - (Optional) Allow-All of the Interface Selector object. Type: String.
  * `end_point_selectors` - (Optional) End Point Selectors of the Interface Selector object. Type: Block.
    * `name` - (Optional) Name of the End Point Selectors object. Type: String.
    * `match_expression` - (Optional) Match Expression of the End Point Selectors object. Type: String.
* `aaa_domain_dn` - (Optional) Represents the relation to a Relation from AAA Domain to Cloud L4L7 Native Load Balancer (class aaaRbacAnnotation). Type: List.
* `relation_cloud_rs_ldev_to_ctx` - (Optional) Represents the relation to a Relation from Cloud LDev to Cloud Context (class fvCtx). Type: String.

## Importing ##

An existing Cloud L4-L7 Third Party Device can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_l4_l7_third_party_device.example <Dn>
```

Starting in Terraform version 1.5, an existing Cloud L4-L7 Third Party Device can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

 ```
 import {
    id = "<Dn>"
    to = aci_cloud_l4_l7_third_party_device.example
 }
 ```