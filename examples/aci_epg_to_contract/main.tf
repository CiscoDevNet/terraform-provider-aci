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

resource "aci_tenant" "tenant_for_epg_contract" {
  name        = "tenant_for_epg_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "tenant_contract" {
  tenant_dn   = aci_tenant.tenant_for_epg_contract.id
  name        = "tenant_contract"
}

resource "aci_application_profile" "tenant_ap" {
  tenant_dn   = aci_tenant.tenant_for_epg_contract.id
  name        = "AP"
}

resource "aci_application_epg" "application_epg" {
  application_profile_dn = aci_application_profile.tenant_ap.id
  name                   = "app_epg"
}

resource "aci_epg_to_contract" "example_provider" {
  application_epg_dn = aci_application_epg.application_epg.id
  contract_dn        = aci_contract.tenant_contract.id
  contract_type      = "provider"
}

resource "aci_epg_to_contract" "example_consumer" {
  application_epg_dn = aci_application_epg.application_epg.id
  contract_dn        = aci_contract.tenant_contract.id
  contract_type      = "consumer"
  prio               = "level1"
}

data "aci_epg_to_contract" "example_consumer_2" {
  application_epg_dn = aci_application_epg.application_epg.id
  contract_name      = "consumer_contract"
  contract_type      = "consumer"
}

