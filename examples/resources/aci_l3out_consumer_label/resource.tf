
resource "aci_l3out_consumer_label" "example" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
  annotations = [
    {
      key   = "test_key"
      value = "test_value"
    },
  ]
}

