
data "aci_access_port_selector" "example_leaf_interface_profile" {
  parent_dn          = aci_leaf_interface_profile.example.id
  name               = "test_name"
  port_selector_type = "range"
}

data "aci_access_port_selector" "example_fex_profile" {
  parent_dn          = aci_fex_profile.example.id
  name               = "test_name"
  port_selector_type = "range"
}
