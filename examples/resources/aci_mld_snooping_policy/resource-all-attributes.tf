
resource "aci_mld_snooping_policy" "full_example_tenant" {
  parent_dn            = aci_tenant.example.id
  admin_state          = "disabled"
  annotation           = "annotation"
  control              = ["fast-leave", "querier"]
  description          = "description_1"
  last_member_interval = "3"
  name                 = "test_name"
  name_alias           = "name_alias_1"
  owner_key            = "owner_key_1"
  owner_tag            = "owner_tag_1"
  query_interval       = "140"
  response_interval    = "11"
  start_query_count    = "9"
  start_query_interval = "2"
  ver                  = "unspecified"
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
