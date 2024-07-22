
data "aci_route_control_profile" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}

data "aci_route_control_profile" "example_l3_outside" {
  parent_dn = aci_l3_outside.example.id
  name      = "test_name"
}
