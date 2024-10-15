
resource "aci_neighbor_discovery_interface_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
