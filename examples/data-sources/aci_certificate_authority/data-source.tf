data "aci_certificate_authority" "example" {
  name = "test_name"
}

data "aci_certificate_authority" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}

