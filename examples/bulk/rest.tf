resource "aci_rest" "rest_l3_ext_out" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/out-testext.json"
  class_name = "l3extOut"

  content = {
    "name" = "testext"
  }
}

resource "aci_rest" "rest_l3_ext_out2" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/out-testext2.json"
  class_name = "l3extOut"

  content = {
    "name" = "testext2"
  }
}

resource "aci_rest" "rest_l3_ext_out3" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/out-testext3.json"
  class_name = "l3extOut"

  content = {
    "name" = "testext3"
  }
}

resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenant_for_benchmark.id
  name                       = "testprof"
}

resource "aci_rest" "rest_dhcp_RelayP" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/relayp-testrelay.json"
  class_name = "dhcpRelayP"

  content = {
    "name" = "testrelay"
  }
}

resource "aci_rest" "rest_mon_epg_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/monepg-testpol.json"
  class_name = "monEPGPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_fhs_bd_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/bdpol-testpolbd.json"
  class_name = "fhsBDPol"

  content = {
    "name" = "testpolbd"
  }
}

resource "aci_rest" "rest_net_flow_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/monitorpol-testpolflow.json"
  class_name = "netflowMonitorPol"

  content = {
    "name" = "testpolflow"
  }
}

resource "aci_rest" "rest_net_flow_pol2" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/monitorpol-testpolflow2.json"
  class_name = "netflowMonitorPol"

  content = {
    "name" = "testpolflow2"
  }
}

resource "aci_rest" "rest_net_flow_pol3" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/monitorpol-testpolflow3.json"
  class_name = "netflowMonitorPol"

  content = {
    "name" = "testpolflow3"
  }
}

resource "aci_rest" "rest_qos_custom_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/qoscustom-testpol.json"
  class_name = "qosCustomPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_qos_dpp_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/qosdpppol-testqospol.json"
  class_name = "qosDppPol"

  content = {
    "name" = "testqospol"
  }
}

resource "aci_rest" "rest_trust_ctrl_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/trustctrlpol-testtrustpol.json"
  class_name = "fhsTrustCtrlPol"

  content = {
    "name" = "testtrustpol"
  }
}

###

resource "aci_rest" "rest_mld_snoop_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/mldsnoopPol-testmldsnooppol.json"
  class_name = "mldSnoopPol"

  content = {
    "name" = "testmldsnooppol"
  }
}

resource "aci_rest" "rest_igmp_snoop_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/snPol-testigmpsnooppol.json"
  class_name = "igmpSnoopPol"

  content = {
    "name" = "testigmpsnooppol"
  }
}

resource "aci_rest" "rest_endpoint_ret_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/epRPol-testendpretpol.json"
  class_name = "fvEpRetPol"

  content = {
    "name" = "testendpretpol"
  }
}

resource "aci_rest" "rest_nd_if_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_benchmark.id}/ndifpol-testndifpol.json"
  class_name = "ndIfPol"

  content = {
    "name" = "testndifpol"
  }
}
/*
resource "aci_rest" "rest_fab_node_pol" {
  path       = "api/node/mo/topology/pod-101.json"
  class_name = "fabricNode"

  content = {
    "name" = "testfabnodepol"
  }
}
*/
