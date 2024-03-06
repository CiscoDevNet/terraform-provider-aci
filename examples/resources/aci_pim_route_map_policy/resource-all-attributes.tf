
resource "aci_pim_route_map_policy" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  annotation  = "annotation"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "tags_1"
      value = "value_1"
    }
  ]
}
