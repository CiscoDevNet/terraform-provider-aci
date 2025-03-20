
resource "aci_application_profile" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
