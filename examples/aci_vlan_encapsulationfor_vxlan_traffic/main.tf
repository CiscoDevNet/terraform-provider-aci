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
resource "aci_attachable_access_entity_profile" "example" {
  description = "AAEP description"
  name        = "demo_entity_prof"
  annotation  = "tag_entity"
  name_alias  = "alias_entity"
}

resource "aci_vlan_encapsulationfor_vxlan_traffic" "foovlan_encapsulationfor_vxlan_traffic" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.example.id
  description                         = "From Terraform"
  annotation                          = "tag_traffic"
  name_alias                          = "name_alias"
}