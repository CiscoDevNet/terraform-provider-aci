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

resource "aci_login_domain_provider" "example" {
  parent_dn  = aci_duo_provider_group.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  order = "0"
  name_alias = "example_name_alias"
  description = "from terraform"
}