
resource "aci_l3out_consumer_label" "example" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

