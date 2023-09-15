
resource "aci_annotation" "example" {
  parent_dn = aci_application_epg.example.id
  key       = "test_annotation"
}

