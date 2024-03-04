---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_connection"
sidebar_current: "docs-aci-data-source-aci_connection"
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

- `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object. Type: String.
- `name` - (Required) Name of Object connection. Type: String.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Connection. Type: String.
- `adj_type` - Connector adjacency type. Type: String.
- `annotation` - Annotation for object connection. Type: String.
- `description` - Description for object connection. Type: String.
- `conn_dir` - Connection direction for object connection. Type: String.
- `conn_type` - Connection type for object connection. Type: String.
- `direct_connect` - Direct connect for object connection. Type: String.
- `name_alias` - Name alias for object connection. Type: String.
- `unicast_route` - Unicast route for object connection. Type: String.
