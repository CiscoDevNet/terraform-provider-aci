resource "aci_access_port_selector" "test_selector" {
    leaf_interface_profile_dn = aci_leaf_interface_profile.test_leaf_profile.id
    name = "tf_test"
    access_port_selector_type = "default"
  
}
