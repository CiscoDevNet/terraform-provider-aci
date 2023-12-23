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

resource "aci_login_domain" "example" {
  name           = "example"
  annotation     = "orchestrator:terraform"
  provider_group = "example"
  realm          = "local"
  realm_sub_type = "default"
  description    = "From Terraform"
  name_alias     = "login_domain_alias"
}