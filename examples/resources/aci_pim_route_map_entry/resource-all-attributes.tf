
resource "aci_pim_route_map_entry" "full_example_pim_route_map_policy" {
  parent_dn           = aci_pim_route_map_policy.example.id
  action              = "deny"
  annotation          = "annotation"
  description         = "description_1"
  group_ip            = "0.0.0.0"
  name                = "name_1"
  name_alias          = "name_alias_1"
  order               = "1"
  rendezvous_point_ip = "0.0.0.0"
  source_ip           = "1.1.1.1/30"
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
