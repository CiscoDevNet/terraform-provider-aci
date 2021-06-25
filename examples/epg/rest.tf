resource "aci_rest" "rest_qos_custom_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/qoscustom-testpol.json"
  class_name = "qosCustomPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_qos_dpp_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/qosdpppol-testqospol.json"
  class_name = "qosDppPol"

  content = {
    "name" = "testqospol"
  }
}
resource "aci_rest" "rest_trust_ctrl_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/trustctrlpol-testtrustpol.json"
  class_name = "fhsTrustCtrlPol"

  content = {
    "name" = "testtrustpol"
  }
}
