resource "aci_interface_fc_policy" "test_fc" {
  name         = "tf_pol"
  rx_bb_credit = "64"
  speed        = "4G"
  trunk_mode   = "trunk-on"
}
