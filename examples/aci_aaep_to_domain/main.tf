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

resource "aci_attachable_access_entity_profile" "foo_attachable_access_entity_profile" {
  name        = "aaep_test"
  description = "attachable_access_entity_profile created while acceptance testing"
}

resource "aci_l3_domain_profile" "foo_l3_domain_profile" {
  name = "l3out_domain_profile_test"
}

resource "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.foo_attachable_access_entity_profile.id
  domain_dn                           = aci_l3_domain_profile.foo_l3_domain_profile.id
}

data "aci_aaep_to_domain" "data_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_aaep_to_domain.foo_aaep_to_domain.attachable_access_entity_profile_dn
  domain_dn                           = aci_aaep_to_domain.foo_aaep_to_domain.domain_dn
}
