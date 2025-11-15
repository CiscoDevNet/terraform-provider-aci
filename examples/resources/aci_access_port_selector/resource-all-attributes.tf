
resource "aci_access_port_selector" "full_example_leaf_interface_profile" {
  parent_dn          = aci_leaf_interface_profile.example.id
  annotation         = "annotation"
  description        = "description_1"
  name               = "test_name"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
  port_selector_type = "range"
  relation_to_leaf_access_port_policy_group = {
    annotation = "annotation_1"
    fex_id     = "102"
    target_dn  = aci_leaf_access_port_policy_group.test_leaf_access_port_policy_group_0.id
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

resource "aci_access_port_selector" "full_example_fex_profile" {
  parent_dn          = aci_fex_profile.example.id
  annotation         = "annotation"
  description        = "description_1"
  name               = "test_name"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
  port_selector_type = "range"
  relation_to_leaf_access_port_policy_group = {
    annotation = "annotation_1"
    fex_id     = "102"
    target_dn  = aci_leaf_access_port_policy_group.test_leaf_access_port_policy_group_0.id
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
