
resource "aci_lacp_enhanced_lag_policy" "example_vswitch_policy" {
  parent_dn = aci_vswitch_policy.example.id
  name      = "test_name"
}
