resource "aci_out_of_band_contract" "example" {
  name = "test_out_of_band_contract"
  annotations = [
    {
      key = "test_annotation"
    },
  ]
}
