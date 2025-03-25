
resource "aci_relation_to_ip_sla_track_member" "full_example_ip_sla_track_list" {
  parent_dn  = aci_ip_sla_track_list.example.id
  annotation = "annotation"
  target_dn  = aci_ip_sla_track_member.example.id
  weight     = "20"
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
