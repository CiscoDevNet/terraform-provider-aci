
resource "aci_epg_useg_criterion_mac_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  mac       = "AA:BB:CC:DD:EE:FF"
  name      = "mac_attr"
}
