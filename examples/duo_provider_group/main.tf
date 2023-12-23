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

resource "aci_duo_provider_group" "example" {
  name                 = "example"
  annotation           = "orchestrator:terraform"
  auth_choice          = "CiscoAVPair"
  ldap_group_map_ref   = "100"
  provider_type        = "radius"
  sec_fac_auth_methods = ["auto"]
  name_alias           = "example_name_alias"
  description          = "from terraform"
}