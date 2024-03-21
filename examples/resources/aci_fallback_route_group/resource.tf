
resource "aci_fallback_route_group" "example_vrf" {
  parent_dn = aci_vrf.example.id
  name      = "fallback_route_group"
}
