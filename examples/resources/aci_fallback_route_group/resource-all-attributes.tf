
resource "aci_fallback_route_group" "full_example_vrf" {
  parent_dn   = aci_vrf.example.id
  annotation  = "annotation"
  description = "description"
  name        = "fallback_route_group"
  name_alias  = "name_alias"
  fallback_members = [
    {
      annotation           = "annotation_1"
      description          = "description_1"
      name                 = "name_1"
      name_alias           = "name_alias_1"
      fallback_member_addr = "2.2.2.2"
    }
  ]
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
