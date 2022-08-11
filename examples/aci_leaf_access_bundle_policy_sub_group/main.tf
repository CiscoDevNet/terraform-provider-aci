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

resource "aci_leaf_access_bundle_policy_group" "example" {
  name        = "example"
  description = "This policy group is created by terraform"
  lag_t       = "link"
}

resource "aci_lacp_member_policy" "example" {
  name        = "example"
  description = "This policy member is created by terraform"
}

resource "aci_leaf_access_bundle_policy_sub_group" "example" {
  leaf_access_bundle_policy_group_dn  = aci_leaf_access_bundle_policy_group.example.id
  name        = "example"
  description = "This policy group is created by terraform"
  port_channel_member = aci_lacp_member_policy.example.id
}