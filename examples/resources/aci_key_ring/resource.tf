
resource "aci_key_ring" "example" {
  name = "test_name"
}

resource "aci_key_ring" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
