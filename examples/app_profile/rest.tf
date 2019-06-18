resource "aci_rest" "rest_mon_epg_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_ap.id}/monepg-testpol.json"
  class_name = "monEPGPol"

  content = {
    "name" = "testpol"
  }
}
