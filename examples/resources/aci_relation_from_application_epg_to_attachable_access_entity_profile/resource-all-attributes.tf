
resource "aci_relation_from_application_epg_to_attachable_access_entity_profile" "full_example_application_epg" {
  parent_dn                             = aci_application_epg.example.id
  encapsulation                         = "vlan-100"
  deployment_immediacy                  = "immediate"
  mode                                  = "native"
  primary_encapsulation                 = "vlan-200"
  attachable_access_entity_profile_name = aci_attachable_access_entity_profile.example.name
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
