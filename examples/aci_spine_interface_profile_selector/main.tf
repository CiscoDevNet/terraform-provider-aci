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

resource "aci_spine_interface_profile_selector" "example" {
  spine_profile_dn = aci_spine_profile.example.id
  tdn              = aci_spine_interface_profile.example.id
  annotation       = "example"
}

resource "aci_spine_profile" "example" {
  name        = "spine_profile_1"
}

resource "aci_spine_interface_profile" "example" {
  name        = "example"
}