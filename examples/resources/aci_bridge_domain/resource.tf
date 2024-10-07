
resource "aci_bridge_domain" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
