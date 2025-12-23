
resource "aci_fibre_channel_interface_policy" "full_example" {
  annotation            = "annotation"
  auto_max_speed        = "16G"
  description           = "description_1"
  fill_pattern          = "ARBFF"
  name                  = "test_name"
  name_alias            = "name_alias_1"
  owner_key             = "owner_key_1"
  owner_tag             = "owner_tag_1"
  port_mode             = "f"
  receive_buffer_credit = "16"
  speed                 = "16G"
  trunk_mode            = "auto"
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
