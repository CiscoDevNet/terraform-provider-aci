
resource "aci_relation_from_bridge_domain_to_l3_outside" "example_bridge_domain" {
  parent_dn       = aci_bridge_domain.example.id
  l3_outside_name = aci_l3_outside.example.name
}
