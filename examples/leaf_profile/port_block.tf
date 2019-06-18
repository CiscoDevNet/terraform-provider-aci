resource "aci_access_port_block" "test_port_block" {
  access_port_selector_dn = "${aci_access_port_selector.test_selector.id}"
  name = "tf_test_block"
}
