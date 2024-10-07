
resource "aci_relation_from_bridge_domain_to_l3_outside" "full_example_bridge_domain" {
  parent_dn       = aci_bridge_domain.example.id
  annotation      = "annotation"
  l3_outside_name = aci_l3_outside.example.name
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
