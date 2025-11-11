
resource "aci_leaf_interface_profile" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
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
