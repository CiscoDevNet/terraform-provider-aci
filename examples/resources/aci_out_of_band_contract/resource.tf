resource "aci_out_of_band_contract" "example" {
  name = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
