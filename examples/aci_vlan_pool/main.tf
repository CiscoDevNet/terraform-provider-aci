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

resource "aci_vlan_pool" "example" {
  name        = "example"
  description = "From Terraform"
  alloc_mode  = "static"
  annotation  = "example"
  name_alias  = "example"
}