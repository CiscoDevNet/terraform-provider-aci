resource "aci_tenant" "tenant_for_epg" {
  name        = "tenant_for_epg"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_application_profile" "app_profile_for_epg" {
  tenant_dn   = aci_tenant.tenant_for_epg.id
  name        = "ap_for_epg"
  description = "This app profile is created by terraform ACI providers"
}

resource "aci_application_epg" "inherit_epg" {
  application_profile_dn = aci_application_profile.app_profile_for_epg.id
  name                   = "inherit_epg"
  description            = "epg to create relation sec_inherited"

}

// Creation of Monitoring policy
resource "aci_monitoring_policy" "rest_mon_epg_pol" {
  tenant_dn = aci_tenant.tenant_for_epg.id
  name      = "testpol"
}

resource "aci_application_epg" "demoepg" {
  application_profile_dn       = aci_application_profile.app_profile_for_epg.id
  name                         = "tf_test_epg"
  description                  = "This epg is created by terraform ACI providers"
  flood_on_encap               = "disabled"
  fwd_ctrl                     = "none"
  is_attr_based_epg            = "no"
  match_t                      = "None"
  pc_enf_pref                  = "unenforced"
  pref_gr_memb                 = "exclude"
  prio                         = "unspecified"
  relation_fv_rs_bd            = aci_bridge_domain.bd_for_rel.id            # Relation to Bridge Domain
  relation_fv_rs_cust_qos_pol  = aci_rest.rest_qos_custom_pol.id            # Relation to Custom Quality of Service - QoS traffic policy
  relation_fv_rs_prov          = [aci_contract.rs_prov_contract.id]         # Relation to Provided Contract
  relation_fv_rs_cons_if       = [aci_imported_contract.rest_vz_cons_if.id] # Relation to Imported Contract
  relation_fv_rs_sec_inherited = [aci_application_epg.inherit_epg.id]       # Relation to inherit security configuration from another EPG
  relation_fv_rs_dpp_pol       = aci_rest.rest_qos_dpp_pol.id               # Relation to Data Plane Policing
  relation_fv_rs_cons          = [aci_contract.rs_cons_contract.id]         # Relation to Consumed Contract
  relation_fv_rs_trust_ctrl    = aci_rest.rest_trust_ctrl_pol.id            # Relation to First Hop Security trust control
  relation_fv_rs_prot_by       = [aci_taboo_contract.rest_taboo_con.id]     # Relation to vzTaboo Taboo Contract
  relation_fv_rs_aepg_mon_pol  = aci_monitoring_policy.rest_mon_epg_pol.id  # Relation to Monitoring policy
  relation_fv_rs_intra_epg     = [aci_contract.intra_epg_contract.id]       # Relation to Intra EPG Contract
}
