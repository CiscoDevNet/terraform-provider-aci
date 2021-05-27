resource "aci_lldp_interface_policy" "test_lldp" {
  description = "example description"
  name        = "demo_lldp_pol"
  admin_rx_st = "enabled"
  admin_tx_st = "enabled"
  annotation  = "tag_lldp"
  name_alias  = "alias_lldp"

}
