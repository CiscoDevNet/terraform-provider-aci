
resource "aci_endpoint_tag_ip" "full_example_tenant" {
  parent_dn    = aci_tenant.example.id
  annotation   = "annotation"
  vrf_name     = "test_ctx_name"
  id_attribute = "1"
  ip           = "10.0.0.2"
  name         = "name"
  name_alias   = "name_alias"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "tags_1"
      value = "value_1"
    }
  ]
}
