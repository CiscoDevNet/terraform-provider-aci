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
    name = "ap1"
}

data "aci_application_epg" "epg" {
  application_profile_dn = "${data.aci_application_profile.ap.id}"
  name = "checkepg"
}

data "aci_client_end_point" "example" {
  application_epg_dn  = "${data.aci_application_epg.epg.id}"
  name                = "example"
}