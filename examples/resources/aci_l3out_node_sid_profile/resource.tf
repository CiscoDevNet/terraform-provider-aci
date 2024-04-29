
resource "aci_l3out_node_sid_profile" "example_l3out_loopback_interface_profile" {
  parent_dn  = aci_l3out_loopback_interface_profile.example.id
  sid_offset = "1"
}
