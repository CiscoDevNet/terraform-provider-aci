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
resource "aci_leaf_interface_profile" "example" {
  description = "s"
  name        = "example"
  annotation  = "tag_leaf"
  name_alias  = "s"
}
resource "aci_access_port_selector" "example" {
  leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
  description               = "s"
  name                      = "example"
  access_port_selector_type = "ALL"
  annotation                = "tag_port_selector"
  name_alias                = "alias_port_selector"
}

resource "aci_access_sub_port_block" "fooaccess_sub_port_block" {
  access_port_selector_dn = aci_access_port_selector.example.id
  description             = "s"
  name                    = "example"
  annotation              = "example"
  from_card               = "1"
  from_port               = "1"
  from_sub_port           = "1"
  name_alias              = "example"
  to_card                 = "1"
  to_port                 = "1"
  to_sub_port             = "1"
}