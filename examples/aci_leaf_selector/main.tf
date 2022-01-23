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

resource "aci_leaf_selector" "example" {
  leaf_profile_dn         = aci_leaf_profile.example.id
  name                    = "example_leaf_selector"
  switch_association_type = "range"
  description             = "from terraform"
  relation_infra_rs_acc_node_p_grp = aci_access_switch_policy_group.example.id
}

resource "aci_leaf_selector" "example2" {
  leaf_profile_dn         = aci_leaf_profile.example.id
  name                    = "example2_leaf_selector"
  switch_association_type = "range"
}

resource "aci_leaf_profile" "example" {
  name        = "leaf-profile-example"
  description = "From Terraform"
  relation_infra_rs_acc_port_p = [
    aci_leaf_interface_profile.example.id
  ]
}

resource "aci_leaf_interface_profile" "example" {
  name = "leaf-interface-profile-example"
}

resource "aci_access_switch_policy_group" "example" {
  name  = "policy-group-example"
  description = "example"
}