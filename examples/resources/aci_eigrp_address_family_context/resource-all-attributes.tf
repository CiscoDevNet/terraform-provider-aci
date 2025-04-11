
resource "aci_eigrp_address_family_context" "full_example_tenant" {
  parent_dn          = aci_tenant.example.id
  active_interval    = "3"
  annotation         = "annotation"
  description        = "description_1"
  external_distance  = "170"
  internal_distance  = "90"
  maximum_path_limit = "8"
  metric_style       = "narrow"
  name               = "test_name"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
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
