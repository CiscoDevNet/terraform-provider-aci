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

resource "aci_l3out_floating_svi" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  node_dn                      = "topology/pod-1/node-201"
  encap                        = "vlan-20"
  addr                         = "10.20.30.40/16"
  annotation                   = "example"
  description                  = "from terraform"
  autostate                    = "enabled"
  encap_scope                  = "ctx"
  if_inst_t                    = "ext-svi"
  ipv6_dad                     = "disabled"
  ll_addr                      = "::"
  mac                          = "12:23:34:45:56:67"
  mode                         = "untagged"
  mtu                          = "580"
  target_dscp                  = "CS1"
  userdom                      = ":all:"
  relation_l3ext_rs_dyn_path_att {
    tdn = data.aci_physical_domain.dom.id
    floating_address = "10.20.30.254/16"
    forged_transmit = "Disabled"
    mac_change = "Disabled"
    promiscuous_mode = "Disabled"
  }
}
