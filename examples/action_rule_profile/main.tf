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

resource "aci_action_rule_profile" "example" {

  tenant_dn       = aci_tenant.example.id
  description     = "From Terraform"
  name            = "example"
  annotation      = "example"
  name_alias      = "example"
  set_route_tag   = 100
  set_preference  = 100
  set_weight      = 100
  set_metric      = 100
  set_metric_type = "ospf-type1"
}
