
resource "aci_pim_route_map_policy" "example" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

