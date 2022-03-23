terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://10.23.248.120"
  insecure = true
}

resource "aci_tenant" "tenant_epg" {
  name        = "tenant_epg"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_tenant" "tenant_contract" {
  name        = "tenant_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "tenant_contract" {
  tenant_dn   = aci_tenant.tenant_contract.id
  name        = "tenant_contract"
}
resource "aci_imported_contract" "contract_interface" {
  tenant_dn         = aci_tenant.tenant_epg.id
  name              = "exported_contract_from_tenant_epg"
}

resource "aci_application_profile" "tenant_ap" {
  tenant_dn   = aci_tenant.tenant_epg.id
  name        = "AP"
}

resource "aci_application_epg" "application_epg" {
  application_profile_dn = aci_application_profile.tenant_ap.id
  name                   = "app_epg"
}

resource "aci_epg_to_contract_interface" "epg_contract_interface" {
  application_epg_dn = aci_application_epg.application_epg.id
  contract_interface_dn = aci_imported_contract.contract_interface.id
}

