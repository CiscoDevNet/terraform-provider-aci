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

resource "aci_cdp_interface_policy" "example" {
  name        = "example"
  admin_st    = "enabled"
  annotation  = "tag_cdp"
  name_alias  = "alias_cdp"
  description = "From Terraform"
}