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

resource "aci_access_switch_policy_group" "example" {
  name        = "example"
  annotation  = "example"
  description = "example"
  name_alias  = "example"
}