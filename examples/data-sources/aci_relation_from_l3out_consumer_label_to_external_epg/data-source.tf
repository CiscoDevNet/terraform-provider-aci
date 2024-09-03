
data "aci_relation_from_l3out_consumer_label_to_external_epg" "example_l3out_consumer_label" {
  parent_dn = aci_l3out_consumer_label.example.id
  target_dn = "uni/tn-example_tenant/out-example_l3_outside/instP-inst_profile_2"
}
