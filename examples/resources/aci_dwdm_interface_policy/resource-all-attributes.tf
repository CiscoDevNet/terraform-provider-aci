
resource "aci_dwdm_interface_policy" "full_example" {
  annotation      = "annotation"
  description     = "description_1"
  channel_nummber = "Channel1"
  name            = "test_name"
  name_alias      = "name_alias_1"
  owner_key       = "owner_key_1"
  owner_tag       = "owner_tag_1"
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
