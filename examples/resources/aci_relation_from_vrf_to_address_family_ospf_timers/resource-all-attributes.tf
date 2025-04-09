
resource "aci_relation_from_vrf_to_address_family_ospf_timers" "full_example_vrf" {
  parent_dn        = aci_vrf.example.id
  address_family   = "ipv4-ucast"
  annotation       = "annotation"
  ospf_timers_name = aci_ospf_timers.example.name
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
