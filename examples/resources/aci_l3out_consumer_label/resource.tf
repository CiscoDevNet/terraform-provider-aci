
resource "aci_l3out_consumer_label" "example" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_l3out_consumer_label"
  annotations = [
    {
      key = "test_annotation"
    },
  ]
}

