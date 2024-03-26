
resource "aci_vrf_fallback_route_group_member" "full_example_vrf_fallback_route_group" {
  parent_dn       = aci_vrf_fallback_route_group.example.id
  annotation      = "annotation"
  description     = "description"
  name            = "name"
  name_alias      = "name_alias"
  fallback_member = "2.2.2.3"
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
