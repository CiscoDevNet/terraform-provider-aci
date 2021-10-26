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

resource "aci_spine_switch_association" "fooswitch_association" {
  spine_profile_dn              = aci_spine_profile.foospine_profile.id
  description                   = "from terraform"
  name                          = "spine_switch_association_1"
  spine_switch_association_type = "range"
  annotation                    = "spine_switch_association_tag"
  name_alias                    = "example"
}

data "aci_spine_switch_association" "example3" {
  spine_profile_dn              = aci_spine_profile.foospine_profile.id
  name                          = aci_spine_switch_association.fooswitch_association.name
  spine_switch_association_type = aci_spine_switch_association.fooswitch_association.spine_switch_association_type
}

output "name3" {
  value = data.aci_spine_switch_association.example3
}