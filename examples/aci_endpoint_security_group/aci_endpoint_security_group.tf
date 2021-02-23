resource "aci_tenant" "terraform_tenant" {
    name        = "tf_tenant"
    description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
    tenant_dn  = aci_tenant.terraform_tenant.id
    name       = "tf_ap"
}

# # create ESG
# resource "aci_endpoint_security_group" "terraform_esg" {
#     application_profile_dn  = aci_application_profile.terraform_ap.id
#     name                    = "tf_esg"
# }

# # create another ESG_2, inheriting from ESG
# resource "aci_endpoint_security_group" "terraform_esg_2" {
#     application_profile_dn       = aci_application_profile.terraform_ap.id
#     name                         = "tf_esg_2"
#     description                  = "create relation sec_inherited"
#     relation_fv_rs_sec_inherited = [aci_endpoint_security_group.terraform_esg.id]
# }

# # ESG with provider contract
# resource "aci_endpoint_security_group" "terraform_esg_3" {
#     application_profile_dn  = aci_application_profile.terraform_ap.id
#     name                    = "tf_esg"
#     relation_fv_rs_prov     = [aci_contract.rs_prov_contract.id]
# }

# # ESG with consumer contract
# resource "aci_endpoint_security_group" "terraform_esg_4" {
#     application_profile_dn  = aci_application_profile.terraform_ap.id
#     name                    = "tf_esg"
#     relation_fv_rs_cons     = [aci_contract.rs_cons_contract.id]
# }

# # ESG with intra epg contract
# resource "aci_endpoint_security_group" "terraform_esg_5" {
#     application_profile_dn       = aci_application_profile.terraform_ap.id
#     name                         = "tf_esg"
#     relation_fv_rs_intra_epg     = [aci_contract.intra_epg_contract.id]
# }

# # ESG with contract interface
# resource "aci_endpoint_security_group" "terraform_esg_6" {
#     application_profile_dn       = aci_application_profile.terraform_ap.id
#     name                         = "tf_esg"
#     relation_fv_rs_cons_if       = [aci_rest.rest_cons_if.id]
# }

