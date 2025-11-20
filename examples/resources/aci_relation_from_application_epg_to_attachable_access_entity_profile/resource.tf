
resource "aci_relation_from_application_epg_to_attachable_access_entity_profile" "example_application_epg" {
  parent_dn                             = aci_application_epg.example.id
  encapsulation                         = "vlan-100"
  attachable_access_entity_profile_name = aci_attachable_access_entity_profile.example.name
}
