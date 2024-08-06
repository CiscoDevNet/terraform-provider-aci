
resource "aci_relation_to_static_path" "example_application_epg" {
  parent_dn     = aci_application_epg.example.id
  encapsulation = "vlan-201"
  target_dn     = "topology/pod-1/paths-101/pathep-[eth1/1]"
}
