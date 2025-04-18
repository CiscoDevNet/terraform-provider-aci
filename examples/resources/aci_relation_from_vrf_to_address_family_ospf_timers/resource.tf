
resource "aci_relation_from_vrf_to_address_family_ospf_timers" "example_vrf" {
  parent_dn        = aci_vrf.example.id
  address_family   = "ipv4-ucast"
  ospf_timers_name = aci_ospf_timers.example.name
}
