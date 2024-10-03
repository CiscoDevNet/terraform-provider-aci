
resource "aci_key_ring" "example" {
  name = "test_name"
}

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_key_ring" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
