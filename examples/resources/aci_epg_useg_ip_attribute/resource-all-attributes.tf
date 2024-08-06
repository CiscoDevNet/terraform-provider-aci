
resource "aci_epg_useg_ip_attribute" "full_example_epg_useg_block_statement" {
  parent_dn      = aci_epg_useg_block_statement.example.id
  annotation     = "annotation"
  description    = "description_1"
  ip             = "131.107.1.200"
  name           = "131"
  name_alias     = "name_alias_1"
  owner_key      = "owner_key_1"
  owner_tag      = "owner_tag_1"
  use_epg_subnet = "yes"
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
