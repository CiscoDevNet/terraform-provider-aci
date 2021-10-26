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

resource "aci_leaf_selector" "example" {
  leaf_profile_dn         = aci_leaf_profile.example.id
  name                    = "example_leaf_selector"
  switch_association_type = "range"
  annotation              = "orchestrator:terraform"
  description             = "from terraform"
  name_alias              = "tag_leaf_selector"
}