
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
  relation_from_l3out_consumer_label_to_external_epgs = [
    {
      annotation = "annotation_1"
      target_dn  = "uni/tn-test_tenant/out-test_l3_outside/instP-inst_profile"
    }
  ]
  relation_from_l3out_consumer_label_to_route_control_profiles = [
    {
      annotation = "annotation_1"
      direction  = "export"
      target_dn  = "uni/tn-test_tenant/prof-rt_ctrl_profile"
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
