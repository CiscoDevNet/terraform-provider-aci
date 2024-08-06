
data "aci_application_epg" "example_application_profile" {
  parent_dn = aci_application_profile.example.id
  name      = "test_name"
}
