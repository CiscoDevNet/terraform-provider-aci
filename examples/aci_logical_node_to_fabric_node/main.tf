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


resource "aci_logical_node_to_fabric_node" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  tdn               = "topology/pod-1/paths-101/pathep-[eth1/4]"
  annotation        = "annotation"
  config_issues     = "none"
  rtr_id            = "10.0.1.1"
  rtr_id_loop_back  = "no"
}