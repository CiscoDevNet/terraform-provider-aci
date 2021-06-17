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

resource "aci_bgp_peer_connectivity_profile" "example" {
  logical_node_profile_dn = aci_logical_node_profile.example.id
  addr                    = "10.0.0.1"
  addr_t_ctrl             = "af-mcast,af-ucast"
  allowed_self_as_cnt     = "3"
  annotation              = "example"
  ctrl                    = "allow-self-as"
  name_alias              = "example"
  password                = "example"
  peer_ctrl               = "bfd"
  private_a_sctrl         = "remove-all,remove-exclusive"
  ttl                     = "1"
  weight                  = "1"
  as_number               = "1"
  local_asn               = "15"
  local_asn_propagate     = "dual-as"
}
