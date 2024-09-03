
resource "aci_relation_from_l3out_consumer_label_to_external_epg" "full_example_l3out_consumer_label" {
  parent_dn  = aci_l3out_consumer_label.example.id
  annotation = "annotation"
  target_dn  = "uni/tn-example_tenant/out-example_l3_outside/instP-inst_profile_2"
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
