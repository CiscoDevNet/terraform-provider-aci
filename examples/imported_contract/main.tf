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

resource "aci_tenant" "tenant_epg" {
  name        = "tenant_epg"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_tenant" "tenant_contract" {
  name        = "tenant_contract"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "tenant_contract" {
  tenant_dn = aci_tenant.tenant_contract.id
  name      = "tenant_contract"
}

resource "aci_imported_contract" "contract_interface" {
  tenant_dn = aci_tenant.tenant_epg.id
  name      = "exported_contract_from_tenant_epg"
}

