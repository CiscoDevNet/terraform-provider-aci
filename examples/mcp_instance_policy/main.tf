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

resource "aci_mcp_instance_policy" "example" {
  admin_st         = "disabled"
  annotation       = "orchestrator:terraform"
  name_alias       = "mcp_instance_alias"
  description      = "From Terraform"
  ctrl             = []
  init_delay_time  = "180"
  key              = "example"
  loop_detect_mult = "3"
  loop_protect_act = "port-disable"
  tx_freq          = "2"
  tx_freq_msec     = "0"
}