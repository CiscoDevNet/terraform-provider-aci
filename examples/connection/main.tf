terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_connection" "conn" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.example.id
  name                            = "connection"
  adj_type                        = "L3"
  description                     = "from terraform"
  annotation                      = "example"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "yes"
  name_alias                      = "example"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.example.term_prov_dn,
    aci_function_node.example1.conn_provider_dn
  ]
}
