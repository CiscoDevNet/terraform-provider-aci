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

resource "aci_tenant" "tenant_for_route_control" {
  name        = "tenant_for_route_control"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenant_for_route_control.id
  name                       = "route_profile01"
  description                = "from terraform"
  route_control_profile_type = "global"
}

resource "aci_action_rule_profile" "set_rule1" {
  tenant_dn = aci_tenant.tenant_for_route_control.id
  name      = "Rule01"
}

resource "aci_match_rule" "rule" {
  tenant_dn  = aci_tenant.tenant_for_route_control.id
  name  = "match_rule"
}

resource "aci_route_control_context" "control" {
  route_control_profile_dn  = aci_route_control_profile.example.id
  name  = "control"
  action = "permit"
  order = "0"
  set_rule = aci_action_rule_profile.set_rule1.id
  relation_rtctrl_rs_ctx_p_to_subj_p = [aci_match_rule.rule.id]
}

resource "aci_match_route_destination_rule" "destination" {
  match_rule_dn  = aci_match_rule.rule.id
  ip  = "10.1.1.1/24"
  aggregate = "yes"
  greater_than_mask = "25"
  less_than_mask = "26"
}