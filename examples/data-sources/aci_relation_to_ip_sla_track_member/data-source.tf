
data "aci_relation_to_ip_sla_track_member" "example_ip_sla_track_list" {
  parent_dn = aci_ip_sla_track_list.example.id
  target_dn = aci_ip_sla_track_member.example.id
}
