---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_function_node"
sidebar_current: "docs-aci-data-source-function_node"
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

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
- `name` - (Required) Name of the Function Node object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Function Node object.
- `annotation` - (Read-Only) Annotation of the Function Node object.
- `description` - (Read-Only) Description of the Function Node object.
- `func_template_type` - (Read-Only) Function Template type of the Function Node object.
- `func_type` - (Read-Only) A function type. A GoThrough node is a transparent device, where a packet goes through without being addressed to the device, and the endpoints are not aware of that device. A GoTo device has a specific destination.
- `is_copy` - (Read-Only) If the device is a copy device.
- `managed` - (Read-Only) Specified if the function is using a managed device.
- `name_alias` - (Read-Only) Name alias of the Function Node object.
- `routing_mode` - (Read-Only) Routing mode of the Function Node object.
- `sequence_number` - (Read-Only) Internal property incremented when aaa user logs in.
- `share_encap` - (Read-Only) Enables encap sharing on node.
- `l4_l7_device_interface_consumer_name` - (Read-Only) The device interface is used to map with a service graph Function Node Connector consumer object.
- `l4_l7_device_interface_provider_name` - (Read-Only) The device interface is used to map with a service graph Function Node Connector provider object.
- `conn_consumer_dn` - (Read-Only) Dn of the Function Node Connector consumer object.
- `conn_provider_dn` - (Read-Only) Dn of the Function Node Connector provider object.
- `relation_vns_rs_node_to_abs_func_prof` - (Read-Only) Relation to class vnsAbsFuncProf. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_node_to_l_dev` - (Read-Only) Relation to class vnsALDevIf. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_node_to_m_func` - (Read-Only) Relation to class vnsMFunc. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_default_scope_to_term` - (Read-Only) Relation to class vnsATerm. Cardinality - ONE_TO_ONE. Type - String.
- `relation_vns_rs_node_to_cloud_l_dev` - (Read-Only) Relation to class cloudALDev. Cardinality - N_TO_ONE. Type - String.