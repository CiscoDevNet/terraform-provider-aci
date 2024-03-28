---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_function_node"
sidebar_current: "docs-aci-resource-aci_function_node"
description: |-
  Manages ACI Function Node
---

# aci_function_node

Manages ACI Function Node

## API Information ##

* `Class` - vnsAbsNode
* `Distinguished Name` - uni/tn-{tenant_name}/AbsGraph-{sg_name}/AbsNode-{node_name}

## GUI Information ##

* `Location` - Tenants -> Services -> L4-L7 -> Service Graph Templates -> {Service_Graph} -> {Function_Node}

## Example Usage

```hcl
resource "aci_function_node" "example" {
  l4_l7_service_graph_template_dn      = aci_l4_l7_service_graph_template.example.id
  name                                 = "example"
  annotation                           = "example"
  description                          = "from terraform"
  func_template_type                   = "CLOUD_NATIVE_LB"
  func_type                            = "GoTo"
  is_copy                              = "yes"
  managed                              = "no"
  name_alias                           = "example"
  routing_mode                         = "Redirect"
  sequence_number                      = "1"
  share_encap                          = "yes"
  l4_l7_device_interface_consumer_name = "interface1"
  l4_l7_device_interface_provider_name = "interface2"
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object. Type: String.
- `name` - (Required) Name of the Function Node object. Type: String.
    - The valid function node name format for cloud APICs is `NX`, where X is a number starting with 0.
    - The valid function node name format for on-prem APICs is `NX`, where X is a number starting with 1.
    - The valid copy function node name format for on-premises APICs is `CPX`, where X is a number starting with 1.
- `annotation` - (Optional) Annotation of the Function Node object. Type: String.
- `description` - (Optional) Description of the Function Node object. Type: String.
- `func_template_type` - (Optional) Function Template type of the Function Node object. Allowed values: "OTHER", "FW_TRANS", "FW_ROUTED", "CLOUD_VENDOR_LB", "CLOUD_VENDOR_FW", "CLOUD_NATIVE_LB", "CLOUD_NATIVE_FW", "ADC_TWO_ARM", "ADC_ONE_ARM". Default value: "OTHER". Type: String.
- `func_type` - (Optional) A function type. A GoThrough node is a transparent device, where a packet goes through without being addressed to the device, and the endpoints are not aware of that device. A GoTo device has a specific destination. Allowed values: "GoThrough", "GoTo", "L1", "L2", "None". Default value: "GoTo". Type: String.
- `is_copy` - (Optional) If the device is a copy device. Allowed values: "yes", "no". Default value: "no". Type: String.
- `managed` - (Optional) Specified if the function is using a managed device. Allowed values: "yes", "no". Default value: "yes". In ACI version: 5.2 and greater, `managed` is not supported so we need to set it to "no". Type: String.
- `name_alias` - (Optional) Name alias of the Function Node object. Type: String.
- `routing_mode` - (Optional) Routing mode of the Function Node object. Allowed values: "Redirect", "unspecified". Default value: "unspecified". Type: String.
- `sequence_number` - (Optional) Internal property incremented when aaa user logs in. Type: String.
- `share_encap` - (Optional) Enables encap sharing on node. Allowed values: "yes", "no". Default value: "no". Type: String.
- `relation_vns_rs_node_to_abs_func_prof` - (Optional) Represents the relation to L4-L7 Services Function Profile (class vnsAbsFuncProf). Type: String.
- `relation_vns_rs_node_to_l_dev` - (Optional) Represents the relation to Logical Device Abstraction (class vnsALDevIf). Type: String.
- `relation_vns_rs_node_to_m_func` - (Optional) Represents the relation to Meta Function (class vnsMFunc). Type: String.
- `relation_vns_rs_default_scope_to_term` - (Optional) Represents the relation to Terminal Abstract Class (class vnsATerm). Type: String.
- `relation_vns_rs_node_to_cloud_l_dev` - (Optional) Represents the relation to Cloud L4-L7 Abstract Devices (class cloudALDev). Type: String.
- `l4_l7_device_interface_consumer_name` - (Optional) The device interface is used to map with a service graph Function Node Connector consumer object. Type: String.
- `l4_l7_device_interface_consumer_connector_type` - (Optional) The device interface connector type used to map with a service graph Function Node Connector consumer object. Allowed values: "none", "redir". Default value: "none". Type: String.
- `l4_l7_device_interface_consumer_attachment_notification` - (Optional) Represents the consumer attachment notification. Allowed values: "yes", "no". Default value: "no". Type: String.
- `l4_l7_device_interface_provider_name` - (Optional) The device interface is used to map with a service graph Function Node Connector provider object. Type: String.
- `l4_l7_device_interface_provider_connector_type` - (Optional) The device interface connector type used to map with a service graph Function Node Connector provider object. Allowed values: "none", "redir", "dnat", "snat", "snat_dnat". Default value: "none". Type: String.
- `l4_l7_device_interface_provider_attachment_notification` - (Optional) Represents the provider attachment notification. Allowed values: "yes", "no". Default value: "no". Type: String.

## Attribute Reference
- `conn_consumer_dn` - (Read-Only) Distinguished name of the Function Node consumer connector. Type: String.
- `conn_provider_dn` - (Read-Only) Distinguished name of the Function Node provider connector. Type: String.

## Importing

An existing Function Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_function_node.example <Dn>
```

Starting in Terraform version 1.5, an existing Function Node can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

 ```
 import {
    id = "<Dn>"
    to = aci_function_node.example
 }
 ```