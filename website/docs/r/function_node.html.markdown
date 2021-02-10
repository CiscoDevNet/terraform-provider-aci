---
layout: "aci"
page_title: "ACI: aci_function_node"
sidebar_current: "docs-aci-resource-function_node"
description: |-
  Manages ACI Function Node
---

# aci_function_node

Manages ACI Function Node

## Example Usage

```hcl
resource "aci_function_node" "example" {
  l4_l7_service_graph_template_dn  = "${aci_l4-l7_service_graph_template.example.id}"
  name  = "example"
  annotation  = "example"
  func_template_type  = "example"
  func_type  = "example"
  is_copy  = "example"
  managed  = "example"
  name_alias  = "example"
  routing_mode  = "example"
  sequence_number  = "example"
  share_encap  = "example"
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
- `name` - (Required) name of Object function_node.
- `annotation` - (Optional) Annotation for object function_node.
- `func_template_type` - (Optional) func_template_type for object function_node.
  Allowed values: "OTHER","FW_TRANS","FW_ROUTED","CLOUD_VENDOR_LB","CLOUD_VENDOR_FW","CLOUD_NATIVE_LB","CLOUD_NATIVE_FW","ADC_TWO_ARM","ADC_ONE_ARM". Default value: "OTHER"
- `func_type` - (Optional) Function type
  Allowed values: "GoThrough","GoTo","L1","L2","None". Default value: "GoTo"
- `is_copy` - (Optional) is_copy for object function_node.
  Allowed values: "yes","no". Default value: "no"
- `managed` - (Optional) Managed for object function_node.
  Allowed values: "yes","no". Default value: "yes"
- `name_alias` - (Optional) name_alias for object function_node.
- `routing_mode` - (Optional) Routing_mode for object function_node.
  Allowed values: "Redirect","unspecified". Default value: "unspecified"
- `sequence_number` - (Optional) Internal property incremented when aaa user logs in
- `share_encap` - (Optional) Enables encap sharing on node
  Allowed values: "yes","no". Default value: "no"

- `relation_vns_rs_node_to_abs_func_prof` - (Optional) Relation to class vnsAbsFuncProf. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_node_to_l_dev` - (Optional) Relation to class vnsALDevIf. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_node_to_m_func` - (Optional) Relation to class vnsMFunc. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_default_scope_to_term` - (Optional) Relation to class vnsATerm. Cardinality - ONE_TO_ONE. Type - String.
- `relation_vns_rs_node_to_cloud_l_dev` - (Optional) Relation to class cloudALDev. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Function Node.

## Importing

An existing Function Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_function_node.example <Dn>
```
