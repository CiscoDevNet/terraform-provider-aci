provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "tenentcheck" {
  name       = "test"
  annotation = "atag"
  name_alias = "alias_tenant"
}

resource "aci_bgp_peer_prefix" "example" {
  tenant_dn    = "${aci_tenant.tenentcheck.id}"
  name         = "one"
  description  = "from terraform"
  action       = "shut"
  annotation   = "example"
  max_pfx      = "200"
  name_alias   = "example"
  restart_time = "200"
  thresh       = "85"
}