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

resource "aci_tenant" "footenant" {
	name 		= "tenant_bgp_protp"
	description = "tenant created while acceptance testing"
}

resource "aci_l3_outside" "fool3_outside" {
  name 		= "l3_bgp_protp"
  description = "l3_outside created while acceptance testing"
  tenant_dn = aci_tenant.footenant.id
}

resource "aci_logical_node_profile" "foological_node_profile" {
  name 		= "logical_node_profile"
  description = "logical_node_profile created while acceptance testing"
  l3_outside_dn = aci_l3_outside.fool3_outside.id
}

resource "aci_bgp_best_path_policy" "foobgp_best_path_policy" {
  tenant_dn  = aci_tenant.footenant.id
  name       = "rs_bgp_best_path_policy"
  annotation = "best_path_example"
  ctrl       = "asPathMultipathRelax"
  name_alias = "best_path_example"
}

resource "aci_bgp_timers" "foobgp_timer" {
  tenant_dn    = aci_tenant.footenant.id
  name         = "one"
  annotation   = "bgp_timer_example"
  gr_ctrl      = "helper"
  hold_intvl   = "189"
  ka_intvl     = "65"
  max_as_limit = "70"
  name_alias   = "bgp_timer_aliasing"
  stale_intvl  = "15"
}

resource "aci_l3out_bgp_protocol_profile" "foolbgp_protp" {
  logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
  annotation              = "bgp_protp_annotation"
  name_alias              = "bgp_protp_name_alias"
  name = "bgp_protp_test"
  relation_bgp_rs_best_path_ctrl_pol = aci_bgp_best_path_policy.foobgp_best_path_policy.id
  relation_bgp_rs_bgp_node_ctx_pol = aci_bgp_timers.foobgp_timer.id
}
