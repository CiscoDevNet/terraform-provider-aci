
resource "aci_pim_route_map_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
