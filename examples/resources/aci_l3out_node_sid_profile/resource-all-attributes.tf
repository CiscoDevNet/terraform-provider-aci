
resource "aci_l3out_node_sid_profile" "full_example_l3out_loopback_interface_profile" {
  parent_dn        = aci_l3out_loopback_interface_profile.example.id
  annotation       = "annotation"
  description      = "description"
  loopback_address = "1.1.1.1"
  name             = "node_sid_profile"
  name_alias       = "name_alias"
  segment_id       = "1"
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
