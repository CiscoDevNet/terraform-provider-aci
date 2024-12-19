
resource "aci_dot1p_classifier" "full_example_custom_qos_policy" {
  parent_dn   = aci_custom_qos_policy.example.id
  annotation  = "annotation"
  description = "description_1"
  from        = "1"
  name        = "name_1"
  name_alias  = "name_alias_1"
  priority    = "level1"
  target      = "AF11"
  target_cos  = "0"
  to          = "2"
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
