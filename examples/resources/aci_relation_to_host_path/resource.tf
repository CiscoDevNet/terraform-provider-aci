
resource "aci_relation_to_host_path" "example_access_interface_override" {
  parent_dn = aci_access_interface_override.example.id
  target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
}
