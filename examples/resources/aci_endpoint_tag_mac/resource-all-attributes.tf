
resource "aci_endpoint_tag_mac" "full_example_tenant" {
  parent_dn    = aci_tenant.example.id
  annotation   = "annotation"
  bd_name      = "test_bd_name"
  id_attribute = "1"
  mac          = "00:00:00:00:00:01"
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
