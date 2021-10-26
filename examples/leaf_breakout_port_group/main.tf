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

resource "aci_leaf_breakout_port_group" "example" {
  name       = "first"
  annotation = "example"
  brkout_map = "100g-4x"
  name_alias = "asfa"
}
