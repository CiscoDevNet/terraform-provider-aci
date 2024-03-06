
resource "aci_tag" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  key       = "test_key"
  value     = "test_value"
}

resource "aci_tag" "example_endpoint_tag_ip" {
  parent_dn = aci_endpoint_tag_ip.example.id
  key       = "test_key"
  value     = "test_value"
}
