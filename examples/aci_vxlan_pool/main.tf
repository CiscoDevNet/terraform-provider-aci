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

resource "aci_vxlan_pool" "example" {
  name        = "example"
  annotation  = "example"
  name_alias  = "name_alias_for_vxlan"
  description = "From Teraform"
}