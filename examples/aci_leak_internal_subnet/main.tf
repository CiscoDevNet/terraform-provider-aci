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
}


resource "aci_vrf" "vrf2" {
  tenant_dn          = aci_tenant.terraform_vrf.id
  bd_enforced_enable = "no"
  knw_mcast_act      = "permit"
  name               = "vrf2"
  pc_enf_dir         = "ingress"
  pc_enf_pref        = "enforced"
}


resource "aci_leak_internal_subnet" "internal_subnet" {
  vrf_dn = aci_vrf.vrf1.id
  ip     = "1.1.20.2/24"
  leak_to {
    destination_vrf_name    = "default"
    destination_tenant_name = "common"
  }
  leak_to {
    destination_vrf_name    = aci_vrf.vrf2.name
    destination_tenant_name = aci_tenant.terraform_vrf.name
    destination_vrf_scope   = "private"
  }
}
