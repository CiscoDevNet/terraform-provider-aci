
resource "aci_relation_to_host_path" "full_example_host_path_selector" {
  parent_dn  = aci_host_path_selector.example.id
  annotation = "annotation"
  target_dn  = "topology/pod-1/paths-101/pathep-[eth1/1]"
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
