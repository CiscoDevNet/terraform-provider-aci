
data "aci_tag" "example_tenant" {
  parent_dn = aci_tenant.example.id
  key       = "test_key"
}

data "aci_tag" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  key       = "test_key"
}
