terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_vrf" {
  name = "terraform_vrf"
}

data "aci_tenant" "common" {
  name = "common"
}

resource "aci_vrf" "vrf1" {
  tenant_dn          = aci_tenant.terraform_vrf.id
  bd_enforced_enable = "no"
  knw_mcast_act      = "permit"
  name               = "vrf1"
  pc_enf_dir         = "ingress"
  pc_enf_pref        = "enforced"
  relation_fv_rs_ctx_to_bgp_ctx_af_pol {
    af                     = "ipv4-ucast"
    tn_bgp_ctx_af_pol_name = aci_bgp_address_family_context.one.id
  }
}

data "aci_vrf" "common_default" {
  tenant_dn          = data.aci_tenant.common.id
  name               = "default"
}

resource "aci_bgp_address_family_context" "one" {
  tenant_dn     = aci_tenant.terraform_vrf.id
  name          = "one"
  ctrl          = "host-rt-leak"
  e_dist        = "25"
  i_dist        = "198"
  local_dist    = "100"
  max_ecmp      = "18"
  max_ecmp_ibgp = "25"
}

resource "aci_bgp_address_family_context" "two" {
  tenant_dn     = data.aci_tenant.common.id
  name          = "two"
  ctrl          = "host-rt-leak"
  e_dist        = "25"
  i_dist        = "198"
  local_dist    = "100"
  max_ecmp      = "18"
  max_ecmp_ibgp = "25"
}


resource "aci_vrf_to_bgp_address_family_context" "example_v4" {
  vrf_dn  = data.aci_vrf.common_default.id
  bgp_address_family_context_dn = aci_bgp_address_family_context.two.id
  address_family  = "ipv4-ucast"
}

resource "aci_vrf_to_bgp_address_family_context" "example_v6" {
  vrf_dn  = data.aci_vrf.common_default.id
  bgp_address_family_context_dn = aci_bgp_address_family_context.two.id
  address_family  = "ipv6-ucast"
}

data "aci_vrf_to_bgp_address_family_context" "example_v4" {
  vrf_dn  = aci_vrf_to_bgp_address_family_context.example_v4.vrf_dn
  bgp_address_family_context_dn  = aci_bgp_address_family_context.two.id
  address_family  = "ipv4-ucast"
}

data "aci_vrf_to_bgp_address_family_context" "example_v6" {
  vrf_dn  = aci_vrf_to_bgp_address_family_context.example_v6.vrf_dn
  bgp_address_family_context_dn  = aci_bgp_address_family_context.two.id
  address_family  = "ipv6-ucast"
}
