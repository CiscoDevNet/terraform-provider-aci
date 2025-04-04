
data "aci_certificate_authority" "example" {
  name = "test_name"
}

// This example is only applicable to Cisco Cloud Network Controller
data "aci_certificate_authority" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
