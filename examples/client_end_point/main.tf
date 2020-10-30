provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

data "aci_tenant" "tenant" {
  name = "tenant"
}

data "aci_application_profile" "ap" {
  tenant_dn = "${data.aci_tenant.tenant.id}"
  name      = "ap1"
}

data "aci_application_epg" "epg" {
  application_profile_dn = "${data.aci_application_profile.ap.id}"
  name                   = "checkepg"
}

data "aci_client_end_point" "check" {
  application_epg_dn = "${aci_application_epg.epg.id}"
  name               = "25:56:68:78:98:74"
  mac                = "25:56:68:78:98:74"
  ip                 = "1.2.3.4"
  vlan               = "5"
}
