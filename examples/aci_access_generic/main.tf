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
  description = "From Terraform"
  name        = "demo_entity_prof"
  annotation  = "tag_entity"
  name_alias  = "Name_Alias"
}

resource "aci_access_generic" "fooaccess_generic" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  description                         = "from terraform"
  name                                = "default"
  annotation                          = "access_generic_tag"
  name_alias                          = "access_generic"
}

data "aci_access_generic" "example" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  name                                = aci_access_generic.fooaccess_generic.name
}

output "name" {
  value = data.aci_access_generic.example
}

