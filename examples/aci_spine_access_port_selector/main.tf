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

resource "aci_spine_interface_profile" "example" {
  name        = "example"
  description = "from terraform"
  annotation  = "example"
  name_alias  = "example"
}

resource "aci_spine_port_policy_group" "example" {
  name = "example"
}

resource "aci_spine_access_port_selector" "example" {
  spine_interface_profile_dn      = aci_spine_interface_profile.example.id
  name                            = "example"
  spine_access_port_selector_type = "range"

  relation_infra_rs_sp_acc_grp = aci_spine_port_policy_group.example.id
}

resource "aci_spine_access_port_selector" "example2" {
  spine_interface_profile_dn      = aci_spine_interface_profile.example.id
  name                            = "example2"
  spine_access_port_selector_type = "range"
}

resource "aci_access_port_block" "fooaccess_port_block" {
  access_port_selector_dn = aci_spine_access_port_selector.example.id
  description             = "from terraform"
  name                    = "demo_port_block"
  from_card               = "1"
  from_port               = "1"
  to_card                 = "3"
  to_port                 = "3"
}