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

resource "aci_tenant" "example" {
  name = "tf_tenant_l3out"
}

resource "aci_l3_outside" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "demo_l3out"
}

resource "aci_logical_node_profile" "node_profile" {
  l3_outside_dn = aci_l3_outside.example.id
  name          = "demo_node"
}

resource "aci_logical_interface_profile" "interface_profile" {
  logical_node_profile_dn = aci_logical_node_profile.node_profile.id
  name                    = "demo_int_prof"
}

resource "aci_l3out_floating_svi" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.interface_profile.id
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
  relation_l3ext_rs_dyn_path_att {
    tdn              = aci_physical_domain.example.id
    floating_address = "10.20.30.254/16"
    forged_transmit  = "Disabled"
    mac_change       = "Enabled"
    promiscuous_mode = "Disabled"
  }
  relation_l3ext_rs_dyn_path_att {
    tdn              = aci_vmm_domain.example_vmm.id
    floating_address = "10.20.30.254/16"
    enhanced_lag_policy_tdn = "uni/vmmp-VMware/dom-example/vswitchpolcont/enlacplagp-test"
  }
}

resource "aci_l3out_floating_svi" "example2" {
  logical_interface_profile_dn = aci_logical_interface_profile.interface_profile.id
  node_dn                      = "topology/pod-1/node-202"
  encap                        = "vlan-21"
  addr                         = "10.21.30.40/16"
  autostate                    = "enabled"
  encap_scope                  = "local"
  if_inst_t                    = "ext-svi"
  mtu                          = "9000"
  relation_l3ext_rs_dyn_path_att {
    tdn              = aci_physical_domain.example.id
    floating_address = "10.21.30.254/16"
  }
}

resource "aci_physical_domain" "example" {
  name = "example"
}

resource "aci_vmm_domain" "example_vmm" {
  provider_profile_dn = "uni/vmmp-VMware"
  name = "example"
}
