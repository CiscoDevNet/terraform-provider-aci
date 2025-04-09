
data "aci_relation_from_vrf_to_eigrp_address_family_context" "example_vrf" {
  parent_dn                         = aci_vrf.example.id
  address_family                    = "ipv4-ucast"
  eigrp_address_family_context_name = aci_eigrp_address_family_context.example.name
}
