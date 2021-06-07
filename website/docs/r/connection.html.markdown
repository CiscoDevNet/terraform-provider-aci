---
layout: "aci"
page_title: "ACI: aci_connection"
sidebar_current: "docs-aci-resource-connection"
description: |-
  Manages ACI Connection
---

# aci_connection

Manages ACI Connection

## Example Usage

```hcl
resource "aci_connection" "conn2" {
  l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.example.id
  name  = "conn2"
  adj_type  = "L3"
  description = "from terraform"
  annotation  = "example"
  conn_dir  = "unknown"
  conn_type  = "internal"
  direct_connect  = "yes"
  name_alias  = "example"
  unicast_route  = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.example.term_cons_dn,
    aci_function_node.example.conn_consumer_dn
  ]
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
- `name` - (Required) Name of object connection.
- `adj_type` - (Optional) Connector adjacency type. Allowed values are "L2", "L3". Default value is "L2".
- `annotation` - (Optional) Annotation for object connection.
- `description` - (Optional) Description for object connection.
- `conn_dir` - (Optional) Connection directory for object connection. Allowed values are "consumer", "provider". Default value is "provider".
- `conn_type` - (Optional) Connection type of connection object. Allowed values are "external", "internal". Default value is "external".
- `direct_connect` - (Optional) Direct connect for object connection. Allowed values are "yes" and "no". Default value is "no".
- `name_alias` - (Optional) Name alias for object connection.
- `unicast_route` - (Optional) Unicast route for connection object. Allowed values are "yes" and "no". Default value is "yes".

- `relation_vns_rs_abs_copy_connection` - (Optional) List of relation to class vnsAConn. Cardinality - ONE_TO_M. Type - Set of String.
- `relation_vns_rs_abs_connection_conns` - (Optional) list of relation to class vnsAConn. Cardinality - ONE_TO_M. Type - Set of String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Connection.

## Importing

An existing Connection can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_connection.example <Dn>
```
