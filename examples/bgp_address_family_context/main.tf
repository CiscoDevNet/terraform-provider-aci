provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_bgp_address_family_context" "example" {
  tenant_dn     = "${aci_tenant.tenentcheck.id}"
  name          = "one"
  description   = "from terraform"
  annotation    = "example"
  ctrl          = "host-rt-leak"
  e_dist        = "25"
  i_dist        = "198"
  local_dist    = "100"
  max_ecmp      = "18"
  max_ecmp_ibgp = "25"
  name_alias    = "example"
}