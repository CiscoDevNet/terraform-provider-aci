
data "aci_endpoint_security_group" "example_application_profile" {
  parent_dn = aci_application_profile.example.id
  name      = "test_name"
}
