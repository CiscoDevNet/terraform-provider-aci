data "aci_rest" "fvTenant" {
  dn = "uni/tn-EXAMPLE_TENANT_DATA"
}

resource "aci_rest" "fvTenant" {
  dn         = "uni/tn-EXAMPLE_TENANT"
  class_name = "fvTenant"
  content = {
    name  = "EXAMPLE_TENANT"
    descr = "Example description"
  }
}

resource "aci_rest" "mgmtConnectivityPrefs" {
  dn         = "uni/fabric/connectivityPrefs"
  class_name = "mgmtConnectivityPrefs"
  content = {
    interfacePref = "ooband"
  }
}

resource "aci_rest" "fvTenant" {
  dn         = "uni/tn-EXAMPLE_TENANT"
  class_name = "fvTenant"
  content = {
    name = "EXAMPLE_TENANT"
  }

  child {
    rn         = "ctx-VRF1"
    class_name = "fvCtx"
    content = {
      name = "VRF1"
    }
  }
}