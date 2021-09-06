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

resource "aci_l3_interface_policy" "example" {
  name  = "example"
  annotation  = "example"
  bfd_isis = "disabled"
  name_alias  = "example"
  description = "example"
}