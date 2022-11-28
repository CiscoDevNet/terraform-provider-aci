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
  name        = "tenant"
  description = "tenant created while acceptance testing"

}
resource "aci_vrf" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "demo_vrf"
}

resource "aci_contract" "web_contract" {
  tenant_dn   = aci_tenant.example.id
  name        = "web_contract"
}

resource "aci_contract" "db_contract" {
  tenant_dn   = aci_tenant.example.id
  name        = "db_contract"
}

resource "aci_imported_contract" "contract_interface" {
  tenant_dn         = aci_tenant.example.id
  name              = "exported_contract_from_tenant_epg"
}

resource "aci_any" "example_vzany" {
  vrf_dn       = aci_vrf.example.id
  description  = "vzAny Description"
  annotation   = "tag_any"
  match_t      = "AtleastOne"
  name_alias   = "alias_any"
  pref_gr_memb = "disabled"
  relation_vz_rs_any_to_cons = [aci_contract.db_contract.id, aci_contract.web_contract.id]
  relation_vz_rs_any_to_cons_if = [aci_imported_contract.contract_interface.id]
  relation_vz_rs_any_to_prov = [aci_contract.web_contract.id]
}