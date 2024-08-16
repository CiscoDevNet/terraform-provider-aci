
resource "aci_vrf_fallback_route" "full_example_vrf_fallback_route_group" {
  parent_dn      = aci_vrf_fallback_route_group.example.id
  annotation     = "annotation"
  description    = "description_1"
  prefix_address = "2.2.2.3/24"
  name           = "name_1"
  name_alias     = "name_alias_1"
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
