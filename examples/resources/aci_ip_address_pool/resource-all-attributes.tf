
resource "aci_ip_address_pool" "full_example" {
  gateway_address         = "10.0.0.1/24"
  address_type            = "regular"
  annotation              = "annotation"
  description             = "description_1"
  name                    = "test_name"
  name_alias              = "name_alias_1"
  owner_key               = "owner_key_1"
  owner_tag               = "owner_tag_1"
  skip_gateway_validation = "no"
  ip_address_blocks = [
    {
      annotation  = "annotation_1"
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
  ]
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
