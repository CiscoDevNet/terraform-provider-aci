
data "aci_spine_access_port_selector" "example_spine_interface_profile" {
  parent_dn          = aci_spine_interface_profile.example.id
  name               = "test_name"
  port_selector_type = "ALL"
}
