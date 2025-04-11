
resource "aci_eigrp_address_family_context" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
