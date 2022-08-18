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

resource "aci_fabric_wide_settings" "example" {
  name = "example"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
  disable_ep_dampening = "yes"
  domain_validation = "yes"
  enable_mo_streaming = "yes"
  enable_remote_leaf_direct = "yes"
  enforce_subnet_check = "yes"
  leaf_opflexp_authenticate_clients = "yes"
  leaf_opflexp_use_ssl = "yes"
  opflexp_authenticate_clients = "yes"
  opflexp_ssl_protocols = "TLSv1,TLSv1.1,TLSv1.2"
  opflexp_use_ssl = "yes"
  policy_sync_node_bringup = "yes"
  reallocate_gipo = "yes"
  restrict_infra_vlan_traffic = "yes"
  unicast_xr_ep_learn_disable = "yes"
  validate_overlapping_vlans = "yes"
}