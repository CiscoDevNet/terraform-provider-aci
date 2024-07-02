
resource "aci_epg_useg_block_statement" "full_example_application_epg" {
  parent_dn   = aci_application_epg.example.id
  annotation  = "annotation"
  description = "description"
  match       = "all"
  name        = "criterion"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  precedence  = "1"
  scope       = "scope-bd"
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
