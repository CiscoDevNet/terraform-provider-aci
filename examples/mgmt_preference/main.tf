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

resource "aci_mgmt_preference" "example" {
  interface_pref = "inband"
  annotation     = "orchestrator:terraform"
  description    = "from terraform"
  name_alias     = "example_name_alias"
}