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

resource "aci_tenant" "foo_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "terraform_test_tenant"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.foo_tenant.id
  name      = "vrf1"
}

resource "aci_l3_outside" "foo_l3_outside" {
  tenant_dn              = aci_tenant.foo_tenant.id
  name                   = "l3_outside"
  relation_l3ext_rs_ectx = aci_vrf.vrf1.id
}

resource "aci_l3out_bgp_external_policy" "l3out_bgp" {
  l3_outside_dn = aci_l3_outside.foo_l3_outside.id
}

resource "aci_logical_node_profile" "foo_logical_node_profile" {
  l3_outside_dn = aci_l3out_bgp_external_policy.l3out_bgp.l3_outside_dn
  name          = "demo_node"
}

resource "aci_logical_interface_profile" "foo_logical_interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.foo_logical_node_profile.id
  name                    = "demo_int_prof"
}

resource "aci_bfd_multihop_interface_policy" "foo_bfd_multihop_interface_policy" {
  tenant_dn = aci_tenant.foo_tenant.id
  name      = "bfd_mh_interface_policy"
}

resource "aci_bfd_multihop_interface_profile" "foo_bfd_multihop_interface_profile" {
  logical_interface_profile_dn = aci_logical_interface_profile.foo_logical_interface_profile.id
  interface_profile_type       = "sha1"
  key                          = "key123"
  key_id                       = "101"
  relation_bfd_rs_mh_if_pol    = aci_bfd_multihop_interface_policy.foo_bfd_multihop_interface_policy.id
}

data "aci_bfd_multihop_interface_profile" "data_bfd_multihopinterface_profile" {
  logical_interface_profile_dn = aci_bfd_multihop_interface_profile.foo_bfd_multihop_interface_profile.logical_interface_profile_dn
}

output "bfd_multihopinterface_profile" {
  value = data.aci_bfd_multihop_interface_profile.data_bfd_multihopinterface_profile
}
