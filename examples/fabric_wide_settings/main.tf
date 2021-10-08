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
  enable_mo_streaming = "yes"
  enable_remote_leaf_direct = "yes"
  enforce_subnet_check = "yes"
  opflexp_authenticate_clients = "yes"
  opflexp_use_ssl = "yes"
  restrict_infra_vlan_traffic = "yes"
  unicast_xr_ep_learn_disable = "yes"
  validate_overlapping_vlans = "yes"
}