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

resource "aci_fex_bundle_group" "example" {
  fex_profile_dn = aci_fex_profile.example.id
  name           = "example"
  annotation     = "example"
  name_alias     = "example"
  description    = "from terraform"
}