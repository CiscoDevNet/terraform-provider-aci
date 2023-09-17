
resource "aci_l3out_redistribute_policy" "example" {
  parent_dn          = aci_l3_outside.example.id
  src                = "direct"
  route_profile_name = "test_tn_rtctrl_profile_name"
  annotations = [
    {
      key   = "test_key"
      value = "test_value"
    },
  ]
}

