
data "aci_relation_from_l3out_consumer_label_to_route_control_profile" "example_l3out_consumer_label" {
  parent_dn = aci_l3out_consumer_label.example.id
  direction = "import"
  target_dn = aci_route_control_profile.example.id
}
