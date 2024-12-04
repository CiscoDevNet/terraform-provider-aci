
resource "aci_dscp_to_priority_map" "full_example_custom_qos_policy" {
  parent_dn   = aci_custom_qos_policy.example.id
  annotation  = "annotation"
  description = "description_1"
  from        = "AF11"
  name        = "name_1"
  name_alias  = "name_alias_1"
  priority    = "level1"
  target      = "AF11"
  target_cos  = "0"
  to          = "AF22"
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
