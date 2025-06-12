
resource "aci_fabric_external_connection_policy" "full_example_tenant" {
  parent_dn    = aci_tenant.example.id
  annotation   = "annotation"
  description  = "description_1"
  id_attribute = "1"
  name         = "name_1"
  name_alias   = "name_alias_1"
  owner_key    = "owner_key_1"
  owner_tag    = "owner_tag_1"
  community    = "extended:as2-nn4:5:16"
  site_id      = "0"
  peering_profile = {
    annotation  = "annotation_1"
    description = "description_1"
    name        = "name_1"
    name_alias  = "name_alias_1"
    owner_key   = "owner_key_1"
    owner_tag   = "owner_tag_1"
    password    = "password_1"
    type        = "automatic_with_full_mesh"
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
