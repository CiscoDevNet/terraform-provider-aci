
data "aci_relation_from_external_epg_to_route_control_profile" "example_external_epg" {
  parent_dn                  = aci_external_epg.example.id
  direction                  = "import"
  route_control_profile_name = aci_route_control_profile.example.name
}
