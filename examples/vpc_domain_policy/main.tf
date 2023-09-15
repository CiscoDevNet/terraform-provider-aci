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

resource "aci_vpc_domain_policy" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  dead_intvl  = "200"
  name_alias  = "example"
  description = "from terraform"
}