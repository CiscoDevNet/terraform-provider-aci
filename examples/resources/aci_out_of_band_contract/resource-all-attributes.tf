
resource "aci_out_of_band_contract" "full_example" {
  annotation  = "annotation"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
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
