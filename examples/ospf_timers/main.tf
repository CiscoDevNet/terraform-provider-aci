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

resource "aci_tenant" "tenentcheck" {
  name = "tenentcheck"
}

resource "aci_ospf_timers" "example" {
  tenant_dn           = aci_tenant.tenentcheck.id
  name                = "ospf_timers_1"
  annotation          = "ospf_timers_tag"
  description         = "from terraform"
  bw_ref              = "30000"
  ctrl                = ["pfx-suppress", "name-lookup"]
  dist                = "200"
  gr_ctrl             = "helper"
  lsa_arrival_intvl   = "2000"
  lsa_gp_pacing_intvl = "50"
  lsa_hold_intvl      = "1000"
  lsa_max_intvl       = "1000"
  lsa_start_intvl     = "5"
  max_ecmp            = "10"
  max_lsa_action      = "restart"
  max_lsa_num         = "56"
  max_lsa_reset_intvl = "10"
  max_lsa_sleep_cnt   = "10"
  max_lsa_sleep_intvl = "10"
  max_lsa_thresh      = "50"
  name_alias          = "example"
  spf_hold_intvl      = "100"
  spf_init_intvl      = "500"
  spf_max_intvl       = "10"
}
