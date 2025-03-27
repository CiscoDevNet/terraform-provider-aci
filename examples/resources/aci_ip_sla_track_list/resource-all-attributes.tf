
resource "aci_ip_sla_track_list" "full_example_tenant" {
  parent_dn       = aci_tenant.example.id
  annotation      = "annotation"
  description     = "description_1"
  name            = "test_name"
  name_alias      = "name_alias_1"
  owner_key       = "owner_key_1"
  owner_tag       = "owner_tag_1"
  percentage_down = "30"
  percentage_up   = "40"
  type            = "weight"
  weight_down     = "10"
  weight_up       = "20"
  relation_to_ip_sla_track_members = [
    {
      annotation = "annotation_1"
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
  ]
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
