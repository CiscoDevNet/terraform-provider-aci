resource "aci_tenant" "tenant_for_filter" {
  name        = "tenant_for_filter"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_filter" "demo_filter" {
  tenant_dn                      = aci_tenant.tenant_for_filter.id
  name                           = "test_tf_filter"
  description                    = "This filter is created by terraform ACI provider."
  relation_vz_rs_filt_graph_att  = "testinterm"  # Relation to vnsInTerm class.
}
