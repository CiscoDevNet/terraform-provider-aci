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

data "aci_tenant" "tenant_fetch" {
  name = "data_source_test"
}

data "aci_application_profile" "ap_fetch" {
  tenant_dn = data.aci_tenant.tenant_fetch.id
  name      = "data_source_ap"
}

data "aci_application_epg" "epg_fetch" {
  application_profile_dn = data.aci_application_profile.ap_fetch.id
  name                   = "data_source_epg"
}

data "aci_vrf" "vrf_fetch" {
  tenant_dn = data.aci_tenant.tenant_fetch.id
  name      = "data_source_vrf"
}

data "aci_bridge_domain" "bd_fetch" {
  tenant_dn = data.aci_tenant.tenant_fetch.id
  name      = "data_source_bd"
}

data "aci_subnet" "subnet_fetch" {
  bridge_domain_dn = data.aci_bridge_domain.bd_fetch.id
  ip               = "10.0.1.1/24"
  name             = "10.0.1.1/24"
}

data "aci_contract" "contract_fecth" {
  tenant_dn = data.aci_tenant.tenant_fetch.id
  name      = "data_source_contract"
}

data "aci_contract_subject" "subject_fetch" {
  contract_dn = data.aci_contract.contract_fecth.id
  name        = "data_source_subject"
}

data "aci_filter" "fiter_fetch" {
  tenant_dn = data.aci_tenant.tenant_fetch.id
  name      = "data_source_filter"
}

data "aci_filter_entry" "fetch_entry" {
  filter_dn = data.aci_filter.fiter_fetch.id
  name      = "data_source_entry"
}

data "aci_vmm_domain" "fetch_domain" {
  provider_profile_dn = "uni/vmmp-VMware"
  name                = "test"
}
