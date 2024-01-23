
data "aci_certificate_authority" "example_public_key_management" {
  parent_dn = aci_public_key_management.example.id
  name      = "test_name"
}
