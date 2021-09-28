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
  tenant_dn      = aci_tenant.tenant01.id
  name           = "demo_l3out"
  annotation     = "tag_l3out"
  name_alias     = "alias_out"
  target_dscp    = "unspecified"
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

resource "aci_bgp_peer_connectivity_profile" "example" {
  logical_node_profile_dn = aci_logical_node_profile.logical_node_profile.id
  addr                    = "10.0.0.1"
  addr_t_ctrl             = "af-mcast,af-ucast"
  allowed_self_as_cnt     = "1"
  annotation              = "example"
  ctrl                    = "allow-self-as"
  name_alias              = "example"
  password                = "example"
  peer_ctrl               = "bfd"
  private_a_sctrl         = "remove-all,remove-exclusive"
  ttl                     = "7"
  weight                  = "5"
  as_number               = "27500"
  local_asn               = "10"
  local_asn_propagate     = "dual-as"
}
