
resource "aci_relation_to_instance_profile" "example_l3out_consumer_label" {
  parent_dn = aci_l3out_consumer_label.example.id
  target_dn = aci_external_network_instance_profile.test.id
}
