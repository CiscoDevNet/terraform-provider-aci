terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name = "terraform_tenant"
}

resource "aci_l3_outside" "l3_outside" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "l3_outside"
}

resource "aci_logical_node_profile" "logical_node_profile" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  name          = "logical_node_profile"
  config_issues = "none"
  tag           = "black"
  target_dscp   = "unspecified"
}

resource "aci_logical_interface_profile" "logical_interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.logical_node_profile.id
  name                    = "logical_interface_profile"
  prio                    = "unspecified"
  tag                     = "black"
}

resource "aci_l3out_path_attachment" "l3out_path_attachment" {
  logical_interface_profile_dn = aci_logical_interface_profile.logical_interface_profile.id
  target_dn                    = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t                    = "ext-svi"
  addr                         = "0.0.0.0"
  autostate                    = "disabled"
  encap                        = "vlan-1"
  encap_scope                  = "ctx"
  ipv6_dad                     = "disabled"
  ll_addr                      = "::"
  mac                          = "0F:0F:0F:0F:FF:FF"
  mode                         = "native"
  mtu                          = "inherit"
  target_dscp                  = "AF11"
}

resource "aci_l3out_path_attachment_secondary_ip" "l3out_path_attachment_secondary_ip" {
  l3out_path_attachment_dn = aci_l3out_path_attachment.l3out_path_attachment.id
  addr                     = "10.0.0.1/24"
  ipv6_dad                 = "disabled"
  dhcp_relay               = "enabled"
}
