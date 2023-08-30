---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_cloud_l4_l7_third_party_device"
sidebar_current: "docs-aci-data-source-cloud_l4_l7_third_party_device"
description: |-
  Data source for ACI Cloud L4-L7 Third Party Device
---

# aci_cloud_l4_l7_third_party_device #

Data source for ACI Cloud L4-L7 Third Party Device

Note: This data source is supported in Cloud Network Controller only.

## API Information ##

* `Class` - cloudLDev
* `Distinguished Name` - uni/tn-{tenant_name}/cld-{cld_name}

## GUI Information ##

* `Location` - Application Management -> Services -> Devices

## Example Usage ##

```hcl
data "aci_cloud_l4_l7_third_party_device" "example" {
  tenant_dn        = aci_tenant.tf_tenant.id
  name             = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the Cloud L4-L7 Third Party Device object. Type: String.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud L4-L7 Third Party Device. Type: String.
* `annotation` - (Read-Only) Annotation of the Cloud L4-L7 Third Party Device object. Type: String.
* `name_alias` - (Read-Only) Name Alias of the Cloud L4-L7 Third Party Device object. Type: String.
* `version` - (Read-Only) Version of the Cloud L4-L7 Third Party Device object. Type: String.
* `active_active` - (Read-Only) Active-Active mode of the Cloud L4-L7 Third Party Device object. Type: String.
* `context_aware` - (Read-Only) A value to determine if the L4-L7 device cluster supports multiple contexts (VRFs). Type: String.
* `custom_rg` - (Read-Only) Custom RG of the Cloud L4-L7 Third Party Device object. Type: String.
* `device_type` - (Read-Only) Device Type of the Cloud L4-L7 Third Party Device object. Type: String.
* `function_type` - (Read-Only) Function Type of the Cloud L4-L7 Third Party Device object. Type: String.
* `instance_count` - (Read-Only) Instance Count of the Cloud L4-L7 Third Party Device object. Type: String.
* `is_copy` - (Read-Only) Is the device is a copy device. Type: String.
* `is_instantiation` - (Read-Only) Is Instantiation of the Cloud L4-L7 Third Party Device object. Type: String.
* `l4l7_device_application_security_group` - (Read-Only) Naming for the Third Party Device Application Security Group of the Cloud L4-L7 Third Party Device object. Type: String.
* `l4l7_third_party_device` - (Read-Only) Naming for the Third Party Device of the Cloud L4-L7 Third Party Device object. Type: String.
* `managed` - (Read-Only) Is the device is managed. Type: String.
* `mode` - (Read-Only) Mode of the Cloud L4-L7 Third Party Device object. The value for specifying if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS). Type: String.
* `package_model` - (Read-Only) Package Model of the Cloud L4-L7 Third Party Device object. Type: String.
* `prom_mode` - (Read-Only) Promiscuous Mode of the Cloud L4-L7 Third Party Device object. Type: String.
* `service_type` - (Read-Only) Service Type of the Cloud L4-L7 Third Party Device object. Type: String.
* `target_mode` - (Read-Only) Target Mode of the Cloud L4-L7 Third Party Device object. Type: String.
* `trunking` - (Read-Only) For virtual devices, if a trunking port group is to be used. Type: String.
* `interface_selectors` - (Read-Only) Interface Selectors of the Cloud L4-L7 Third Party Device object. Type: Block.
  * `name` - (Read-Only) Name of the Interface Selector object. Type: String.
  * `allow_all` - (Read-Only) Allow-All of the Interface Selector object. Type: String.
  * `end_point_selectors` - (Read-Only) End Point Selectors of the Interface Selector object. Type: Block.
    * `name` - (Read-Only) Name of the End Point Selectors object. Type: String.
    * `match_expression` - (Read-Only) Match Expression of the End Point Selectors object. Type: String.
* `aaa_domain_dn` - (Read-Only) Represents the relation to a Relation from AAA Domain to Cloud L4L7 Native Load Balancer (class aaaRbacAnnotation). Type: List.
* `relation_cloud_rs_ldev_to_ctx` - (Read-Only) Represents the relation to a Relation from Cloud LDev to Cloud Context (class fvCtx). Type: String.
