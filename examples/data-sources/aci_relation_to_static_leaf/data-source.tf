
data "aci_relation_to_static_leaf" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  target_dn = "topology/pod-1/node-101"
}
