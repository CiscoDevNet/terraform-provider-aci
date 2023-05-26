
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

resource "aci_vmm_domain" "foo_vmm_domain" {
  provider_profile_dn = "uni/vmmp-VMware"
  name                = "demo_domp"
  access_mode         = "read-write"
  annotation          = "tag_dom"
  arp_learning        = "disabled"
  ave_time_out        = "30"
  config_infra_pg     = "no"
  ctrl_knob           = "epDpVerify"
  delimiter           = "_"
  enable_ave          = "no"
  enable_tag          = "no"
  enable_vm_folder    = "no"
  encap_mode          = "unknown"
  enf_pref            = "hw"
  ep_inventory_type   = "on-link"
  ep_ret_time         = "0"
  hv_avail_monitor    = "no"
  mcast_addr          = "224.0.1.2"
  mode                = "default"
  name_alias          = "alias_dom"
  pref_encap_mode     = "unspecified"
}

resource "aci_vmm_controller" "example" {
  vmm_domain_dn       = aci_vmm_domain.foo_vmm_domain.id
  name                = "example"
  annotation          = "orchestrator:terraform"
  dvs_version         = "unmanaged"
  host_or_ip          = "10.10.10.10"
  inventory_trig_st   = "untriggered"
  mode                = "default"
  msft_config_err_msg = "Error"
  msft_config_issues  = ["zero-mac-in-inventory", "aaacert-invalid"]
  n1kv_stats_mode     = "enabled"
  port                = "0"
  root_cont_name      = "vmmdc"
  scope               = "vm"
  seq_num             = "0"
  stats_mode          = "disabled"
  vxlan_depl_pref     = "vxlan"
}
