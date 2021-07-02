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

resource "aci_vmm_domain" "foovmm_domain" {
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