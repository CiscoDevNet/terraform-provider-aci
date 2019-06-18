resource "aci_rest" "rest_mon_epg_pol" {
  path       = "api/node/mo/${aci_tenant.test_tenant.id}/monepg-testpol.json"
  class_name = "monEPGPol"

  content = {
    "name" = "testpol"
  }
}
