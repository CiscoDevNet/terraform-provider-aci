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

resource "aci_user_security_domain" "example" {
  local_user_dn  = aci_local_user.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
  description = "from Terraform"
}