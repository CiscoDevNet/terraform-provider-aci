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

resource "aci_tenant" "terraform_ten" {
  name = "terraform_ten"
}

resource "aci_vrf" "vrf1" {
  tenant_dn          = aci_tenant.terraform_ten.id
  bd_enforced_enable = "no"
  knw_mcast_act      = "permit"
  name               = var.vrf_name
  pc_enf_dir         = "ingress"
  pc_enf_pref        = "enforced"
  relation_fv_rs_ctx_to_bgp_ctx_af_pol {
    af                     = "ipv4-ucast"
    tn_bgp_ctx_af_pol_name = "test_bgp"
  }
}
