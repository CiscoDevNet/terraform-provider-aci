
resource "aci_link_level_flow_control_interface_policy" "full_example" {
  annotation   = "annotation"
  description  = "description_1"
  receive_mode = "on"
  send_mode    = "on"
  name         = "test_name"
  name_alias   = "name_alias_1"
  owner_key    = "owner_key_1"
  owner_tag    = "owner_tag_1"
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
