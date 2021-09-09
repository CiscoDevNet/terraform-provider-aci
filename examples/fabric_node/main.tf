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

data "aci_fabric_node" "example" {
  fabric_pod_dn  = "topology/pod-1"
  fabric_node_id  = "101"
}