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

resource "aci_spine_interface_profile" "example" {
  name        = "example"
  description = "from terraform"
  annotation  = "example"
  name_alias  = "example"
}