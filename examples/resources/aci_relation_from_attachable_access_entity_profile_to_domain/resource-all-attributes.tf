
resource "aci_relation_from_attachable_access_entity_profile_to_domain" "full_example_attachable_access_entity_profile" {
  parent_dn  = aci_attachable_access_entity_profile.example.id
  annotation = "annotation"
  target_dn  = aci_physical_domain.example.id
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
