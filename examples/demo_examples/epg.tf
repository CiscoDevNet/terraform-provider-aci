resource "aci_tenant" "tenant_for_epg" {
  name        = "tenant_for_epg"
  description = "This tenant is created by terraform ACI provider"
}
resource "aci_app_profile" "app_profile_for_epg" {
  tenant_dn   = "${aci_tenant.tenant_for_epg.id}"
  name        = "ap_for_epg"
  description = "This app profile is created by terraform ACI providers"
}
resource "aci_epg" "demoepg" {
  application_profile_dn = "${aci_app_profile.app_profile_for_epg.id}"
  name                   = "tf_test_epg"
  description            = "This epg is created by terraform ACI providers"
}

