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

resource "aci_l2_domain" "fool2_domain" {
  name       = "l2_domain_1"
  annotation = "l2_domain_tag"
  name_alias = "l2_domain"
}

data "aci_l2_domain" "example7" {
  name = aci_l2_domain.fool2_domain.name
}

output "name7" {
  value = data.aci_l2_domain.example7
}