
resource "aci_relation_to_host_path" "example_host_path_selector" {
  parent_dn = aci_host_path_selector.example.id
  target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
}
