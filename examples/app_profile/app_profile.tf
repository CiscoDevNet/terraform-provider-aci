resource "aci_tenant" "tenant_for_ap" {
  name        = "tenant_for_ap"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_application_profile" "demo_app_profile" {
  tenant_dn                 = aci_tenant.tenant_for_ap.id
  name                      = "test_tf_ap"
  description               = "This app profile is created by terraform ACI provider"
  prio                      = "unspecified"
  annotation                = "test_ap_anot"
  relation_fv_rs_ap_mon_pol = aci_rest.rest_mon_epg_pol.id # Relation to class monEPGPol. Cardinality - N_TO_ONE
}
