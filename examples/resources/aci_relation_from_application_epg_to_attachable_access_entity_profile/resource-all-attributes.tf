
resource "aci_relation_from_application_epg_to_attachable_access_entity_profile" "full_example_application_epg" {
  parent_dn                  = aci_application_epg.example.id
  annotation                 = "annotation"
  encapsulation              = "encapsulation_1"
  deployment_immediacy       = "immediate"
  mode                       = "native"
  primary_encapsulation      = "primary_encapsulation_1"
  tn_infra_att_entity_p_name = aci_.example.name
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
