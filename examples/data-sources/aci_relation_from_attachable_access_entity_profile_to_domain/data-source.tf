
data "aci_relation_from_attachable_access_entity_profile_to_domain" "example_attachable_access_entity_profile" {
  parent_dn = aci_attachable_access_entity_profile.example.id
  target_dn = aci_physical_domain.example.id
}
