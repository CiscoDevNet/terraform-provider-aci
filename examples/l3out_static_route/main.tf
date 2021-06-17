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

resource "aci_l3out_static_route" "example" {
  fabric_node_dn = aci_logical_node_to_fabric_node.example.id
  ip             = "10.0.0.1"
  aggregate      = "no"
  annotation     = "example"
  name_alias     = "example"
  pref           = "1"
  rt_ctrl        = "bfd"
  description    = "from terraform"
}
