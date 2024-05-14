
resource "aci_relation_to_instance_profile" "full_example_l3out_consumer_label" {
  parent_dn  = aci_l3out_consumer_label.example.id
  annotation = "annotation"
  target_dn  = aci_external_network_instance_profile.example.id
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
