data "aci_rest_managed" "fvTenant" {
  dn = "uni/tn-infra"
}

resource "aci_rest_managed" "fvTenant1" {
  dn         = "uni/tn-EXAMPLE_TENANT"
  class_name = "fvTenant"
  content = {
    name  = "EXAMPLE_TENANT"
    descr = "Example description"
  }
}

resource "aci_rest_managed" "mgmtConnectivityPrefs" {
  dn         = "uni/fabric/connectivityPrefs"
  class_name = "mgmtConnectivityPrefs"
  annotation = "orchestrator:example"
  content = {
    interfacePref = "ooband"
  }
}

resource "aci_rest_managed" "fvTenant2" {
  dn         = "uni/tn-EXAMPLE_TENANT2"
  class_name = "fvTenant"
  content = {
    name       = "EXAMPLE_TENANT2"
    annotation = "orchestrator:class"
  }

  child {
    rn         = "ctx-VRF1"
    class_name = "fvCtx"
    content = {
      name       = "VRF1"
      annotation = "orchestrator:child"
    }
  }
}