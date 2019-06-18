resource "aci_tenant" "tenant_for_filter" {
  name        = "tenant_for_filter"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_filter" "demofilter" {
  tenant_dn                      = "${aci_tenant.tenant_for_filter.id}"
  name                           = "test_tf_filter"
  description                    = "This filter is created by terraform ACI provider."
  relation_vz_rs_filt_graph_att  = "testinterm"                  # Relation to  vnsInTerm class. Cardinality - N_TO_ONE.
  relation_vz_rs_fwd_r_flt_p_att = "${aci_filter.fwd_filter.id}" # Relation to  vzAFilterableUnit class. Cardinality - N_TO_ONE.
  relation_vz_rs_rev_r_flt_p_att = "${aci_filter.rev_filter.id}" # Relation to  vzAFilterableUnit class. Cardinality - N_TO_ONE.
}

resource "aci_filter" "fwd_filter" {
  tenant_dn   = "${aci_tenant.tenant_for_filter.id}"
  name        = "test_tf_filter"
  description = "This filter is created by terraform ACI provider."
}
resource "aci_filter" "rev_filter" {
  tenant_dn   = "${aci_tenant.tenant_for_filter.id}"
  name        = "test_tf_filter"
  description = "This filter is created by terraform ACI provider."
}
