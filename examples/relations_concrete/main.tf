provider "aci" {
  username = "admin"
  password = "cisco123"
  url      = "https://192.168.10.102"
  insecure = true
}

resource "aci_tenant" "tenant_for_epg" {
  name        = "tenant_for_epg"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_application_profile" "app_profile_for_epg" {
  tenant_dn   = aci_tenant.tenant_for_epg.id
  name        = "ap_for_epg"
  description = "This app profile is created by terraform ACI providers"
}

resource "aci_application_epg" "inherit_epg" {
  application_profile_dn = aci_application_profile.app_profile_for_epg.id
  name                   = "inherit_epg"
  description            = "epg to create relation sec_inherited"

}

resource "aci_epg_to_static_path" "epg_to_stat_path" {
  application_epg_dn = aci_application_epg.inherit_epg.id
  tdn                = "topology/pod-1/paths-103/pathep-[eth1/1]"
  mode               = "regular"
  encap              = "vlan-1111"
  instr_imedcy       = "immediate"

}

resource "aci_l3_outside" "fool3_outside" {
  tenant_dn      = aci_tenant.tenant_for_epg.id
  name           = "demo_l3out"
  annotation     = "tag_l3out"
  name_alias     = "alias_out"
  target_dscp    = "unspecified"
}

resource "aci_logical_node_profile" "foological_node_profile" {
  l3_outside_dn = aci_l3_outside.fool3_outside.id
  description   = "sample logical node profile"
  name          = "demo_node"
  annotation    = "tag_node"
  config_issues = "none"
  name_alias    = "alias_node"
  tag           = "black"
  target_dscp   = "unspecified"
}

resource "aci_logical_node_to_fabric_node" "example" {

  logical_node_profile_dn  = aci_logical_node_profile.foological_node_profile.id
  tdn  = "topology/pod-1/node-101"
  annotation  = "example"
  rtr_id  = "12"
}