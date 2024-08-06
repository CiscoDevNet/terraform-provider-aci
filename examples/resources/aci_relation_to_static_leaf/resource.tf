
resource "aci_relation_to_static_leaf" "example_application_epg" {
  parent_dn     = aci_application_epg.example.id
  encapsulation = "vlan-101"
  target_dn     = "topology/pod-1/node-101"
}
