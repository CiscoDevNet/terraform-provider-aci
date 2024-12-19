
resource "aci_custom_qos_policy" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  dot1p_classifiers = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      from        = "0"
      name        = "name_1"
      name_alias  = "name_alias_1"
      priority    = "level1"
      target      = "AF11"
      target_cos  = "0"
      to          = "0"
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
  ]
  dscp_to_priority_maps = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      from        = "AF11"
      name        = "name_1"
      name_alias  = "name_alias_1"
      priority    = "level1"
      target      = "AF11"
      target_cos  = "0"
      to          = "AF11"
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
  ]
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
