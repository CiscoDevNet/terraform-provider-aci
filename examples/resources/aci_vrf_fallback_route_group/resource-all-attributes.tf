
resource "aci_vrf_fallback_route_group" "full_example_vrf" {
  parent_dn   = aci_vrf.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "fallback_route_group"
  name_alias  = "name_alias_1"
  vrf_fallback_route_group_members = [
    {
      annotation      = "annotation_1"
      description     = "description_1"
      name            = "name_1"
      name_alias      = "name_alias_1"
      fallback_member = "2.2.2.2"
    }
  ]
  vrf_fallback_routes = [
    {
      annotation     = "annotation_1"
      description    = "description_1"
      prefix_address = "2.2.2.2/24"
      name           = "name_1"
      name_alias     = "name_alias_1"
    }
  ]
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
