
resource "aci_pim_route_map_entry" "full_example_pim_route_map_policy" {
  parent_dn           = aci_pim_route_map_policy.example.id
  action              = "deny"
  annotation          = "annotation"
  description         = "description"
  group_ip            = "0.0.0.0"
  name                = "name"
  name_alias          = "name_alias"
  order               = "1"
  rendezvous_point_ip = "0.0.0.0"
  source_ip           = "1.1.1.1/30"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "tags_1"
      value = "value_1"
    }
  ]
}
