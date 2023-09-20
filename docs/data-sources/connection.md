---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_connection"
sidebar_current: "docs-aci-data-source-connection"
description: |-
  Data source for ACI Connection
---

# aci_connection

Data source for ACI Connection

## Example Usage

```hcl
data "aci_connection" "check" {
  l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.example.id
  name  = "conn2"
}
```

## Argument Reference

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
- `name` - (Required) Name of Object connection.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Connection.
- `adj_type` - Connector adjacency type.
- `annotation` - Annotation for object connection.
- `description` - Description for object connection.
- `conn_dir` - Connection direction for object connection.
- `conn_type` - Connection type for object connection.
- `direct_connect` - Direct connect for object connection.
- `name_alias` - Name alias for object connection.
- `unicast_route` - Unicast route for object connection.
