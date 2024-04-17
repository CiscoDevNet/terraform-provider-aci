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

resource "aci_rest_managed" "aaaPreLoginBanner" {
  dn          = "uni/userext/preloginbanner"
  class_name  = "aaaPreLoginBanner"
  escape_html = false
  content = {
    message = <<-EOT
<<< WARNING >>>  THE PROGRAMS AND DATA HELD ON THIS SYSTEM
ARE THE PROPERTY OF, OR LICENCED BY, XXXXXXXX. IF THE COMPANY HAS NOT
AUTHORISED YOUR ACCESS TO THIS SYSTEM YOU WILL COMMIT A CRIMINAL OFFENCE IF
YOU DO NOT IMMEDIATELY DISCONNECT. UNAUTHORISED ACCESS IS STRICTLY FORBIDDEN
AND IS A DISCIPLINARY OFFENCE.
      EOT
  }
}