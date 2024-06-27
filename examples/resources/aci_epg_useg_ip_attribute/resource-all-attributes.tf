
resource "aci_epg_useg_ip_attribute" "full_example_epg_useg_block_statement" {
  parent_dn     = aci_epg_useg_block_statement.example.id
  annotation    = "annotation"
  description   = "description"
  ip            = "131.107.1.200"
  name          = "131"
  name_alias    = "name_alias"
  owner_key     = "owner_key"
  owner_tag     = "owner_tag"
  use_fv_subnet = "yes"
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
