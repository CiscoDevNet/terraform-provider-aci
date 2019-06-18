resource "aci_rest" "rest_qos_custom_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/qoscustom-testpol.json"
  class_name = "qosCustomPol"

  content = {
    "name" = "testpol"
  }
}


resource "aci_rest" "rest_infra_domp" {
  path       = "api/node/mo/uni/fc-test.json"
  class_name = "fcDomP"

  content = {
    "name" = "test"
  }
}
resource "aci_rest" "rest_mon_epg_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/monepg-testpol.json"
  class_name = "monEPGPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_vz_cons_if" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/cif-testcontract.json"
  class_name = "vzCPIf"

  content = {
    "name" = "testcontract"
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
resource "aci_rest" "rest_taboo_con" {
  path       = "api/node/mo/${aci_tenant.tenant_for_epg.id}/taboo-testcon.json"
  class_name = "vzTaboo"

  content = {
    "name" = "testcon"
  }
}
