
resource "aci_l3out_redistribute_policy" "full_example_l3_outside" {
  parent_dn                  = aci_l3_outside.example.id
  annotation                 = "annotation"
  source                     = "direct"
  route_control_profile_name = aci_route_control_profile.example.name
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
