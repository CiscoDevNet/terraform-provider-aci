---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_l4_l7_native_load_balancer"
sidebar_current: "docs-aci-resource-aci_cloud_l4_l7_native_load_balancer"
description: |-
  Manages ACI Cloud L4-L7 Native Load Balancer
---

# aci_cloud_l4_l7_native_load_balancer #

Manages ACI Cloud L4-L7 Native Load Balancer

Note: This resource is supported in Azure Cloud Network Controller only.

## API Information ##

* `Class` - cloudLB
* `Distinguished Name` - uni/tn-{tenant_name}/clb-{lb_name}

## GUI Information ##

* `Location` - Application Management -> Services -> Devices


## Example Usage ##

```hcl
resource "aci_cloud_l4_l7_native_load_balancer" "example" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "example"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id,
    aci_aaa_domain.aaa_domain_2.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
  active_active                 = "no"
  allow_all                     = "no"
  auto_scaling                  = "no"
  context_aware                 = "multi-Context"
  device_type                   = "CLOUD"
  function_type                 = "GoTo"
  is_copy                       = "no"
  is_instantiation              = "no"
  is_static_ip                  = "no"
  managed                       = "no"
  promiscuous_mode              = "no"
  scheme                        = "internal"
  size                          = "medium"
  sku                           = "standard"
  service_type                  = "NATIVELB"
  target_mode                   = "primary"
  trunking                      = "no"
  cloud_l4l7_load_balancer_type = "application"
  instance_count                = "2"
  max_instance_count            = "10"
  min_instance_count            = "5"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `annotation` - (Optional) Annotation of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `name_alias` - (Optional) Name Alias of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `version` - (Optional) Version of the Cloud L4-L7 Native Load Balancer object. Default value is "1.0". Type: String.
* `active_active` - (Optional) Active-Active mode of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `allow_all` - (Optional) Allow-All of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `auto_scaling` - (Optional) Auto-Scaling of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `context_aware` - (Optional) A value to determine if the L4-L7 device cluster supports multiple contexts (VRFs). Allowed values are "multi-Context", "single-Context", and default value is "single-Context". Type: String.
* `custom_resource_group` - (Optional) Custom Resource Group of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `device_type` - (Optional) Device Type of the Cloud L4-L7 Native Load Balancer object. Allowed values are "CLOUD", "PHYSICAL", "VIRTUAL", and default value is "CLOUD". Type: String.
* `function_type` - (Optional) Function Type of the Cloud L4-L7 Native Load Balancer object. Allowed values are "GoThrough", "GoTo", "L1", "L2", "None", and default value is "GoTo". Type: String.
* `instance_count` - (Optional) Instance Count of the Cloud L4-L7 Native Load Balancer object. Default value is "2". Type: String.
* `is_copy` - (Optional) Enables the device to be a copy device. Allowed values are "no", "yes", and default value is "no". Type: String.
* `is_instantiation` - (Optional) Enables Instantiation of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `is_static_ip` - (Optional) Enables static IP of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `l4l7_device_application_security_group` - (Optional) Naming for the Third Party Device Application Security Group of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `l4l7_third_party_device` - (Optional) Naming for the Third Party Device of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `managed` - (Optional) Enables the device to be a managed device. Allowed values are "no", "yes", and default value is "yes". Type: String.
* `max_instance_count` - (Optional) Maximum Instance Count of the Cloud L4-L7 Native Load Balancer object. Default value is "10". Type: String.
* `min_instance_count` - (Optional) Minimum Instance Count of the Cloud L4-L7 Native Load Balancer object. Default value is "0". Type: String.
* `mode` - (Read-Only) Mode of the Cloud L4-L7 Native Load Balancer object. The value for specifying if the device is legacy (classical VLAN/VXLAN) or supports service tag switching (STS). Type: String.
* `native_lb_name` - (Optional) Naming for the Native Load Balancer Devices of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `package_model` - (Optional) Package Model of the Cloud L4-L7 Native Load Balancer object. Type: String.
* `promiscuous_mode` - (Optional) Promiscuous Mode of the Cloud L4-L7 Native Load Balancer object. Allowed values are "no", "yes", and default value is "no". Type: String.
* `scheme` - (Optional) Scheme of the Cloud L4-L7 Native Load Balancer object. Allowed values are "internal", "internet", and default value is "internet". Type: String.
* `size` - (Optional) Size of the Cloud L4-L7 Native Load Balancer object. Allowed values are "large", "medium", "small", "v2", and default value is "medium". Type: String.
* `sku` - (Optional) SKU of the Cloud L4-L7 Native Load Balancer object. Allowed values are "WAF", "WAF_v2", "standard", "standard_v2", and default value is "standard". Type: String.
* `service_type` - (Optional) Service Type of the Cloud L4-L7 Native Load Balancer object. Allowed values are "ADC", "COPY", "FW", "NATIVELB", "OTHERS", and default value is "NATIVELB". Type: String.
* `target_mode` - (Optional) Target Mode of the Cloud L4-L7 Native Load Balancer object. Allowed values are "primary", "secondary", "unspecified", and default value is "unspecified". Type: String.
* `trunking` - (Optional) For virtual devices, if a trunking port group is to be used. Allowed values are "no", "yes", and default value is "no". Type: String.
* `cloud_l4l7_load_balancer_type` - (Optional) Type of the Cloud L4-L7 Native Load Balancer object. Allowed values are "application", "network", and default value is "application". Type: String.
* `relation_cloud_rs_ldev_to_cloud_subnet` - (Optional) Represents the relation to a Relation from Cloud LDev to Cloud Subnet (class cloudSubnet). Type: List.
* `aaa_domain_dn` - (Optional) Represents the relation to a Relation from AAA Domain to Cloud L4L7 Native Load Balancer (class aaaRbacAnnotation). Type: List.
* `static_ip_addresses` - (Optional) A list of unique static IP addresses of the Cloud L4-L7 Native Load Balancer object. Type: List.

## Importing ##

An existing Cloud L4-L7 Native Load Balancer can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_l4_l7_native_load_balancer.example <Dn>
```

Starting in Terraform version 1.5, an existing Cloud L4-L7 Native Load Balancer can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

 ```
 import {
    id = "<Dn>"
    to = aci_cloud_l4_l7_native_load_balancer.example
 }
 ```