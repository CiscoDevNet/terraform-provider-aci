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

## Example Usage

```hcl
data "aci_function_node" "example" {
  l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.example.id
  name  = "example"
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
- `name` - (Required) Name of Object function node.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Function Node.
- `annotation` - (Optional) Annotation for object function node.
- `description` - (Optional) Description for object function node.
- `func_template_type` - (Optional) Function Template type for object function node.
- `func_type` - (Optional) A function type. A GoThrough node is a transparent device, where a packet goes through without being addressed to the device, and the endpoints are not aware of that device. A GoTo device has a specific destination.
- `is_copy` - (Optional) If the device is a copy device.
- `managed` - (Optional) Specified if the function is using a managed device.
- `name_alias` - (Optional) Name alias for object function node.
- `routing_mode` - (Optional) Routing mode for object function node.
- `sequence_number` - (Optional) Internal property incremented when aaa user logs in.
- `share_encap` - (Optional) Enables encap sharing on node.
