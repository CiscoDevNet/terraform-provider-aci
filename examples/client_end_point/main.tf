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

data "aci_client_end_point" "check" {
  name               = "25:56:68:78:98:74"
  mac                = "25:56:68:78:98:74"
  ip                 = "0.0.0.0"
  vlan               = "5"
}
