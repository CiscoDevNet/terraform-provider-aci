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

data "aci_interface_blacklist" "check" {
  pod_id    = 1
  node_id   = 101
  interface = "eth1/1"
}

resource "aci_interface_blacklist" "disable" {
  pod_id    = 1
  node_id   = 101
  interface = "eth1/2"
}

resource "aci_interface_blacklist" "disable-string" {
  pod_id    = "1"
  node_id   = "101"
  interface = "eth1/3"
}

resource "aci_interface_blacklist" "disable-fex" {
  pod_id    = "1"
  node_id   = "101"
  fex_id    = "100"
  interface = "eth1/4"
}