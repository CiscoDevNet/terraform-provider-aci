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

resource "aci_spine_profile" "foospine_profile" {
  name        = "spine_profile_1"
  description = "from terraform"
  annotation  = "spine_profile_tag"
  name_alias  = "example"
}