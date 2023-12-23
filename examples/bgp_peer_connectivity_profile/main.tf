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
resource "aci_tenant" "tenant01" {
  name        = "tenant01"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_l3_outside" "l3_outside" {
  tenant_dn   = aci_tenant.tenant01.id
  name        = "demo_l3out"
  annotation  = "tag_l3out"
  name_alias  = "alias_out"
  target_dscp = "unspecified"
}

resource "aci_logical_node_profile" "logical_node_profile" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  description   = "sample logical node profile"
  name          = "demo_node"
  annotation    = "tag_node"
  config_issues = "none"
  name_alias    = "alias_node"
  tag           = "black"
  target_dscp   = "unspecified"
}

# Logical Node level BGP Peer
resource "aci_bgp_peer_connectivity_profile" "node_bgp_peer" {
  parent_dn           = aci_logical_node_profile.logical_node_profile.id
  addr                = "10.0.0.1"
  addr_t_ctrl         = ["af-mcast", "af-ucast"]
  allowed_self_as_cnt = "1"
  annotation          = "example"
  ctrl                = ["allow-self-as"]
  name_alias          = "node_bgp_peer"
  password            = "example"
  peer_ctrl           = ["bfd"]
  private_a_sctrl     = ["remove-all", "remove-exclusive"]
  ttl                 = "7"
  weight              = "5"
  as_number           = "27500"
  local_asn           = "10"
  local_asn_propagate = "dual-as"
  admin_state         = "enabled"

  relation_bgp_rs_peer_to_profile {
    direction = "import"
    target_dn = "uni/tn-tenant01/prof-test"
  }
  relation_bgp_rs_peer_to_profile {
    direction = "export"
    target_dn = "uni/tn-tenant01/prof-data"
  }
}

resource "aci_logical_interface_profile" "logical_interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.logical_node_profile.id
  description             = "aci_logical_interface_profile from terraform"
  name                    = "demo_int_prof"
  annotation              = "tag_prof"
  name_alias              = "alias_prof"
  prio                    = "unspecified"
  tag                     = "black"
}

resource "aci_l3out_path_attachment" "l3out_path_attach" {
  logical_interface_profile_dn = aci_logical_interface_profile.logical_interface_profile.id
  target_dn                    = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t                    = "ext-svi"
  addr                         = "10.0.0.254/24"
  annotation                   = "example"
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

# Logical Interface Path Attchement level BGP Peer
resource "aci_bgp_peer_connectivity_profile" "interface_bgp_peer" {
  parent_dn           = aci_l3out_path_attachment.l3out_path_attach.id
  addr                = "10.0.0.2"
  addr_t_ctrl         = ["af-mcast", "af-ucast"]
  allowed_self_as_cnt = "1"
  annotation          = "example"
  ctrl                = ["allow-self-as"]
  name_alias          = "interface_bgp_peer"
  password            = "example"
  peer_ctrl           = ["bfd"]
  private_a_sctrl     = ["remove-all", "remove-exclusive"]
  ttl                 = "7"
  weight              = "5"
  as_number           = "27500"
  local_asn           = "10"
  local_asn_propagate = "dual-as"
  admin_state         = "enabled"

  relation_bgp_rs_peer_to_profile {
    direction = "import"
    target_dn = "uni/tn-tenant01/prof-test"
  }
  relation_bgp_rs_peer_to_profile {
    direction = "export"
    target_dn = "uni/tn-tenant01/prof-data"
  }
}
