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

resource "aci_aaa_domain" "foosecurity_domain" {
  name        = "aaa_domain_1"
  description = "from terraform"
  annotation  = "aaa_domain_tag"
  name_alias  = "example"
}

data "aci_aaa_domain" "example6" {
  name = aci_aaa_domain.foosecurity_domain.name
}

output "name6" {
  value = data.aci_aaa_domain.example6
}