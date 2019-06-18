resource "aci_tenant" "tenant_for_subject" {
  name        = "tenant_for_subject"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_contract" "contract_for_subject" {
  tenant_dn   = "${aci_tenant.tenant_for_subject.id}"
  name        = "contract_for_subject"
  description = "This contract is created by terraform ACI provider"
  scope       = "context"
  target_dscp = "VA"
}

resource "aci_contract_subject" "demosubject" {
  contract_dn                   = "${aci_contract.contract_for_subject.id}"
  name                          = "test_tf_subject"
  description                   = "This subject is created by terraform ACI provider"
  cons_match_t                  = "None"
  prio                          = "unspecified"
  prov_match_t                  = "None"
  rev_flt_ports                 = "no"
  target_dscp                   = "unspecified"
  relation_vz_rs_subj_graph_att = "${aci_rest.rest_abs_graph.id}"  # Relation to vnsAbsGraph class. Cardinality - N_TO_ONE.
  relation_vz_rs_subj_filt_att  = ["${aci_filter.subj_filter.id}"] # Relation to vzFilter class. Cardinality - N_TO_M.
}
