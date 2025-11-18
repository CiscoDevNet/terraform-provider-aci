
resource "aci_relation_from_application_epg_to_attachable_access_entity_profile" "example_application_epg" {
  parent_dn                  = aci_application_epg.example.id
  tn_infra_att_entity_p_name = aci_.example.name
}
