
resource "aci_spine_access_port_selector" "full_example_spine_interface_profile" {
  parent_dn          = aci_spine_interface_profile.example.id
  annotation         = "annotation"
  description        = "description_1"
  name               = "test_name"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
  port_selector_type = "ALL"
  relation_to_spine_port_policy_group = {
    annotation = "annotation_1"
    target_dn  = aci_spine_port_policy_group.test_spine_port_policy_group_0.id
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
