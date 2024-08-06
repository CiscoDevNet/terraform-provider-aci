
resource "aci_endpoint_tag_ip" "full_example_tenant" {
  parent_dn    = aci_tenant.example.id
  annotation   = "annotation"
  vrf_name     = "test_ctx_name"
  id_attribute = "1"
  ip           = "10.0.0.2"
  name         = "name_1"
  name_alias   = "name_alias_1"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
