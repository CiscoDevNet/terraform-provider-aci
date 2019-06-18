resource "aci_app_profile" "demo_app_profile" {
  tenant_dn   = "${aci_tenant.demotenant.id}"
  name        = "test_tf_ap"
  description = "This app profile is created by terraform ACI provider"
}
