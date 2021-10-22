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

resource "aci_any" "example_vzany" {
  vrf_dn       = aci_vrf.example.id
  description  = "vzAny Description"
  annotation   = "tag_any"
  match_t      = "AtleastOne"
  name_alias   = "alias_any"
  pref_gr_memb = "disabled"
}