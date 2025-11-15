
data "aci_access_port_block" "example_access_port_selector" {
  parent_dn = aci_access_port_selector.example.id
  name      = "test_name"
}

data "aci_access_port_block" "example_spine_access_port_selector" {
  parent_dn = aci_spine_access_port_selector.example.id
  name      = "test_name"
}
