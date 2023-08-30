---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_cloud_l4_l7_native_load_balancer"
sidebar_current: "docs-aci-data-source-cloud_l4_l7_native_load_balancer"
description: |-
  Data source for ACI Cloud L4-L7 Native Load Balancer
---

# aci_cloud_l4_l7_native_load_balancer #

Data source for ACI Cloud L4-L7 Native Load Balancer

Note: This data source is supported in Azure Cloud Network Controller only.

## API Information ##

* `Class` - cloudLB
* `Distinguished Name` - uni/tn-{tenant_name}/clb-{lb_name}

## GUI Information ##

* `Location` - Application Management -> Services -> Devices


## Example Usage ##

```hcl
data "aci_cloud_l4_l7_native_load_balancer" "example" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the Cloud L4-L7 Native Load Balancer object. Type: String.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `annotation` - (Read-Only) Annotation of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `name_alias` - (Read-Only) Name Alias of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `version` - (Read-Only) Version of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `active_active` - (Read-Only) Active-Active mode of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `allow_all` - (Read-Only) Allow-All of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `auto_scaling` - (Read-Only) Auto-Scaling of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `context_aware` - (Read-Only) A value to determine if the L4-L7 device cluster supports multiple contexts (VRFs). Type: String.
* `custom_rg` - (Read-Only) Custom RG of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `device_type` - (Read-Only) Device Type of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `function_type` - (Read-Only) Function Type of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `instance_count` - (Read-Only) Instance Count of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `is_copy` - (Read-Only) Is the device is a copy device. Type: String.
* `is_instantiation` - (Read-Only) Is Instantiation of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `is_static_ip` - (Read-Only) Is Static IP of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `l4l7_device_application_security_group` - (Read-Only) Naming for the Third Party Device Application Security Group of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `l4l7_third_party_device` - (Read-Only) Naming for the Third Party Device of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `managed` - (Read-Only) Is the device is managed. Type: String.
* `max_instance_count` - (Read-Only) Maximum Instance Count of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `min_instance_count` - (Read-Only) Minimum Instance Count of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `mode` - (Read-Only) Mode of the Cloud L4-L7 Native Load Balancer object. The value for specifying if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS). Type: String.
* `native_lb_name` - (Read-Only) Naming for the Native Load Balancer Devices of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `package_model` - (Read-Only) Package Model of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `prom_mode` - (Read-Only) Promiscuous Mode of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `scheme` - (Read-Only) Scheme of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `size` - (Read-Only) Size of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `sku` - (Read-Only) SKU of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `service_type` - (Read-Only) Service Type of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `target_mode` - (Read-Only) Target Mode of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `trunking` - (Read-Only) For virtual devices, if a trunking port group is to be used. Type: String.
* `cloud_l4l7_load_balancer_type` - (Read-Only) Type of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `relation_cloud_rs_ldev_to_cloud_subnet` - (Read-Only) Represents the relation to a Relation from Cloud LDev to Cloud Subnet (class cloudSubnet). Type: List.
* `aaa_domain_dn` - (Read-Only) Represents the relation to a Relation from AAA Domain to Cloud L4L7 Native Load Balancer (class aaaRbacAnnotation). Type: List.