---
layout: "aci"
page_title: "ACI: aci_connection"
sidebar_current: "docs-aci-data-source-connection"
description: |-
  Data source for ACI Connection
---

# aci_connection #
Data source for ACI Connection

## Example Usage ##

```hcl
data "aci_connection" "check" {
  l4_l7_service_graph_template_dn  = "${aci_l4_l7_service_graph_template.example.id}"
  name  = "conn2"
}
```


## Argument Reference ##
* `l4_l7_service_graph_template_dn` - (Required) Distinguished name of parent L4-L7 Service Graph Template object.
* `name` - (Required) name of Object connection.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Connection.
* `adj_type` - (Optional) connector adjacency type
* `annotation` - (Optional) annotation for object connection.
* `conn_dir` - (Optional) conn_dir for object connection.
* `conn_type` - (Optional) 
* `direct_connect` - (Optional) direct_connect for object connection.
* `name_alias` - (Optional) name_alias for object connection.
* `unicast_route` - (Optional) unicast route
