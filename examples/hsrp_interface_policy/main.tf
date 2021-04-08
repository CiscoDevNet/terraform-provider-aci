provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_hsrp_interface_policy" "example" {
  tenant_dn    = "${aci_tenant.tenentcheck.id}"
  name         = "one"
  annotation   = "example"
  description  = "from terraform"
  ctrl         = "bia"
  delay        = "10"
  name_alias   = "example"
  reload_delay = "10"
}