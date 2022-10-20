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

data "aci_tenant" "common" {
  name = "common"
}

data "aci_vrf" "default" {
  tenant_dn = data.aci_tenant.common.id
  name      = "default"
}

resource "aci_tenant" "terraform_vrf" {
  name = "terraform_vrf"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.terraform_vrf.id
  name      = "vrf1"
}

resource "aci_vrf" "vrf2" {
  tenant_dn = aci_tenant.terraform_vrf.id
  name      = "vrf2"
}

resource "aci_vrf_leak_epg_bd_subnet" "vrf_leak_epg_bd_subnet" {
  vrf_dn                    = aci_vrf.vrf1.id
  ip                        = "1.1.20.2/24"
  allow_l3out_advertisement = true
  leak_to {
    vrf_dn                    = data.aci_vrf.default.id
    allow_l3out_advertisement = "inherit"
  }
  leak_to {
    vrf_dn                    = aci_vrf.vrf2.id
    allow_l3out_advertisement = true # Must be set as true for Cloud APIC
  }
}
