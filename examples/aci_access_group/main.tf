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
  name = "foo_leaf_int_prof"
}

resource "aci_access_port_selector" "example" {
  leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
  description               = "from terraform"
  name                      = "demo_port_selector"
  access_port_selector_type = "ALL"
}

resource "aci_fex_profile" "example" {
  name = "fex_prof"
}

resource "aci_fex_bundle_group" "example" {
  fex_profile_dn = aci_fex_profile.example.id
  name           = "example"
}

resource "aci_access_group" "example" {
  access_port_selector_dn = aci_access_port_selector.example.id
  annotation              = "one"
  fex_id                  = "101"
  tdn                     = aci_fex_bundle_group.example.id
}
