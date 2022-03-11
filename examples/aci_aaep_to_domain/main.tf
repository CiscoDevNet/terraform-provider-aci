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

resource "aci_attachable_access_entity_profile" "fooattachable_access_entity_profile" {
  name        = "aaep_test"
  description = "attachable_access_entity_profile created while acceptance testing"

}

output "test_aaep_name" {
  value = resource.aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
}

resource "aci_l3_domain_profile" "fool3_domain_profile" {
  name       = "l3out_domain_profile_test"
  annotation = "example"
  name_alias = "example"
}

output "test_l3out_domain_profile" {
  value = resource.aci_l3_domain_profile.fool3_domain_profile.id
}

resource "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = resource.aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  t_dn                                = resource.aci_l3_domain_profile.fool3_domain_profile.id
}

output "test_aaep_to_domain" {
  value = resource.aci_aaep_to_domain.foo_aaep_to_domain.id
}
