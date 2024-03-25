---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_function_node"
sidebar_current: "docs-aci-data-source-aci_function_node"
description: |-
  Data source for ACI Function Node
---

# aci_function_node

Data source for ACI Function Node

## API Information ##

* `Class` - vnsAbsNode
* `Distinguished Name` - uni/tn-{tenant_name}/AbsGraph-{sg_name}/AbsNode-{node_name}

## GUI Information ##

* `Location` - Tenants -> Services -> L4-L7 -> Service Graph Templates -> {Service_Graph} -> {Function_Node}

## Example Usage

```hcl
data "aci_function_node" "example" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.example.id
  name                            = "example"
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object. Type: String.
- `name` - (Required) Name of the Function Node object. Type: String.

## Attribute Reference

- `id` - (Read-Only) Attribute id set to the Dn of the Function Node object. Type: String.
- `annotation` - (Read-Only) Annotation of the Function Node object. Type: String.
- `description` - (Read-Only) Description of the Function Node object. Type: String.
- `func_template_type` - (Read-Only) Function Template type of the Function Node object. Type: String.
- `func_type` - (Read-Only) A function type. A GoThrough node is a transparent device, where a packet goes through without being addressed to the device, and the endpoints are not aware of that device. A GoTo device has a specific destination. Type: String.
- `is_copy` - (Read-Only) If the device is a copy device. Type: String.
- `managed` - (Read-Only) Specified if the function is using a managed device. Type: String.
- `name_alias` - (Read-Only) Name alias of the Function Node object. Type: String.
- `routing_mode` - (Read-Only) Routing mode of the Function Node object. Type: String.
- `sequence_number` - (Read-Only) Internal property incremented when aaa user logs in. Type: String.
- `share_encap` - (Read-Only) Enables encap sharing on node. Type: String.
- `relation_vns_rs_node_to_abs_func_prof` - (Read-Only) Represents the relation to L4-L7 Services Function Profile (class vnsAbsFuncProf). Type: String.
- `relation_vns_rs_node_to_l_dev` - (Read-Only) Represents the relation to Logical Device Abstraction (class vnsALDevIf). Type: String.
- `relation_vns_rs_node_to_m_func` - (Read-Only) Represents the relation to Meta Function (class vnsMFunc). Type: String.
- `relation_vns_rs_default_scope_to_term` - (Read-Only) Represents the relation to Terminal Abstract Class (class vnsATerm). Type: String.
- `relation_vns_rs_node_to_cloud_l_dev` - (Read-Only) Represents the relation to Cloud L4-L7 Abstract Devices (class cloudALDev). Type: String.
- `l4_l7_device_interface_consumer_name` - (Read-Only) The device interface is used to map with a service graph Function Node Connector consumer object. Type: String.
- `l4_l7_device_interface_consumer_connector_type` - (Read-Only) The device interface connector type used to map with a service graph Function Node Connector consumer object. Type: String.
- `l4_l7_device_interface_consumer_attachment_notification` - (Read-Only) Represents the consumer attachment notification. Type: String.
- `l4_l7_device_interface_provider_name` - (Read-Only) The device interface is used to map with a service graph Function Node Connector provider object. Type: String.
- `l4_l7_device_interface_provider_connector_type` - (Read-Only) The device interface connector type used to map with a service graph Function Node Connector provider object. Type: String.
- `l4_l7_device_interface_provider_attachment_notification` - (Read-Only) Represents the provider attachment notification. Type: String.
- `conn_consumer_dn` - (Read-Only) Distinguished name of the Function Node consumer connector. Type: String.
- `conn_provider_dn` - (Read-Only) Distinguished name of the Function Node provider connector. Type: String.
- `conn_copy_dn` - (Read-Only) Distinguished name of the Function Node copy connector. Type: String.