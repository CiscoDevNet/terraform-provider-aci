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

resource "aci_tenant" "example" {
  name = "tf_l3out_tenant"
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
  set_next_hop    = "1.1.1.1"
  set_communities = {
    community = "no-advertise"
    criteria  = "replace"
  }
  next_hop_propagation    = "yes"
  multipath               = "yes"
  saspath_prepend_last_as = 10
  saspath_prepend_asn = {
    order = 20
    asn   = 30
  }
}
