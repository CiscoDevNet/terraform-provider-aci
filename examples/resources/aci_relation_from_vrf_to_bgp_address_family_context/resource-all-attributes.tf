
resource "aci_relation_from_vrf_to_bgp_address_family_context" "full_example_vrf" {
  parent_dn                       = aci_vrf.example.id
  address_family                  = "ipv4-ucast"
  annotation                      = "annotation"
  bgp_address_family_context_name = aci_bgp_address_family_context.example.name
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
