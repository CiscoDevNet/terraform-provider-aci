resource "aci_rest" "rest_l3_ext_out" {
  path       = "api/node/mo/${aci_tenant.tenant_for_bridge_domain.id}/out-testext.json"
  class_name = "l3extOut"

  content = {
    "name" = "testext"
  }
}

resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenant_for_bridge_domain.id
  name                       = "testprof"
}

resource "aci_rest" "rest_dhcp_RelayP" {
  path       = "api/node/mo/${aci_tenant.tenant_for_bridge_domain.id}/relayp-testrelay.json"
  class_name = "dhcpRelayP"

  content = {
    "name" = "testrelay"
  }
}

resource "aci_rest" "rest_mon_epg_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_bridge_domain.id}/monepg-testpol.json"
  class_name = "monEPGPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_fhs_bd_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_bridge_domain.id}/bdpol-testpolbd.json"
  class_name = "fhsBDPol"

  content = {
    "name" = "testpolbd"
  }
}
resource "aci_rest" "rest_net_flow_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_bridge_domain.id}/monitorpol-testpolflow.json"
  class_name = "netflowMonitorPol"

  content = {
    "name" = "testpolflow"
  }
}
