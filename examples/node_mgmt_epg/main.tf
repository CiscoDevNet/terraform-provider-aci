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

resource "aci_node_mgmt_epg" "in_band_example" {
  type = "in_band"
  management_profile_dn  = "uni/tn-mgmt/mgmtp-default"
  name  = "inb_example"
  annotation  = "example"
  encap  = "vlan-1"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
}

resource "aci_node_mgmt_epg" "out_of_band_example" {
  type = "out_of_band"
  management_profile_dn  = "uni/tn-mgmt/mgmtp-default"
  name  = "oob_example"
  annotation  = "example"
  name_alias  = "example"
  prio = "level1"
}
