
data "aci_vrf" "example_tenant" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
}
