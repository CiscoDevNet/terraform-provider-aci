
resource "aci_l3out_consumer_label" "full_example_l3_outside" {
  parent_dn   = aci_l3_outside.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner       = "infra"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  tag         = "lemon-chiffon"
  relation_to_external_epgs = [
    {
      annotation = "annotation_1"
      target_dn  = aci_external_network_instance_profile.example.id
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
  relation_to_route_control_profiles = [
    {
      annotation = "annotation_1"
      direction  = "export"
      target_dn  = aci_route_control_profile.example.id
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
