
resource "aci_fallback_member" "example_fallback_route_group" {
  parent_dn            = aci_fallback_route_group.example.id
  fallback_member_addr = "2.2.2.3"
}
