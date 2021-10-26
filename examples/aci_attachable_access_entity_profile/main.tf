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
  name_alias  = "Name Alias"
}