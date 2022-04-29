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

resource "aci_function_node" "foofunction_node" {

  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.serviceGraphTemp.id
  name                            = "functionNodeOne"
  func_template_type              = "OTHER"
  func_type                       = "None"
  is_copy                         = "no"
  managed                         = "no"
  routing_mode                    = "unspecified"
  sequence_number                 = "3"
  share_encap                     = "yes"
}
