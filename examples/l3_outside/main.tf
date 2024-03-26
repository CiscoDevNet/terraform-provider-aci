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


# Tenant Setup
resource "aci_tenant" "terraform_tenant" {
  name = "terraform_tenant"
}

data "aci_tenant" "common" {
  name = "common"
}

# Route Control Profile Setup
data "aci_route_control_profile" "shared_route_control_profile" {
  parent_dn = data.aci_tenant.common.id
  name      = "ok"
}

# VRF Setup
data "aci_vrf" "default_vrf" {
  tenant_dn = data.aci_tenant.common.id
  name      = "default"
}

# L3 Domain Setup
resource "aci_l3_domain_profile" "l3_domain_profile" {
  name = "l3_domain_profile"
}

resource "aci_l3_outside" "foo_l3_outside" {
  tenant_dn      = aci_tenant.terraform_tenant.id
  name           = "foo_l3_outside"
  enforce_rtctrl = ["export", "import"]
  target_dscp    = "unspecified"
  mpls_enabled   = "yes"
  pim            = ["ipv4", "ipv6"]
  // Relation to Route Control for Dampening - can't configure multiple Dampening Policies for the same address-family.
  relation_l3ext_rs_dampening_pol {
    tn_rtctrl_profile_dn = data.aci_route_control_profile.shared_route_control_profile.id
    af                   = "ipv6-ucast"
  }

  relation_l3ext_rs_dampening_pol {
    tn_rtctrl_profile_dn = data.aci_route_control_profile.shared_route_control_profile.id
    af                   = "ipv4-ucast"
  }

  // Target VRF object should belong to the parent tenant or be a shared object.
  relation_l3ext_rs_ectx = data.aci_vrf.default_vrf.id

  // Relation to Route Profile for Interleak - L3 Out Context Interleak Policy object should belong to the parent tenant or be a shared object.
  relation_l3ext_rs_interleak_pol = data.aci_route_control_profile.shared_route_control_profile.id

  // Relation to L3 Domain
  relation_l3ext_rs_l3_dom_att = aci_l3_domain_profile.l3_domain_profile.id

  // Relation to Route Profile for Redistribution
  relation_l3extrs_redistribute_pol {
    target_dn = data.aci_route_control_profile.shared_route_control_profile.id
    source    = "static"
  }

  relation_l3extrs_redistribute_pol {
    target_dn = data.aci_route_control_profile.shared_route_control_profile.id
    source    = "direct"
  }
  default_route_leak_policy {
    always     = "no"
    annotation = "orchestrator:terraform"
    criteria   = "only"
    scope      = ["l3-out", "ctx"]
  }

}
