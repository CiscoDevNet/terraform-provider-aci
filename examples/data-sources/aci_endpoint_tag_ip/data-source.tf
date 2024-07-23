
data "aci_endpoint_tag_ip" "example_tenant" {
  parent_dn = aci_tenant.example.id
  vrf_name = "test_ctx_name"
  ip = "10.0.0.2"
}
