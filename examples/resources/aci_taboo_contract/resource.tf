
resource "aci_taboo_contract" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
