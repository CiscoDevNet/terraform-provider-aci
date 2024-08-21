
resource "aci_relation_to_route_control_profile" "full_example_l3out_consumer_label" {
  parent_dn  = aci_l3out_consumer_label.example.id
  annotation = "annotation"
  direction  = "import"
  target_dn  = "uni/tn-example_tenant/prof-rt_ctrl_profile"
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