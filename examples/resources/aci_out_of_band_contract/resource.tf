resource "aci_out_of_band_contract" "example" {
  name = "test_name"
  annotations = [
    {
      key   = "test_key"
      value = "test_value"
    },
  ]
}
