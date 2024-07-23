
data "aci_endpoint_tag_mac" "example_tenant" {
  parent_dn = aci_tenant.example.id
  bd_name   = "test_bd_name"
  mac       = "00:00:00:00:00:01"
}
