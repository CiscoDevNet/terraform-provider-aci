
resource "aci_lacp_enhanced_lag_policy" "full_example_vswitch_policy" {
  parent_dn           = aci_vswitch_policy.example.id
  annotation          = "annotation"
  load_balancing_mode = "dst-ip"
  mode                = "active"
  name                = "test_name"
  name_alias          = "name_alias_1"
  number_of_links     = "2"
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
