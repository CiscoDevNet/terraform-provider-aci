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
  name = "demo_leaf_profile"
}
resource "aci_access_port_selector" "fooaccess_port_selector" {
  leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
  description               = "From Terraform"
  name                      = "tf_test"
  access_port_selector_type = "default"
  annotation                = "tag_port_selector"
  name_alias                = "alias_port_selector"
}

resource "aci_access_port_block" "test_port_block" {
  access_port_selector_dn = aci_access_port_selector.fooaccess_port_selector.id
  name                    = "tf_test_block"
  description             = "From Terraform"
  annotation              = "tag_port_block"
  from_card               = "1"
  from_port               = "1"
  name_alias              = "alias_port_block"
  to_card                 = "3"
  to_port                 = "3"
}