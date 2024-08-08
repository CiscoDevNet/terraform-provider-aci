
resource "aci_route_control_profile" "full_example_tenant" {
  parent_dn                  = aci_tenant.example.id
  annotation                 = "annotation"
  route_map_continue         = "no"
  description                = "description_1"
  name                       = "test_name"
  name_alias                 = "name_alias_1"
  owner_key                  = "owner_key_1"
  owner_tag                  = "owner_tag_1"
  route_control_profile_type = "combinable"
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

resource "aci_route_control_profile" "full_example_l3_outside" {
  parent_dn                  = aci_l3_outside.example.id
  annotation                 = "annotation"
  route_map_continue         = "no"
  description                = "description_1"
  name                       = "test_name"
  name_alias                 = "name_alias_1"
  owner_key                  = "owner_key_1"
  owner_tag                  = "owner_tag_1"
  route_control_profile_type = "combinable"
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
