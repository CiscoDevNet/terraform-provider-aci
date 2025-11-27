
data "aci_copp_interface_protocol_policy" "example_copp_interface_policy" {
  parent_dn = aci_copp_interface_policy.example.id
  name      = "test_name"
}
