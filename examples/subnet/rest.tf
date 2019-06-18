resource "aci_rest" "rest_l3_ext_out" {
  path       = "api/node/mo/${aci_tenant.tenant_for_subnet.id}/out-testext.json"
  class_name = "l3extOut"

  content = {
    "name" = "testext"
  }
}

resource "aci_rest" "rest_rt_ctrl_profile" {
  path       = "api/node/mo/${aci_rest.rest_l3_ext_out.id}/prof-testprof.json"
  class_name = "rtctrlProfile"

  content = {
    "name" = "testprof"
  }
}


resource "aci_rest" "rest_nd_pfx_pol" {
  path       = "api/node/mo/${aci_tenant.tenant_for_subnet.id}/ndpfxpol-testpol.json"
  class_name = "ndPfxPol"

  content = {
    "name" = "testpol"
  }
}
