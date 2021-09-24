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

resource "aci_node_mgmt_epg" "example" {
  type = "in_band"
  description = "from terraform"
  name  = "example"
  annotation  = "example"
  encap  = "vlan-1"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
}

resource "aci_node_mgmt_epg" "example2" {
  type = "out_of_band"
  description = "from terraform"
  name  = "example"
  annotation  = "example"
  encap  = "vlan-1"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
}

resource "aci_mgmt_zone" "example" {
  managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.example.id
  type = "in_band"
  name = "inb_zone"
  name_alias = "zone_tag"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  relation_mgmt_rs_addr_inst = "uni/tn-mgmt/addrinst-Zabbix_SNMPoobaddr"

  relation_mgmt_rs_in_b = aci_node_mgmt_epg.example.id // type = "in_band"

  relation_mgmt_rs_inb_epg = aci_node_mgmt_epg.example.id // type = "in_band"
}

resource "aci_mgmt_zone" "example2" {
  managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.example.id
  type = "out_of_band"
  name = "oob_zone"
  name_alias = "zone_tag"
  annotation = "orchestrator:terraform"
  description = "from terraform"

  relation_mgmt_rs_oo_b = aci_node_mgmt_epg.example.id

  relation_mgmt_rs_oob_epg = aci_node_mgmt_epg.example.id
}
