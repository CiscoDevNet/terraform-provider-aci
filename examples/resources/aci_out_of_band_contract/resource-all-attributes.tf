
resource "aci_out_of_band_contract" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  intent      = "estimate_add"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  priority    = "level1"
  scope       = "application-profile"
  target_dscp = "AF11"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
