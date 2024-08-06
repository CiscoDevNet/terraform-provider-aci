
resource "aci_epg_useg_ad_group_attribute" "full_example_epg_useg_block_statement" {
  parent_dn   = aci_epg_useg_block_statement.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "name_1"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  selector    = "adepg/authsvr-common-sg1-ISE_1/grpcont/dom-cisco.com/grp-Eng"
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
