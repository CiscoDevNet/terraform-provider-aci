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

data "aci_dhcp_option" "example" {
  dhcp_option_policy_dn  = aci_dhcp_option_policy.example.id
  name  = "example"
}