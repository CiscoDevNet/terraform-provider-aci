provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_bgp_timers" "example1" {
  tenant_dn    = "${aci_tenant.tenentcheck.id}"
  name         = "one"
  annotation   = "example"
  gr_ctrl      = "helper"
  hold_intvl   = "189"
  ka_intvl     = "65"
  max_as_limit = "70"
  name_alias   = "aliasing"
  stale_intvl  = "15"
}