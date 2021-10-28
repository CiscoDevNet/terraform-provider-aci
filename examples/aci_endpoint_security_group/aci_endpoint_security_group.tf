resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_ap"
}

# create ESG
resource "aci_endpoint_security_group" "terraform_esg" {
  application_profile_dn      = aci_application_profile.terraform_ap.id
  name                        = "tf_esg"
  relation_fv_rs_scope        = aci_vrf.vrf.id
  relation_fv_rs_prov         = [aci_contract.rs_prov_contract.id]
  relation_fv_rs_cons         = [aci_contract.rs_cons_contract.id]
  relation_fv_rs_intra_epg    = [aci_contract.intra_epg_contract.id]
  relation_fv_rs_cons_if      = [aci_contract.exported_contract.id]
  relation_fv_rs_cust_qos_pol = aci_rest.rest_qos_custom_pol.id
  relation_fv_rs_prot_by      = [aci_rest.rest_taboo_con.id]
}

# create another ESG_2, inheriting from ESG
resource "aci_endpoint_security_group" "terraform_inherit_esg" {
  application_profile_dn       = aci_application_profile.terraform_ap.id
  name                         = "tf_inherit_esg"
  description                  = "create relation sec_inherited"
  match_t                      = "None"
  pc_enf_pref                  = "unenforced"
  pref_gr_memb                 = "exclude"
  relation_fv_rs_sec_inherited = [aci_endpoint_security_group.terraform_esg.id]
}

# query an existing ESG. In this case, creating an ESG named 'test' before terraform apply
data "aci_endpoint_security_group" "query_esg" {
  application_profile_dn = "uni/tn-test_esg/ap-esg_ap"
  name                   = "test"
}

output "data_source_esg" {
  description = "ESG queried by data source"
  value       = data.aci_endpoint_security_group.query_esg.id
}