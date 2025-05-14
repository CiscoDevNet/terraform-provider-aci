
resource "aci_relation_from_external_epg_to_route_control_profile" "full_example_external_epg" {
  parent_dn                  = aci_external_epg.example.id
  annotation                 = "annotation"
  direction                  = "import"
  route_control_profile_name = aci_route_control_profile.example.name
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
