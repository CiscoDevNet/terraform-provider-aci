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
  fabric_pod_dn  = aci_fabric_pod.example.id
  fabric_node_id  = "example"
}