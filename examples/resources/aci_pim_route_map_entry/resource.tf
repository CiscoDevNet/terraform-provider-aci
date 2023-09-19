
resource "aci_pim_route_map_entry" "example" {
  parent_dn = aci_pim_route_map_policy.example.id
  order     = "1"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

