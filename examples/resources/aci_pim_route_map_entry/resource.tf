
resource "aci_pim_route_map_entry" "example_pim_route_map_policy" {
  parent_dn = aci_pim_route_map_policy.example.id
  order     = "1"
}
