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

resource "aci_endpoint_ip_aging_profile" "example" {
  admin_st = "disabled"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
}