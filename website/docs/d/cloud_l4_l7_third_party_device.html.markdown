---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_cloud_l4_l7_third_party_device"
sidebar_current: "docs-aci-data-source-cloud_l4_l7_third_party_device"
description: |-
  Data source for ACI Cloud L4-L7 Device
---

# aci_cloud_l4_l7_third_party_device #

Data source for ACI Cloud L4-L7 Device

Note: This resource is supported in Cloud Network Controller only.

## API Information ##

* `Class` - cloudLDev
* `Distinguished Name` - uni/tn-{tenant_name}/cld-{cld_name}

## GUI Information ##

* `Location` - Application Management -> Services -> Devices

## Example Usage ##

```hcl
data "aci_cloud_l4_l7_third_party_device" "cloud_third_party_fw" {
  tenant_dn        = aci_tenant.tf_tenant.id
  name             = "cloud_third_party_fw"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Cloud L4-L7 Device object.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud L4-L7 Device.
* `annotation` - (Read-Only) Annotation of the Cloud L4-L7 Device object.
* `name_alias` - (Read-Only) Name Alias of the Cloud L4-L7 Device object.
* `version` - (Read-Only) Version of the Cloud L4-L7 Device object.
* `active_active` - (Read-Only) Active-Active mode of the Cloud L4-L7 Device object.
* `context_aware` - (Read-Only) A value to determine if the L4-L7 device cluster supports multiple contexts (VRFs).
* `custom_rg` - (Read-Only) Custom RG of the Cloud L4-L7 Device object.
* `device_type` - (Read-Only) Device Type of the Cloud L4-L7 Device object.
* `function_type` - (Read-Only) Function Type of the Cloud L4-L7 Device object.
* `instance_count` - (Read-Only) Instance Count of the Cloud L4-L7 Device object.
* `is_copy` - (Read-Only) Is the device is a copy device.
* `is_instantiation` - (Read-Only) Is Instantiation of the Cloud L4-L7 Device object.
* `l4l7_device_application_security_group` - (Read-Only) Naming for the Third Party Device Application Security Group of the Cloud L4-L7 Device object.
* `l4l7_third_party_device` - (Read-Only) Naming for the Third Party Device of the Cloud L4-L7 Device object.
* `managed` - (Read-Only) Is the device is managed.
* `mode` - (Read-Only) Mode of the Cloud L4-L7 Device object. The value for specifying if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS).
* `package_model` - (Read-Only) Package Model of the Cloud L4-L7 Device object.
* `prom_mode` - (Read-Only) Promiscuous Mode of the Cloud L4-L7 Device object.
* `service_type` - (Read-Only) Service Type of the Cloud L4-L7 Device object.
* `target_mode` - (Read-Only) Target Mode of the Cloud L4-L7 Device object.
* `trunking` - (Read-Only) For virtual devices, if a trunking port group is to be used.
* `interface_selectors` - (Read-Only) Interface Selectors of the Cloud L4-L7 Device object. Type: Block.
  * `name` - (Read-Only) Name of the Interface Selector object.
  * `allow_all` - (Read-Only) Allow-All of the Interface Selector object.
  * `end_point_selectors` - (Read-Only) End Point Selectors of the Interface Selector object. Type: Block.
    * `name` - (Read-Only) Name of the End Point Selectors object.
    * `match_expression` - (Read-Only) Match Expression of the End Point Selectors object.
* `aaa_domain_dn` - (Read-Only) Represents the relation to a Relation from AAA Domain to Cloud L4L7 Native Load Balancer (class aaaRbacAnnotation). Type: List.
* `relation_cloud_rs_ldev_to_ctx` - (Read-Only) Represents the relation to a Relation from Cloud LDev to Cloud Context (class fvCtx). Type: String.
