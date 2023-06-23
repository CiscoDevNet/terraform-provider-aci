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

resource "aci_leaf_access_bundle_policy_group" "test_access_bundle_policy_group" {
  description = "From Terraform"
  name        = "tf_bundle_group"
  annotation  = "tag_if_pol"
  lag_t       = "link"
  name_alias  = "alias_if_pol"
}

resource "aci_leaf_access_bundle_policy_sub_group" "test_access_bundle_policy_sub_group" {
  leaf_access_bundle_policy_group_dn = aci_leaf_access_bundle_policy_group.test_access_bundle_policy_group.id
  name                               = "tf_sub_group"
}

resource "aci_access_port_selector" "fooaccess_port_selector" {
  leaf_interface_profile_dn      = aci_leaf_interface_profile.example.id
  description                    = "From Terraform"
  name                           = "tf_test"
  access_port_selector_type      = "range"
  annotation                     = "tag_port_selector"
  name_alias                     = "alias_port_selector"
  relation_infra_rs_acc_base_grp = aci_leaf_access_bundle_policy_group.test_access_bundle_policy_group.id
}

resource "aci_access_port_block" "test_port_block" {
  access_port_selector_dn           = aci_access_port_selector.fooaccess_port_selector.id
  name                              = "tf_test_block"
  description                       = "From Terraform"
  annotation                        = "tag_port_block"
  from_card                         = "1"
  from_port                         = "1"
  name_alias                        = "alias_port_block"
  to_card                           = "3"
  to_port                           = "3"
  relation_infra_rs_acc_bndl_subgrp = aci_leaf_access_bundle_policy_sub_group.test_access_bundle_policy_sub_group.id
}