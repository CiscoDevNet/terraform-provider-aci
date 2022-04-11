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

resource "aci_tenant" "foo_tenant" {
  name        = "tf_tenant"
  description = "tenant created while acceptance testing"
}

resource "aci_action_rule_profile" "foo_action_rule_profile" {
  name        = "tf_rule"
  description = "action_rule_profile created while acceptance testing"
  tenant_dn   = aci_tenant.foo_tenant.id
}

resource "aci_action_rule_profile_additional_communities" "foo_rtctrl_set_add_comm" {
  community              = "no-advertise"
  description            = "additional communities created while acceptance testing"
  action_rule_profile_dn = aci_action_rule_profile.foo_action_rule_profile.id
}
