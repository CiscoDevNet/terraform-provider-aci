
resource "aci_first_hop_security_policy" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
