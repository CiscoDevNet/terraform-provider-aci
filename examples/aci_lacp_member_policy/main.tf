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

resource "aci_lacp_member_policy" "example" {
  name          = "example"
  description   = "This policy member is created by terraform"
  priority      = "32768"
  transmit_rate = "normal"
}