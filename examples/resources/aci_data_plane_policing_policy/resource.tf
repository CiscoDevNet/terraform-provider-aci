
resource "aci_data_plane_policing_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
