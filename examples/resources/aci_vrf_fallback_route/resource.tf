
resource "aci_vrf_fallback_route" "example_vrf_fallback_route_group" {
  parent_dn      = aci_vrf_fallback_route_group.example.id
  prefix_address = "2.2.2.3/24"
}
