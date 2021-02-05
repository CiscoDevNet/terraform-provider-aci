---
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

  l4_l7_service_graph_template_dn  = "${aci_l4_l7_service_graph_template.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7ServiceGraphTemplate object.
* `name` - (Required) name of Object function_node.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Function Node.
* `annotation` - (Optional) annotation for object function_node.
* `func_template_type` - (Optional) func_template_type for object function_node.
* `func_type` - (Optional) function type
* `is_copy` - (Optional) is_copy for object function_node.
* `managed` - (Optional) managed for object function_node.
* `name_alias` - (Optional) name_alias for object function_node.
* `routing_mode` - (Optional) routing_mode for object function_node.
* `sequence_number` - (Optional) internal property incremented when aaa user logs in
* `share_encap` - (Optional) enables encap sharing on node
