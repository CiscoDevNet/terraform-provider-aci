
resource "aci_l3out_to_fallback_route_group" "full_example_l3_outside" {
  parent_dn  = aci_l3_outside.example.id
  annotation = "annotation"
  target_dn  = aci_vrf_fallback_route_group.test.id
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
