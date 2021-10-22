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

resource "aci_spine_port_selector" "example" {
  spine_profile_dn = aci_spine_profile.example.id
  tdn              = aci_spine_interface_profile.example.id
  annotation       = "example"
}