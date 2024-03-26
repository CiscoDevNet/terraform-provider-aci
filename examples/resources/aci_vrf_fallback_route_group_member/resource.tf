
resource "aci_vrf_fallback_route_group_member" "example_vrf_fallback_route_group" {
  parent_dn       = aci_vrf_fallback_route_group.example.id
  fallback_member = "2.2.2.3"
}
