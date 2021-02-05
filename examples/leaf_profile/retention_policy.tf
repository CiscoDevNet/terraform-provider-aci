resource "aci_end_point_retention_policy" "test_ret_policy" {
    tenant_dn = aci_tenant.test_tenant.id
    name = "tf_test"
  
}
