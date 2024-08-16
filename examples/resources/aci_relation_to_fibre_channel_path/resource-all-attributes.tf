
resource "aci_relation_to_fibre_channel_path" "full_example_application_epg" {
  parent_dn   = aci_application_epg.example.id
  annotation  = "annotation"
  description = "description_1"
  target_dn   = "topology/pod-1/paths-101/pathep-[eth1/1]"
  vsan        = "vsan-10"
  vsan_mode   = "native"
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
