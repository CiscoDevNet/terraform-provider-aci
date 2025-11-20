
resource "aci_ip_address_block" "full_example_ip_address_pool" {
  parent_dn   = aci_ip_address_pool.example.id
  annotation  = "annotation"
  description = "description_1"
  from        = "10.0.0.2"
  name        = "name_1"
  name_alias  = "name_alias_1"
  to          = "10.0.0.5"
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
