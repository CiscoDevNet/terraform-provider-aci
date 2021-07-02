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

resource "aci_fex_profile" "example" {
  name        = "fex_prof"
  annotation  = "example"
  name_alias  = "example"
  description = "from terraform"
}