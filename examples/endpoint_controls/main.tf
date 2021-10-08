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

resource "aci_endpoint_control" "example" {
  admin_st = "disabled"
  annotation = "orchestrator:terraform"
  hold_intvl = "1800"
  rogue_ep_detect_intvl = "60"
  rogue_ep_detect_mult = "4"
  description = "from terraform"
  name_alias = "example_name_alias"
}