resource "aci_tenant" "test_tenant" {
  name        = "tf_test_rel_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_ospf_interface_policy" "test_ospf" {
  tenant_dn = aci_tenant.test_tenant.id
  name      = "tf_ospf"
}
