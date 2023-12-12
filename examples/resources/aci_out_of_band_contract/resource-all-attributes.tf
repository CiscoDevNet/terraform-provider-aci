
resource "aci_out_of_band_contract" "full_example" {
  annotation  = "annotation"
  description = "description"
  intent      = "estimate_add"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  priority    = "level1"
  scope       = "application-profile"
  target_dscp = "AF11"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
