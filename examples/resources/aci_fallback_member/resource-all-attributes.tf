
resource "aci_fallback_member" "full_example_fallback_route_group" {
  parent_dn            = aci_fallback_route_group.example.id
  annotation           = "annotation"
  description          = "description"
  name                 = "name"
  name_alias           = "name_alias"
  fallback_member_addr = "2.2.2.3"
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
