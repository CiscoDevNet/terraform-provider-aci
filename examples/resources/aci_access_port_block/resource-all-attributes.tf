
resource "aci_access_port_block" "full_example_access_port_selector" {
  parent_dn   = aci_access_port_selector.example.id
  annotation  = "annotation"
  description = "description_1"
  from_card   = "2"
  from_port   = "3"
  name        = "test_name"
  name_alias  = "name_alias_1"
  to_card     = "4"
  to_port     = "5"
  relation_to_pc_vpc_override_policy = {
    annotation = "annotation_1"
    target_dn  = aci_leaf_access_bundle_policy_sub_group.example.id
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

resource "aci_access_port_block" "full_example_spine_access_port_selector" {
  parent_dn   = aci_spine_access_port_selector.example.id
  annotation  = "annotation"
  description = "description_1"
  from_card   = "2"
  from_port   = "3"
  name        = "test_name"
  name_alias  = "name_alias_1"
  to_card     = "4"
  to_port     = "5"
  relation_to_pc_vpc_override_policy = {
    annotation = "annotation_1"
    target_dn  = aci_leaf_access_bundle_policy_sub_group.example.id
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
