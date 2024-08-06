
resource "aci_trust_control_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
