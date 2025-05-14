
resource "aci_external_epg" "example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
}
