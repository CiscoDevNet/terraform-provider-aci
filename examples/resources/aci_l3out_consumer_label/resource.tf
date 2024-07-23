
resource "aci_l3out_consumer_label" "example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
}
