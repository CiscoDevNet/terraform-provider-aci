terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://173.36.219.25"
  insecure = true
}


resource "aci_vmm_domain" "vmm_domain_vds_01" {
//   // access_mode         = "read-write"
//   // ave_time_out        = "30"
//   // config_infra_pg     = "no"
//   // ctrl_knob           = "epDpVerify"
//   // enable_ave          = "no"
//   // enable_tag          = "no"
//   // encap_mode          = "unknown"
//   // enf_pref            = "hw"
//   // ep_inventory_type   = "on-link"
//   // ep_ret_time         = "0"
//   // hv_avail_monitor    = "no"
//   // mcast_addr          = "0.0.0.0"
//   // mode                = "default"
  name                = "vds_01"
//   // pref_encap_mode     = "unspecified"
  provider_profile_dn = "uni/vmmp-VMware"
  relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
}

resource "aci_vlan_pool" "vmm_vlan_pool" {
  name       = "vlan_pool_P12"
  alloc_mode = "dynamic"
}

// // VMWare vmm domain resources
// resource "aci_vmm_domain" "vds" {
//   provider_profile_dn = var.vds
//   relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
//   name                = var.vmm_domain
// }

