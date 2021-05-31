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
	name        = "demo_leaf_profile"
}	
resource "aci_access_port_selector" "fooaccess_port_selector" {
	leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
	description               = "From Terraform"
	name                      = "tf_test"
	access_port_selector_type = "default"
	annotation                = "tag_port_selector"
	name_alias                = "alias_port_selector"
} 
