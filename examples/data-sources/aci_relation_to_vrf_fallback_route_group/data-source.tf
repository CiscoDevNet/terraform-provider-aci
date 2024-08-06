
data "aci_relation_to_vrf_fallback_route_group" "example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  target_dn = aci_vrf_fallback_route_group.example_vrf_fallback_route_group.id
}
